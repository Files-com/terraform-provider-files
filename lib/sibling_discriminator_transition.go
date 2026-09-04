package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type deprecatedSiblingDiscriminatorValueValidator struct {
	rule     string
	selector string
	variants []JSONSchemaVariant
	warning  string
}

func DeprecatedSiblingDiscriminatorValue(rule, selector string, variants []JSONSchemaVariant, warning string) validator.Dynamic {
	return deprecatedSiblingDiscriminatorValueValidator{rule: rule, selector: selector, variants: variants, warning: warning}
}

func (v deprecatedSiblingDiscriminatorValueValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v deprecatedSiblingDiscriminatorValueValidator) MarkdownDescription(_ context.Context) string {
	return v.warning
}

func (v deprecatedSiblingDiscriminatorValueValidator) ValidateDynamic(ctx context.Context, req validator.DynamicRequest, resp *validator.DynamicResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() || req.ConfigValue.IsUnderlyingValueNull() || req.ConfigValue.IsUnderlyingValueUnknown() {
		return
	}
	terraformValue, err := req.ConfigValue.ToTerraformValue(ctx)
	if err != nil {
		resp.Diagnostics.AddAttributeError(req.Path, "Failed to inspect value", err.Error())
		return
	}
	if !terraformValue.IsFullyKnown() {
		return
	}

	var selector types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root(v.selector), &selector)...)
	if resp.Diagnostics.HasError() || selector.IsNull() || selector.IsUnknown() {
		return
	}

	value, diags := DynamicToInterface(ctx, req.Path, req.ConfigValue)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	legacy, diags := siblingDiscriminatorValueIsLegacy(req.Path, value, selector.ValueString(), v.variants)
	resp.Diagnostics.Append(diags...)
	if legacy && !resp.Diagnostics.HasError() {
		detail := fmt.Sprintf("Wrap value under %q: value = { %s = <documented %s value> }. Remove jsonencode(...) if present. %s", selector.ValueString(), selector.ValueString(), selector.ValueString(), v.warning)
		resp.Diagnostics.AddAttributeWarning(req.Path, fmt.Sprintf("Deprecated %s format: %s", v.rule, selector.ValueString()), detail)
	}
}

func DynamicSiblingDiscriminatorToAPI(ctx context.Context, attributePath path.Path, source types.Dynamic, selector string, variants []JSONSchemaVariant) (any, diag.Diagnostics) {
	value, diags := DynamicToInterface(ctx, attributePath, source)
	if diags.HasError() || value == nil {
		return value, diags
	}

	value, _ = decodeJSONString(value)
	legacy, shapeDiags := siblingDiscriminatorValueIsLegacy(attributePath, value, selector, variants)
	diags.Append(shapeDiags...)
	if diags.HasError() || legacy {
		return value, diags
	}

	variant, _ := findJSONSchemaVariant(variants, selector)
	return value.(map[string]interface{})[variant.Name], diags
}

func DynamicSiblingDiscriminatorUpdateToAPI(ctx context.Context, attributePath path.Path, source types.Dynamic, prior types.Dynamic, selector string, variants []JSONSchemaVariant) (any, diag.Diagnostics) {
	value, diags := DynamicSiblingDiscriminatorToAPI(ctx, attributePath, source, selector, variants)
	if diags.HasError() || value == nil {
		return value, diags
	}
	priorValue, priorDiags := DynamicSiblingDiscriminatorToAPI(ctx, attributePath, prior, selector, variants)
	diags.Append(priorDiags...)
	if diags.HasError() || priorValue == nil {
		return value, diags
	}

	valueObject, valueIsObject := value.(map[string]interface{})
	priorObject, priorIsObject := priorValue.(map[string]interface{})
	if !valueIsObject || !priorIsObject {
		return value, diags
	}
	variant, _ := findJSONSchemaVariant(variants, selector)
	for key := range priorObject {
		if _, exists := valueObject[key]; !exists && slices.Contains(variant.Writable, key) {
			valueObject[key] = nil
		}
	}
	return valueObject, diags
}

func APIToDynamicSiblingDiscriminator(ctx context.Context, attributePath path.Path, source any, selector string, variants []JSONSchemaVariant, prior types.Dynamic) (types.Dynamic, diag.Diagnostics) {
	if prior.IsNull() || prior.IsUnknown() || prior.IsUnderlyingValueNull() || prior.IsUnderlyingValueUnknown() {
		if source == nil {
			return types.DynamicNull(), nil
		}
		return ToDynamic(ctx, attributePath, source, nil)
	}

	priorValue, diags := DynamicToInterface(ctx, attributePath, prior)
	if diags.HasError() {
		return types.DynamicNull(), diags
	}
	decodedPriorValue, jsonEncoded := decodeJSONString(priorValue)
	value, shapeDiags := siblingDiscriminatorStateValue(ctx, attributePath, source, decodedPriorValue, selector, variants)
	diags.Append(shapeDiags...)
	if diags.HasError() {
		return types.DynamicNull(), diags
	}
	if value == nil {
		return types.DynamicNull(), diags
	}
	if jsonEncoded {
		if reflect.DeepEqual(value, decodedPriorValue) {
			return prior, diags
		}
		encoded, err := json.Marshal(value)
		if err != nil {
			diags.AddAttributeError(attributePath, "Failed to convert API value", "Could not encode the legacy JSON value: "+err.Error())
			return types.DynamicNull(), diags
		}
		return types.DynamicValue(types.StringValue(string(encoded))), diags
	}

	return ToDynamic(ctx, attributePath, value, prior.UnderlyingValue())
}

func siblingDiscriminatorValueIsLegacy(attributePath path.Path, source any, selector string, variants []JSONSchemaVariant) (bool, diag.Diagnostics) {
	variant, ok := findJSONSchemaVariant(variants, selector)
	if !ok {
		return false, siblingDiscriminatorTransitionError(attributePath, "unknown sibling discriminator "+selector)
	}
	object, ok := source.(map[string]interface{})
	if !ok {
		return true, nil
	}
	if len(object) == 1 {
		if value, ok := object[variant.Name]; ok {
			actualType := jsonSchemaRootType(value)
			if actualType != variant.RootType && !(variant.RootType == "number" && actualType == "integer") {
				return false, siblingDiscriminatorTransitionError(attributePath, fmt.Sprintf("value wrapper %q must contain %s, got %s; keep the legacy value unwrapped or use the documented new format", variant.Name, variant.RootType, actualType))
			}
			return false, nil
		}
	}

	for _, candidate := range variants {
		if _, ok := object[candidate.Name]; ok {
			return false, siblingDiscriminatorTransitionError(attributePath, fmt.Sprintf("value wrapper must contain exactly the %q key selected by %s", variant.Name, selector))
		}
	}
	return true, nil
}

func jsonSchemaRootType(value any) string {
	switch value := value.(type) {
	case map[string]interface{}:
		return "object"
	case []interface{}:
		return "array"
	case string:
		return "string"
	case bool:
		return "boolean"
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return "integer"
	case float32:
		if math.Trunc(float64(value)) == float64(value) {
			return "integer"
		}
		return "number"
	case float64:
		if math.Trunc(value) == value {
			return "integer"
		}
		return "number"
	case nil:
		return "null"
	default:
		return fmt.Sprintf("unsupported %T", value)
	}
}

func siblingDiscriminatorStateValue(ctx context.Context, attributePath path.Path, source any, prior any, selector string, variants []JSONSchemaVariant) (any, diag.Diagnostics) {
	legacy, diags := siblingDiscriminatorValueIsLegacy(attributePath, prior, selector, variants)
	if diags.HasError() {
		return nil, diags
	}
	variant, _ := findJSONSchemaVariant(variants, selector)
	if !legacy {
		wrapped, wrapDiags := WrapSiblingDiscriminator(ctx, attributePath, source, selector, variants)
		if wrapped == nil {
			wrapped = map[string]interface{}{}
		}
		wrappedValue := wrapped.(map[string]interface{})
		configuredValue := prior.(map[string]interface{})
		sourceValue, exists := wrappedValue[variant.Name]
		value, keep := configuredDynamicStateValue(sourceValue, exists, configuredValue[variant.Name], variant.WriteOnly, "")
		if !keep {
			return nil, wrapDiags
		}
		return map[string]interface{}{variant.Name: value}, wrapDiags
	}

	value := source
	_, priorIsObject := prior.(map[string]interface{})
	if variant.LegacyProperty != "" && !priorIsObject {
		if object, ok := source.(map[string]interface{}); ok {
			if projected, exists := object[variant.LegacyProperty]; exists {
				value = projected
			}
		}
		if variant.LegacyList {
			if _, priorIsList := prior.([]interface{}); !priorIsList {
				if values, ok := value.([]interface{}); ok && len(values) > 0 {
					value = values[0]
				}
			}
		}
	}
	value, keep := configuredDynamicStateValue(value, value != nil, prior, variant.WriteOnly, "")
	if !keep {
		return nil, diags
	}
	return value, diags
}

func configuredDynamicStateValue(source any, sourceExists bool, configured any, writeOnly []string, currentPath string) (any, bool) {
	switch configuredValue := configured.(type) {
	case map[string]interface{}:
		sourceValue := map[string]interface{}{}
		if sourceExists {
			var ok bool
			sourceValue, ok = source.(map[string]interface{})
			if !ok {
				return source, true
			}
		}
		result := make(map[string]interface{}, len(configuredValue))
		for key, configuredEntry := range configuredValue {
			entryPath := key
			if currentPath != "" {
				entryPath = currentPath + "." + key
			}
			sourceEntry, exists := sourceValue[key]
			if value, keep := configuredDynamicStateValue(sourceEntry, exists, configuredEntry, writeOnly, entryPath); keep {
				result[key] = value
			}
		}
		return result, sourceExists || len(result) > 0
	case []interface{}:
		if !sourceExists {
			return nil, false
		}
		sourceValue, ok := source.([]interface{})
		if !ok {
			return source, true
		}
		result := make([]interface{}, len(sourceValue))
		for index, sourceEntry := range sourceValue {
			if index < len(configuredValue) {
				result[index], _ = configuredDynamicStateValue(sourceEntry, true, configuredValue[index], writeOnly, currentPath)
			} else {
				result[index] = sourceEntry
			}
		}
		return result, true
	default:
		if sourceExists {
			return source, true
		}
		if slices.Contains(writeOnly, currentPath) {
			return configured, true
		}
		return nil, false
	}
}

func decodeJSONString(source any) (any, bool) {
	value, ok := source.(string)
	if !ok {
		return source, false
	}
	var decoded any
	if json.Unmarshal([]byte(value), &decoded) != nil {
		return source, false
	}
	return decoded, true
}

func siblingDiscriminatorTransitionError(attributePath path.Path, detail string) diag.Diagnostics {
	return diag.Diagnostics{diag.NewAttributeErrorDiagnostic(attributePath, "Invalid Value Wrapper", detail)}
}
