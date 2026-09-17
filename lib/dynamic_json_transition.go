package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	tfpath "github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type deprecatedJSONEncodingValidator struct {
	rule, example, deadline string
}

func DeprecatedJSONEncoding(rule, example, deadline string) validator.Dynamic {
	return deprecatedJSONEncodingValidator{rule: rule, example: example, deadline: deadline}
}

func (v deprecatedJSONEncodingValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v deprecatedJSONEncodingValidator) MarkdownDescription(_ context.Context) string {
	return "Use native HCL instead of a JSON-encoded string. JSON strings stop being accepted on " + v.deadline + "."
}

func (v deprecatedJSONEncodingValidator) ValidateDynamic(ctx context.Context, req validator.DynamicRequest, resp *validator.DynamicResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() || req.ConfigValue.IsUnderlyingValueNull() || req.ConfigValue.IsUnderlyingValueUnknown() {
		return
	}
	if _, encoded := req.ConfigValue.UnderlyingValue().(types.String); !encoded {
		return
	}
	_, diags := JSONValueToAPI(ctx, req.Path, req.ConfigValue)
	resp.Diagnostics.Append(diags...)
	if !resp.Diagnostics.HasError() {
		resp.Diagnostics.AddAttributeWarning(req.Path, "Deprecated JSON encoding: "+v.rule,
			fmt.Sprintf("Remove jsonencode(...) or replace the JSON string with native HCL. Keep the same keys and nesting: %s. Native values are compatible with the typed schema. JSON strings stop being accepted on %s, when the typed schema ships in a major provider release. No additional wrapper is needed.", v.example, v.deadline))
	}
}

func JSONValueToAPI(ctx context.Context, attributePath tfpath.Path, source attr.Value) (any, diag.Diagnostics) {
	if source.IsNull() || source.IsUnknown() {
		return nil, nil
	}
	value, diags := attributeToInterface(ctx, attributePath, source, false, true)
	if diags.HasError() || value == nil {
		return value, diags
	}
	if encoded, ok := value.(string); ok {
		if err := json.Unmarshal([]byte(encoded), &value); err != nil {
			diags.AddAttributeError(attributePath, "Invalid JSON value", "Expected valid JSON or a native HCL value. Remove jsonencode(...) when migrating. The JSON string could not be decoded.")
			return nil, diags
		}
	}
	return value, diags
}

func JSONToStringMapSlice(attributePath tfpath.Path, value any) ([]map[string]any, diag.Diagnostics) {
	if value == nil {
		return nil, nil
	}
	items, ok := value.([]any)
	if !ok {
		return nil, diag.Diagnostics{diag.NewAttributeErrorDiagnostic(attributePath, "Invalid JSON value", "Expected an array of objects.")}
	}
	result := make([]map[string]any, len(items))
	for index, item := range items {
		object, ok := item.(map[string]any)
		if !ok {
			return nil, diag.Diagnostics{diag.NewAttributeErrorDiagnostic(attributePath.AtListIndex(index), "Invalid JSON value", "Expected an object.")}
		}
		result[index] = object
	}
	return result, nil
}

func APIToDynamicJSON(ctx context.Context, attributePath tfpath.Path, source any, prior types.Dynamic, mapPaths, writeOnly []string, normalization string) (types.Dynamic, diag.Diagnostics) {
	// SDK array properties use []map[string]any; decode once into the same JSON tree as configuration.
	encodedSource, err := json.Marshal(source)
	if err != nil {
		return types.DynamicNull(), diag.Diagnostics{diag.NewAttributeErrorDiagnostic(attributePath, "Invalid API value", "Could not encode the API response as JSON.")}
	}
	if err := json.Unmarshal(encodedSource, &source); err != nil {
		return types.DynamicNull(), diag.Diagnostics{diag.NewAttributeErrorDiagnostic(attributePath, "Invalid API value", "Could not decode the API response as JSON.")}
	}
	if prior.IsNull() || prior.IsUnknown() || prior.IsUnderlyingValueNull() || prior.IsUnderlyingValueUnknown() {
		return ToDynamic(ctx, attributePath, source, nil)
	}
	configured, diags := JSONValueToAPI(ctx, attributePath, prior)
	if diags.HasError() {
		return types.DynamicNull(), diags
	}
	if normalization != "" {
		source = preserveJSONNormalization(source, configured, normalizedJSON(configured, normalization), normalization == "mount_mappings")
	}
	value, _ := configuredDynamicStateValue(source, source != nil, configured, writeOnly, "", mapPaths...)
	if reflect.DeepEqual(configured, value) {
		return prior, nil
	}
	if _, encoded := prior.UnderlyingValue().(types.String); encoded {
		result, err := json.Marshal(value)
		if err != nil {
			return types.DynamicNull(), diag.Diagnostics{diag.NewAttributeErrorDiagnostic(attributePath, "Invalid API value", "Could not encode the API value as a JSON string.")}
		}
		return types.DynamicValue(types.StringValue(string(result))), nil
	}
	result, diags := dynamicJSONAttribute(ctx, attributePath, value, prior.UnderlyingValue(), mapPaths, "")
	if diags.HasError() {
		return types.DynamicNull(), diags
	}
	return types.DynamicValue(result), nil
}

// Keep equivalent configured spellings, but still record changes to other fields.
func preserveJSONNormalization(source, configured, normalized any, normalizedKeys bool) any {
	if reflect.DeepEqual(source, normalized) {
		return configured
	}
	switch prior := configured.(type) {
	case map[string]any:
		values, ok := source.(map[string]any)
		canonical, canonicalObject := normalized.(map[string]any)
		if !ok || !canonicalObject {
			return source
		}
		for key, value := range prior {
			canonicalKey := key
			if normalizedKeys {
				canonicalKey = strings.TrimSpace(key)
			}
			current, exists := values[canonicalKey]
			if exists || canonical[canonicalKey] == nil {
				values[key] = preserveJSONNormalization(current, value, canonical[canonicalKey], false)
				if key != canonicalKey {
					delete(values, canonicalKey)
				}
			}
		}
	case []any:
		values, ok := source.([]any)
		canonical, canonicalList := normalized.([]any)
		if ok && canonicalList {
			for index := range values {
				if index < len(prior) && index < len(canonical) {
					values[index] = preserveJSONNormalization(values[index], prior[index], canonical[index], false)
				}
			}
		}
	}
	return source
}

func dynamicJSONAttribute(ctx context.Context, attributePath tfpath.Path, source any, prior attr.Value, mapPaths []string, currentPath string) (attr.Value, diag.Diagnostics) {
	if prior == nil {
		value, diags := ToDynamic(ctx, attributePath, source, nil)
		if value.IsNull() {
			return types.DynamicNull(), diags
		}
		return value.UnderlyingValue(), diags
	}
	if dynamic, ok := prior.(types.Dynamic); ok {
		value, diags := dynamicJSONAttribute(ctx, attributePath, source, dynamic.UnderlyingValue(), mapPaths, currentPath)
		return types.DynamicValue(value), diags
	}
	if source == nil {
		value, err := prior.Type(ctx).ValueFromTerraform(ctx, tftypes.NewValue(prior.Type(ctx).TerraformType(ctx), nil))
		if err != nil {
			return nil, diag.Diagnostics{diag.NewAttributeErrorDiagnostic(attributePath, "Invalid API value", err.Error())}
		}
		return value, nil
	}
	if object, ok := prior.(types.Object); ok {
		values, ok := source.(map[string]any)
		if !ok {
			return nil, diag.Diagnostics{diag.NewAttributeErrorDiagnostic(attributePath, "Invalid API value", "Expected the configured object representation.")}
		}
		if !slices.Contains(mapPaths, currentPath) {
			for name := range object.AttributeTypes(ctx) {
				if _, exists := values[name]; !exists {
					values[name] = nil
				}
			}
		}
		attributes := make(map[string]attr.Value, len(values))
		attributeTypes := make(map[string]attr.Type, len(values))
		var diags diag.Diagnostics
		for name, value := range values {
			entryPath := name
			if currentPath != "" {
				entryPath = currentPath + "." + name
			}
			entry, entryDiags := dynamicJSONAttribute(ctx, attributePath.AtName(name), value, object.Attributes()[name], mapPaths, entryPath)
			diags.Append(entryDiags...)
			if !entryDiags.HasError() {
				attributes[name], attributeTypes[name] = entry, entry.Type(ctx)
			}
		}
		if diags.HasError() {
			return nil, diags
		}
		return types.ObjectValue(attributeTypes, attributes)
	}
	if tuple, ok := prior.(types.Tuple); ok {
		values, ok := source.([]any)
		if !ok {
			return nil, diag.Diagnostics{diag.NewAttributeErrorDiagnostic(attributePath, "Invalid API value", "Expected the configured array representation.")}
		}
		elements := make([]attr.Value, len(values))
		elementTypes := make([]attr.Type, len(values))
		var diags diag.Diagnostics
		for index, value := range values {
			var old attr.Value
			if index < len(tuple.Elements()) {
				old = tuple.Elements()[index]
			}
			entry, entryDiags := dynamicJSONAttribute(ctx, attributePath.AtListIndex(index), value, old, mapPaths, currentPath)
			diags.Append(entryDiags...)
			if !entryDiags.HasError() {
				elements[index], elementTypes[index] = entry, entry.Type(ctx)
			}
		}
		if diags.HasError() {
			return nil, diags
		}
		return types.TupleValue(elementTypes, elements)
	}
	if collection, ok := prior.(types.Map); ok {
		values, ok := source.(map[string]any)
		if !ok {
			return nil, diag.Diagnostics{diag.NewAttributeErrorDiagnostic(attributePath, "Invalid API value", "Expected the configured map representation.")}
		}
		elements := make(map[string]attr.Value, len(values))
		elementTypes := make(map[string]attr.Type, len(values))
		var diags diag.Diagnostics
		for name, value := range values {
			entryPath := name
			if currentPath != "" {
				entryPath = currentPath + "." + name
			}
			entry, entryDiags := dynamicJSONAttribute(ctx, attributePath.AtMapKey(name), value, collection.Elements()[name], mapPaths, entryPath)
			diags.Append(entryDiags...)
			if !entryDiags.HasError() {
				elements[name], elementTypes[name] = entry, entry.Type(ctx)
			}
		}
		if diags.HasError() {
			return nil, diags
		}
		if uniformElementType(ctx, collection.ElementType(ctx), slices.Collect(maps.Values(elements))) {
			return types.MapValue(collection.ElementType(ctx), elements)
		}
		return types.ObjectValue(elementTypes, elements)
	}
	if list, ok := prior.(types.List); ok {
		return dynamicJSONCollection(ctx, attributePath, source, list.Elements(), list.ElementType(ctx), false, mapPaths, currentPath)
	}
	if set, ok := prior.(types.Set); ok {
		return dynamicJSONCollection(ctx, attributePath, source, set.Elements(), set.ElementType(ctx), true, mapPaths, currentPath)
	}
	switch prior.(type) {
	case types.Number, types.Int64, types.Float64, types.String, types.Bool:
		value, diags := ToDynamic(ctx, attributePath, source, prior)
		return value.UnderlyingValue(), diags
	}
	return schemaValue(ctx, attributePath, source, prior.Type(ctx))
}

// Typed lists and sets keep their element type while every element still matches it. Otherwise the value becomes a tuple so a changed API type shows as a diff.
func dynamicJSONCollection(ctx context.Context, attributePath tfpath.Path, source any, priorElements []attr.Value, elementType attr.Type, isSet bool, mapPaths []string, currentPath string) (attr.Value, diag.Diagnostics) {
	values, ok := source.([]any)
	if !ok {
		return nil, diag.Diagnostics{diag.NewAttributeErrorDiagnostic(attributePath, "Invalid API value", "Expected the configured array representation.")}
	}
	elements := make([]attr.Value, len(values))
	elementTypes := make([]attr.Type, len(values))
	var diags diag.Diagnostics
	for index, value := range values {
		var old attr.Value
		if index < len(priorElements) {
			old = priorElements[index]
		}
		entry, entryDiags := dynamicJSONAttribute(ctx, attributePath.AtListIndex(index), value, old, mapPaths, currentPath)
		diags.Append(entryDiags...)
		if !entryDiags.HasError() {
			elements[index], elementTypes[index] = entry, entry.Type(ctx)
		}
	}
	if diags.HasError() {
		return nil, diags
	}
	if !uniformElementType(ctx, elementType, elements) {
		return types.TupleValue(elementTypes, elements)
	}
	if isSet {
		return types.SetValue(elementType, elements)
	}
	return types.ListValue(elementType, elements)
}

func uniformElementType(ctx context.Context, elementType attr.Type, elements []attr.Value) bool {
	for _, element := range elements {
		if !element.Type(ctx).Equal(elementType) {
			return false
		}
	}
	return true
}

// Compare the documented Rails normalizations without rewriting requests or hiding unrelated drift.
// Each case mirrors a Rails method in files-rails. Change both sides together.
//
//	"mount_mappings": DesktopConfigurationProfile#normalize_mount_mappings and Path.normalize
//	"expected_remote_servers": IntegrationCentricProfile#normalize_expected_remote_server
//	"watermark": FolderBehavior::Watermark.normalize_integer_keys
//	"empty_object": json_property columns that store {} when the request sends null
func normalizedJSON(source any, normalization string) any {
	encoded, _ := json.Marshal(source)
	var value any
	if json.Unmarshal(encoded, &value) != nil {
		return source
	}
	if value == nil && (normalization == "empty_object" || normalization == "watermark") {
		return map[string]any{}
	}
	switch normalization {
	case "mount_mappings":
		if object, ok := value.(map[string]any); ok {
			result := make(map[string]any, len(object))
			for key, value := range object {
				if text, ok := value.(string); ok {
					parts := strings.FieldsFunc(strings.TrimSpace(text), func(char rune) bool { return char == '/' || char == '\\' })
					cleaned := parts[:0]
					for _, part := range parts {
						part = strings.ReplaceAll(part, "\x00", "")
						if part != "" && part != "." && part != ".." {
							cleaned = append(cleaned, part)
						}
					}
					value = strings.Join(cleaned, "/")
				}
				result[strings.TrimSpace(key)] = value
			}
			return result
		}
	case "expected_remote_servers":
		if entries, ok := value.([]any); ok {
			for _, entry := range entries {
				if object, ok := entry.(map[string]any); ok {
					switch name := object["name"].(type) {
					case float64:
						object["name"] = strconv.FormatFloat(name, 'f', -1, 64)
					case bool:
						object["name"] = strconv.FormatBool(name)
					}
					for _, key := range []string{"name", "server_type"} {
						if text, ok := object[key].(string); ok {
							object[key] = strings.TrimSpace(text)
						}
					}
					if object["name"] == "" {
						delete(object, "name")
					}
				}
			}
		}
	case "watermark":
		if object, ok := value.(map[string]any); ok {
			for _, key := range []string{"transparency", "max_height_or_width"} {
				if text, ok := object[key].(string); ok {
					if number, err := strconv.ParseInt(strings.TrimSpace(text), 10, 64); err == nil {
						object[key] = float64(number)
					}
				}
			}
		}
	}
	return value
}
