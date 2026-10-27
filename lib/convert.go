package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"math"
	"math/big"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type JSONSchemaVariant struct {
	Name           string
	Value          string
	RootType       string
	Required       []string
	Allowed        []string
	Writable       []string
	WriteOnly      []string
	LegacyProperty string
	LegacyList     bool
	ValueTypes     map[string]string
}

func TimeToStringType(ctx context.Context, path path.Path, source *time.Time, dest *types.String) error {
	ctx = setAttributePath(ctx, path)

	if source == nil {
		*dest = types.StringValue("")
	} else {
		if dest.IsNull() || dest.ValueString() == "" {
			*dest = types.StringValue(source.Format(time.RFC3339))
		} else {
			parsedDest, err := time.Parse(time.RFC3339, dest.ValueString())
			if err != nil {
				return err
			}

			utcSrc := source.UTC()
			utcDest := parsedDest.UTC()

			if utcSrc == utcDest {
				tflog.Info(ctx, "Skipping updating state with matching UTC time")
			} else {
				tflog.Info(ctx, "Updating state with new time")
				*dest = types.StringValue(source.Format(time.RFC3339))
			}
		}
	}

	return nil
}

func DynamicToStringMapSlice(ctx context.Context, path path.Path, source types.Dynamic) ([]map[string]interface{}, diag.Diagnostics) {
	if source.IsNull() || source.IsUnknown() || source.IsUnderlyingValueNull() || source.IsUnderlyingValueUnknown() {
		return nil, nil
	}

	ctx = setAttributePath(ctx, path)

	switch underlyingValue := source.UnderlyingValue().(type) {
	case types.Tuple:
		tflog.Info(ctx, "Converting TupleValue to StringMapSlice")
		return ListToStringMapSlice(ctx, path, underlyingValue.Elements())
	default:
		return nil, diag.Diagnostics{
			diag.NewAttributeErrorDiagnostic(
				path,
				"Failed to convert DynamicValue",
				"Unhandled type: "+underlyingValue.Type(ctx).String(),
			),
		}
	}
}

func ListToStringMapSlice(ctx context.Context, path path.Path, elements []attr.Value) (dest []map[string]interface{}, diags diag.Diagnostics) {
	dest = make([]map[string]interface{}, 0, len(elements))

	for i, element := range elements {
		attrValue, attrDiags := DynamicToStringMap(ctx, path.AtListIndex(i), types.DynamicValue(element))
		if attrDiags.HasError() {
			diags.Append(attrDiags...)
		} else {
			dest = append(dest, attrValue)
		}
	}

	return
}

func DynamicToStringMap(ctx context.Context, path path.Path, source types.Dynamic) (map[string]interface{}, diag.Diagnostics) {
	if source.IsNull() || source.IsUnknown() || source.IsUnderlyingValueNull() || source.IsUnderlyingValueUnknown() {
		return nil, nil
	}

	ctx = setAttributePath(ctx, path)

	switch underlyingValue := source.UnderlyingValue().(type) {
	case types.Object:
		tflog.Info(ctx, "Converting ObjectValue to StringMap")
		return ElementsToStringMap(ctx, path, underlyingValue.Attributes())
	default:
		return nil, diag.Diagnostics{
			diag.NewAttributeErrorDiagnostic(
				path,
				"Failed to convert DynamicValue",
				"Unhandled type: "+source.Type(ctx).String(),
			),
		}
	}
}

func DynamicToInterface(ctx context.Context, path path.Path, source types.Dynamic) (interface{}, diag.Diagnostics) {
	if source.IsNull() || source.IsUnknown() || source.IsUnderlyingValueNull() || source.IsUnderlyingValueUnknown() {
		return nil, nil
	}

	ctx = setAttributePath(ctx, path)

	return AttributeToInterface(ctx, path, source.UnderlyingValue())
}

func ElementsToStringMap(ctx context.Context, path path.Path, attrs map[string]attr.Value) (dest map[string]interface{}, diags diag.Diagnostics) {
	return attributesToMap(ctx, path, attrs, false, false)
}

func attributesToMap(ctx context.Context, path path.Path, attrs map[string]attr.Value, omitNulls, preserveNulls bool) (dest map[string]interface{}, diags diag.Diagnostics) {
	dest = make(map[string]interface{})
	for key, value := range attrs {
		if omitNulls && (value.IsNull() || value.IsUnknown()) {
			continue
		}
		attrValue, attrDiags := attributeToInterface(ctx, path.AtMapKey(key), value, omitNulls, preserveNulls)
		if attrDiags.HasError() {
			diags.Append(attrDiags...)
		} else {
			dest[key] = attrValue
		}
	}

	return
}

func AttributeToInterface(ctx context.Context, path path.Path, source attr.Value) (dest interface{}, diags diag.Diagnostics) {
	return attributeToInterface(ctx, path, source, false, false)
}

func SchemaAttributeToInterface(ctx context.Context, path path.Path, source attr.Value) (dest interface{}, diags diag.Diagnostics) {
	return attributeToInterface(ctx, path, source, true, false)
}

func attributeToInterface(ctx context.Context, path path.Path, source attr.Value, omitNulls, preserveNulls bool) (dest interface{}, diags diag.Diagnostics) {
	if (omitNulls || preserveNulls) && (source.IsNull() || source.IsUnknown()) {
		return nil, nil
	}
	ctx = setAttributePath(ctx, path)

	switch actualValue := source.(type) {
	case types.Bool:
		tflog.Info(ctx, "Converting BoolValue to bool")
		dest = actualValue.ValueBool()
	case types.String:
		tflog.Info(ctx, "Converting StringValue to string")
		dest = actualValue.ValueString()
	case types.Number:
		tflog.Info(ctx, "Converting NumberValue to float64")
		dest, _ = actualValue.ValueBigFloat().Float64()
	case types.Int64:
		tflog.Info(ctx, "Converting Int64Value to int64")
		dest = actualValue.ValueInt64()
	case types.Float64:
		tflog.Info(ctx, "Converting Float64Value to float64")
		dest = actualValue.ValueFloat64()
	case types.Object:
		tflog.Info(ctx, "Converting ObjectValue to map")
		dest, diags = attributesToMap(ctx, path, actualValue.Attributes(), omitNulls, preserveNulls)
	case types.Map:
		tflog.Info(ctx, "Converting MapValue to map")
		dest, diags = attributesToMap(ctx, path, actualValue.Elements(), omitNulls, preserveNulls)
	case types.List:
		tflog.Info(ctx, "Converting ListValue to interface slice")
		dest, diags = elementsToInterfaces(ctx, path, actualValue.Elements(), omitNulls, preserveNulls)
	case types.Set:
		tflog.Info(ctx, "Converting SetValue to interface slice")
		dest, diags = elementsToInterfaces(ctx, path, actualValue.Elements(), omitNulls, preserveNulls)
	case types.Tuple:
		tflog.Info(ctx, "Converting TupleValue to interface slice")
		dest, diags = elementsToInterfaces(ctx, path, actualValue.Elements(), false, preserveNulls)
	case types.Dynamic:
		if actualValue.UnderlyingValue() == nil {
			return nil, nil
		}
		dest, diags = attributeToInterface(ctx, path, actualValue.UnderlyingValue(), false, preserveNulls)
	default:
		diags.AddAttributeError(
			path,
			"Failed to convert Element",
			"Unhandled type: "+actualValue.Type(ctx).String(),
		)
	}

	return
}

func ElementsToInterfaces(ctx context.Context, path path.Path, elements []attr.Value) (dest []interface{}, diags diag.Diagnostics) {
	return elementsToInterfaces(ctx, path, elements, false, false)
}

func elementsToInterfaces(ctx context.Context, path path.Path, elements []attr.Value, omitNulls, preserveNulls bool) (dest []interface{}, diags diag.Diagnostics) {
	dest = make([]interface{}, 0, len(elements))
	for i, element := range elements {
		value, valueDiags := attributeToInterface(ctx, path.AtListIndex(i), element, omitNulls, preserveNulls)
		if valueDiags.HasError() {
			diags.Append(valueDiags...)
		} else {
			dest = append(dest, value)
		}
	}
	return
}

// UnwrapDiscriminatedUnionAtPath turns typed list union items, such as {"file_expiration": {"value": ...}}, into API
// items, such as {"behavior": "file_expiration", "value": ...}. Items that already have the discriminator are API items
// and stay as they are, so a Dynamic value can use either nesting.
func UnwrapDiscriminatedUnionAtPath(ctx context.Context, attributePath path.Path, source any, names []string, discriminator string, variants []JSONSchemaVariant) (any, diag.Diagnostics) {
	return transformDiscriminatedItemsAtPath(attributePath, source, names, func(itemPath path.Path, _ int, item any) (any, diag.Diagnostics) {
		wrapper, ok := item.(map[string]interface{})
		if !ok {
			return jsonSchemaTransformError(itemPath, "expected an object")
		}
		if wrapper[discriminator] != nil {
			return item, nil
		}
		var selected *JSONSchemaVariant
		var inner map[string]interface{}
		for variantIndex := range variants {
			candidate, exists := wrapper[variants[variantIndex].Name]
			if !exists || candidate == nil {
				continue
			}
			if selected != nil {
				return jsonSchemaTransformError(itemPath, "expected exactly one variant")
			}
			selected = &variants[variantIndex]
			inner, ok = candidate.(map[string]interface{})
			if !ok {
				return jsonSchemaTransformError(itemPath.AtName(variants[variantIndex].Name), "expected an object")
			}
		}
		if selected == nil {
			return jsonSchemaTransformError(itemPath, "expected exactly one variant")
		}
		unwrapped := maps.Clone(inner)
		unwrapped[discriminator] = selected.Value
		return unwrapped, nil
	})
}

// WrapDiscriminatedUnionAtPath turns API list union items into typed items under their variant name. Items that are
// already typed, as in state that a Dynamic value stored with the typed nesting, stay as they are.
func WrapDiscriminatedUnionAtPath(ctx context.Context, attributePath path.Path, source any, names []string, discriminator string, variants []JSONSchemaVariant) (any, diag.Diagnostics) {
	return transformDiscriminatedItemsAtPath(attributePath, source, names, func(itemPath path.Path, _ int, item any) (any, diag.Diagnostics) {
		return wrapDiscriminatedItem(itemPath, item, discriminator, variants)
	})
}

// WrapDiscriminatedUnionLikeConfigAtPath wraps an API list union item only where the configured item at the same
// position uses the typed nesting. A Dynamic value read back then keeps the nesting of each configured item.
func WrapDiscriminatedUnionLikeConfigAtPath(ctx context.Context, attributePath path.Path, source any, configured attr.Value, names []string, discriminator string, variants []JSONSchemaVariant) (any, diag.Diagnostics) {
	value, diags := JSONValueToAPI(ctx, attributePath, configured)
	if diags.HasError() {
		return source, diags
	}
	var configuredItems []any
	if _, diags := transformJSONSchemaAtPath(attributePath, value, names, func(_ path.Path, items any) (any, diag.Diagnostics) {
		configuredItems, _ = items.([]any)
		return items, nil
	}); diags.HasError() {
		return source, diags
	}
	return transformDiscriminatedItemsAtPath(attributePath, source, names, func(itemPath path.Path, index int, item any) (any, diag.Diagnostics) {
		if index >= len(configuredItems) || !typedDiscriminatedItem(configuredItems[index], variants) {
			return item, nil
		}
		return wrapDiscriminatedItem(itemPath, item, discriminator, variants)
	})
}

func wrapDiscriminatedItem(itemPath path.Path, item any, discriminator string, variants []JSONSchemaVariant) (any, diag.Diagnostics) {
	if typedDiscriminatedItem(item, variants) {
		return item, nil
	}
	object, ok := item.(map[string]interface{})
	if !ok {
		return jsonSchemaTransformError(itemPath, "expected an object")
	}
	variantValue, ok := object[discriminator].(string)
	if !ok {
		return jsonSchemaTransformError(itemPath.AtName(discriminator), "expected a string discriminator")
	}
	variant, ok := findJSONSchemaVariant(variants, variantValue)
	if !ok {
		return jsonSchemaTransformError(itemPath.AtName(discriminator), "unknown variant "+variantValue)
	}
	wrapped := maps.Clone(object)
	delete(wrapped, discriminator)
	return map[string]interface{}{variant.Name: wrapped}, nil
}

// typedDiscriminatedItem reports whether item is a typed list union item: an object whose keys are all variant names,
// with exactly one of them set. Unused variants may be null, as in a value of an object-typed Terraform variable.
func typedDiscriminatedItem(item any, variants []JSONSchemaVariant) bool {
	object, ok := item.(map[string]interface{})
	if !ok {
		return false
	}
	selected := 0
	for name, value := range object {
		if !slices.ContainsFunc(variants, func(variant JSONSchemaVariant) bool { return variant.Name == name }) {
			return false
		}
		if value != nil {
			selected++
		}
	}
	return selected == 1
}

func transformDiscriminatedItemsAtPath(attributePath path.Path, source any, names []string, transform func(path.Path, int, any) (any, diag.Diagnostics)) (any, diag.Diagnostics) {
	return transformJSONSchemaAtPath(attributePath, source, names, func(valuePath path.Path, value any) (any, diag.Diagnostics) {
		items, ok := value.([]interface{})
		if !ok {
			return jsonSchemaTransformError(valuePath, "expected an array")
		}
		result := make([]interface{}, 0, len(items))
		for index, item := range items {
			converted, diags := transform(valuePath.AtListIndex(index), index, item)
			if diags.HasError() {
				return nil, diags
			}
			result = append(result, converted)
		}
		return result, nil
	})
}

func UnwrapSiblingDiscriminator(_ context.Context, attributePath path.Path, source any, discriminator string, variants []JSONSchemaVariant) (any, diag.Diagnostics) {
	if source == nil {
		return nil, nil
	}
	wrapper, ok := source.(map[string]interface{})
	if !ok {
		return jsonSchemaTransformError(attributePath, "expected an object")
	}
	variant, ok := findJSONSchemaVariant(variants, discriminator)
	if !ok {
		return jsonSchemaTransformError(attributePath, "unknown sibling discriminator "+discriminator)
	}
	var selected string
	var value any
	for name, candidate := range wrapper {
		if candidate == nil {
			continue
		}
		if selected != "" {
			return jsonSchemaTransformError(attributePath, "expected exactly one variant")
		}
		selected = name
		value = candidate
	}
	if selected != variant.Name {
		return jsonSchemaTransformError(attributePath, "value variant must match sibling discriminator "+discriminator)
	}
	return value, nil
}

// NullRemovedSiblingFields sends null for each writable field that prior, the typed value in state, sets and value, the
// unwrapped request value, no longer sets. Rails merges a Behavior value on update and removes only fields sent as null.
func NullRemovedSiblingFields(ctx context.Context, attributePath path.Path, value any, prior attr.Value, discriminator string, variants []JSONSchemaVariant) (any, diag.Diagnostics) {
	if prior == nil || prior.IsNull() || prior.IsUnknown() {
		return value, nil
	}
	priorValue, diags := SchemaAttributeToInterface(ctx, attributePath, prior)
	if diags.HasError() {
		return value, diags
	}
	// A prior value under another behavior name has nothing to remove from this one.
	priorValue, unwrapDiags := UnwrapSiblingDiscriminator(ctx, attributePath, priorValue, discriminator, variants)
	if unwrapDiags.HasError() {
		return value, nil
	}
	variant, _ := findJSONSchemaVariant(variants, discriminator)
	return nullRemovedWritableFields(value, priorValue, variant), nil
}

func nullRemovedWritableFields(value, prior any, variant JSONSchemaVariant) any {
	valueObject, valueIsObject := value.(map[string]interface{})
	priorObject, priorIsObject := prior.(map[string]interface{})
	if !valueIsObject || !priorIsObject {
		return value
	}
	result := maps.Clone(valueObject)
	for key := range priorObject {
		if _, exists := result[key]; !exists && slices.Contains(variant.Writable, key) {
			result[key] = nil
		}
	}
	return result
}

func WrapSiblingDiscriminator(_ context.Context, attributePath path.Path, source any, discriminator string, variants []JSONSchemaVariant) (any, diag.Diagnostics) {
	if source == nil {
		return nil, nil
	}
	source, _ = decodeJSONString(source)
	variant, ok := findJSONSchemaVariant(variants, discriminator)
	if !ok {
		return jsonSchemaTransformError(attributePath, "unknown sibling discriminator "+discriminator)
	}
	legacy, shapeDiags := siblingDiscriminatorValueIsLegacy(attributePath, source, discriminator, variants)
	if shapeDiags.HasError() {
		return nil, shapeDiags
	}
	if !legacy {
		return source, nil
	}
	if variant.LegacyProperty != "" {
		if value, ok := source.(string); ok && value == "" {
			return nil, nil
		}
		if _, ok := source.(map[string]interface{}); !ok {
			if variant.LegacyList {
				if _, ok := source.([]interface{}); !ok {
					source = []interface{}{source}
				}
			}
			source = map[string]interface{}{variant.LegacyProperty: source}
		}
	}
	return map[string]interface{}{variant.Name: source}, nil
}

func UngroupStructuralUnionAtPath(ctx context.Context, attributePath path.Path, source any, names []string, variants []JSONSchemaVariant) (any, diag.Diagnostics) {
	return transformJSONSchemaAtPath(attributePath, source, names, func(valuePath path.Path, value any) (any, diag.Diagnostics) {
		entries, ok := value.(map[string]interface{})
		if !ok {
			return jsonSchemaTransformError(valuePath, "expected an object")
		}
		result := make(map[string]interface{}, len(entries))
		for key, value := range entries {
			entryPath := valuePath.AtMapKey(key)
			groups, ok := value.(map[string]interface{})
			if !ok {
				return jsonSchemaTransformError(entryPath, "expected an object")
			}
			items := []interface{}{}
			for _, variant := range variants {
				group, exists := groups[variant.Name]
				if !exists || group == nil {
					continue
				}
				groupItems, ok := group.([]interface{})
				if !ok {
					return jsonSchemaTransformError(entryPath.AtName(variant.Name), "expected an array")
				}
				items = append(items, groupItems...)
			}
			result[key] = items
		}
		return result, nil
	})
}

func GroupStructuralUnionAtPath(ctx context.Context, attributePath path.Path, source any, names []string, variants []JSONSchemaVariant) (any, diag.Diagnostics) {
	return transformJSONSchemaAtPath(attributePath, source, names, func(valuePath path.Path, value any) (any, diag.Diagnostics) {
		entries, ok := value.(map[string]interface{})
		if !ok {
			return jsonSchemaTransformError(valuePath, "expected an object")
		}
		result := make(map[string]interface{}, len(entries))
		for key, value := range entries {
			entryPath := valuePath.AtMapKey(key)
			items, ok := value.([]interface{})
			if !ok {
				return jsonSchemaTransformError(entryPath, "expected an array")
			}
			groups := map[string]interface{}{}
			for index, item := range items {
				itemPath := entryPath.AtListIndex(index)
				object, ok := item.(map[string]interface{})
				if !ok {
					return jsonSchemaTransformError(itemPath, "expected an object")
				}
				matches := []JSONSchemaVariant{}
				for _, variant := range variants {
					if jsonSchemaObjectMatches(object, variant) {
						matches = append(matches, variant)
					}
				}
				if len(matches) != 1 {
					return jsonSchemaTransformError(itemPath, "expected exactly one matching variant")
				}
				group, _ := groups[matches[0].Name].([]interface{})
				groups[matches[0].Name] = append(group, object)
			}
			result[key] = groups
		}
		return result, nil
	})
}

// WrapLegacySingleItemAtPath reads one stored item at a list marked x-files-legacy-single-item as a one-item list, so
// 5 becomes [5]. Lists and null stay as they are. The item conversion that follows still turns "5" into 5 and rejects
// values that are not valid items.
func WrapLegacySingleItemAtPath(_ context.Context, attributePath path.Path, source any, names []string) (any, diag.Diagnostics) {
	return transformJSONSchemaAtPath(attributePath, source, names, func(_ path.Path, value any) (any, diag.Diagnostics) {
		switch value.(type) {
		case []interface{}, []map[string]interface{}:
			return value, nil
		}
		return []interface{}{value}, nil
	})
}

// Read marked boolean fields using Ruby truthiness, rather than parsing strings such as "false".
func ReadLegacyTruthyAtPath(_ context.Context, attributePath path.Path, source any, names []string) (any, diag.Diagnostics) {
	return transformJSONSchemaAtPath(attributePath, source, names, func(_ path.Path, value any) (any, diag.Diagnostics) {
		if _, boolean := value.(bool); boolean {
			return value, nil
		}
		return true, nil
	})
}

func transformJSONSchemaAtPath(attributePath path.Path, source any, names []string, transform func(path.Path, any) (any, diag.Diagnostics)) (any, diag.Diagnostics) {
	if source == nil {
		return nil, nil
	}
	if len(names) == 0 {
		return transform(attributePath, source)
	}
	object, ok := source.(map[string]interface{})
	if !ok {
		return jsonSchemaTransformError(attributePath, "expected an object")
	}
	value, exists := object[names[0]]
	if !exists || value == nil {
		return source, nil
	}
	converted, diags := transformJSONSchemaAtPath(attributePath.AtName(names[0]), value, names[1:], transform)
	if diags.HasError() {
		return source, diags
	}
	result := maps.Clone(object)
	result[names[0]] = converted
	return result, nil
}

func jsonSchemaObjectMatches(object map[string]interface{}, variant JSONSchemaVariant) bool {
	for _, required := range variant.Required {
		if _, ok := object[required]; !ok {
			return false
		}
	}
	for key := range object {
		if !slices.Contains(variant.Allowed, key) {
			return false
		}
	}
	return true
}

func findJSONSchemaVariant(variants []JSONSchemaVariant, value string) (JSONSchemaVariant, bool) {
	for _, variant := range variants {
		if variant.Value == value {
			return variant, true
		}
	}
	return JSONSchemaVariant{}, false
}

func jsonSchemaTransformError(attributePath path.Path, detail string) (any, diag.Diagnostics) {
	return nil, diag.Diagnostics{diag.NewAttributeErrorDiagnostic(attributePath, "Failed to convert JSON Schema value", detail)}
}

func ToObject(ctx context.Context, path path.Path, source any, plan types.Object) (types.Object, diag.Diagnostics) {
	value, diags := schemaValue(ctx, path, source, types.ObjectType{AttrTypes: plan.AttributeTypes(ctx)})
	if diags.HasError() {
		return types.ObjectNull(plan.AttributeTypes(ctx)), diags
	}
	return value.(types.Object), nil
}

func ToList(ctx context.Context, path path.Path, source any, plan types.List) (types.List, diag.Diagnostics) {
	value, diags := schemaValue(ctx, path, source, types.ListType{ElemType: plan.ElementType(ctx)})
	if diags.HasError() {
		return types.ListNull(plan.ElementType(ctx)), diags
	}
	return value.(types.List), nil
}

func ToMap(ctx context.Context, path path.Path, source any, plan types.Map) (types.Map, diag.Diagnostics) {
	value, diags := schemaValue(ctx, path, source, types.MapType{ElemType: plan.ElementType(ctx)})
	if diags.HasError() {
		return types.MapNull(plan.ElementType(ctx)), diags
	}
	return value.(types.Map), nil
}

func schemaValue(ctx context.Context, path path.Path, source any, target attr.Type) (attr.Value, diag.Diagnostics) {
	if source == nil {
		return nullSchemaValue(target), nil
	}
	switch source.(type) {
	case map[string]any, []any, []map[string]any, string, bool, int, int64, float64:
	default:
		generic, err := decodeSDKValue(source)
		if err != nil {
			return schemaConversionError(path, target, source)
		}
		return schemaValue(ctx, path, generic, target)
	}

	switch targetType := target.(type) {
	case basetypes.ObjectType:
		sourceMap, ok := source.(map[string]any)
		if text, isString := source.(string); isString {
			// Older records and Dynamic state may hold an object as its JSON text, such as a key file.
			ok = json.Unmarshal([]byte(text), &sourceMap) == nil && sourceMap != nil
		}
		if !ok {
			return schemaConversionError(path, target, source)
		}
		values := make(map[string]attr.Value, len(targetType.AttrTypes))
		var diags diag.Diagnostics
		for name, attributeType := range targetType.AttrTypes {
			value, valueDiags := schemaValue(ctx, path.AtName(name), sourceMap[name], attributeType)
			diags.Append(valueDiags...)
			values[name] = value
		}
		if diags.HasError() {
			return nil, diags
		}
		value, valueDiags := types.ObjectValue(targetType.AttrTypes, values)
		return value, valueDiags
	case basetypes.ListType:
		sourceList, ok := source.([]any)
		if objects, objectsOK := source.([]map[string]any); objectsOK {
			sourceList = make([]any, len(objects))
			for index, object := range objects {
				sourceList[index] = object
			}
			ok = true
		}
		if !ok {
			return schemaConversionError(path, target, source)
		}
		values := make([]attr.Value, 0, len(sourceList))
		var diags diag.Diagnostics
		for index, item := range sourceList {
			value, valueDiags := schemaValue(ctx, path.AtListIndex(index), item, targetType.ElemType)
			diags.Append(valueDiags...)
			values = append(values, value)
		}
		if diags.HasError() {
			return nil, diags
		}
		value, valueDiags := types.ListValue(targetType.ElemType, values)
		return value, valueDiags
	case basetypes.MapType:
		sourceMap, ok := source.(map[string]any)
		if !ok {
			return schemaConversionError(path, target, source)
		}
		values := make(map[string]attr.Value, len(sourceMap))
		var diags diag.Diagnostics
		for name, item := range sourceMap {
			value, valueDiags := schemaValue(ctx, path.AtMapKey(name), item, targetType.ElemType)
			diags.Append(valueDiags...)
			values[name] = value
		}
		if diags.HasError() {
			return nil, diags
		}
		value, valueDiags := types.MapValue(targetType.ElemType, values)
		return value, valueDiags
	case basetypes.StringType:
		value, ok := source.(string)
		if !ok {
			return schemaConversionError(path, target, source)
		}
		return types.StringValue(value), nil
	case basetypes.BoolType:
		switch value := source.(type) {
		case bool:
			return types.BoolValue(value), nil
		case string:
			if value == "true" || value == "false" {
				return types.BoolValue(value == "true"), nil
			}
		}
		return schemaConversionError(path, target, source)
	case basetypes.Int64Type:
		switch value := source.(type) {
		case int:
			return types.Int64Value(int64(value)), nil
		case int64:
			return types.Int64Value(value), nil
		case float64:
			if math.Trunc(value) == value {
				return types.Int64Value(int64(value)), nil
			}
		case string:
			value = strings.TrimSpace(value)
			if value == "" {
				return types.Int64Null(), nil
			}
			if integer, err := strconv.ParseInt(value, 10, 64); err == nil {
				return types.Int64Value(integer), nil
			}
		}
		return schemaConversionError(path, target, source)
	case basetypes.Float64Type:
		switch value := source.(type) {
		case int:
			return types.Float64Value(float64(value)), nil
		case int64:
			return types.Float64Value(float64(value)), nil
		case float64:
			return types.Float64Value(value), nil
		}
		return schemaConversionError(path, target, source)
	case basetypes.DynamicType:
		return ToDynamic(ctx, path, source, types.DynamicNull())
	default:
		return schemaConversionError(path, target, source)
	}
}

func nullSchemaValue(target attr.Type) attr.Value {
	switch targetType := target.(type) {
	case basetypes.ObjectType:
		return types.ObjectNull(targetType.AttrTypes)
	case basetypes.ListType:
		return types.ListNull(targetType.ElemType)
	case basetypes.MapType:
		return types.MapNull(targetType.ElemType)
	case basetypes.StringType:
		return types.StringNull()
	case basetypes.BoolType:
		return types.BoolNull()
	case basetypes.Int64Type:
		return types.Int64Null()
	case basetypes.Float64Type:
		return types.Float64Null()
	case basetypes.DynamicType:
		return types.DynamicNull()
	}
	return nil
}

func schemaConversionError(path path.Path, target attr.Type, source any) (attr.Value, diag.Diagnostics) {
	return nil, diag.Diagnostics{diag.NewAttributeErrorDiagnostic(path, "Failed to convert API value", fmt.Sprintf("Expected %s, got %T", target.String(), source))}
}

func ToDynamic(ctx context.Context, path path.Path, source any, plan attr.Value) (dest types.Dynamic, diags diag.Diagnostics) {
	ctx = setAttributePath(ctx, path)

	switch actualValue := source.(type) {
	case map[string]interface{}:
		tflog.Info(ctx, "Converting map to ObjectValue")
		elementTypes := map[string]attr.Type{}
		elements := map[string]attr.Value{}

		for key, value := range actualValue {
			var planAttr attr.Value = nil
			planSchema, ok := plan.(types.Object)
			if ok {
				planAttr = planSchema.Attributes()[key]
				if planAttr == nil {
					tflog.Info(ctx, "Skipping unknown attribute: "+key)
					continue
				}
			}

			attrValue, attrDiags := ToDynamic(ctx, path.AtMapKey(key), value, planAttr)
			if attrDiags.HasError() {
				diags.Append(attrDiags...)
			} else {
				elementTypes[key] = attrValue.Type(ctx)
				elements[key] = attrValue
			}
		}

		if !diags.HasError() {
			objectValue, diags := types.ObjectValue(elementTypes, elements)
			if !diags.HasError() {
				dest = types.DynamicValue(objectValue)
			}
		}
	case []any, []map[string]interface{}:
		tflog.Info(ctx, "Converting slice to TupleValue")
		var anySlice []any
		if mapSlice, ok := actualValue.([]map[string]interface{}); ok {
			anySlice = make([]any, 0, len(mapSlice))
			for _, element := range mapSlice {
				anySlice = append(anySlice, element)
			}
		} else {
			anySlice = actualValue.([]any)
		}

		elementTypes := make([]attr.Type, 0, len(anySlice))
		elements := make([]attr.Value, 0, len(anySlice))

		for i, element := range anySlice {
			dynamic, attrDiags := ToDynamic(ctx, path.AtListIndex(i), element, nil)
			if attrDiags.HasError() {
				diags.Append(attrDiags...)
			} else {
				elementTypes = append(elementTypes, types.DynamicType)
				elements = append(elements, dynamic)
			}
		}

		if !diags.HasError() {
			tupleValue, diags := types.TupleValue(elementTypes, elements)
			if !diags.HasError() {
				dest = types.DynamicValue(tupleValue)
			}
		}
	case bool:
		tflog.Info(ctx, "Converting bool to BoolValue")
		dest = types.DynamicValue(types.BoolValue(actualValue))
	case string:
		tflog.Info(ctx, "Converting string to StringValue")
		dest = types.DynamicValue(types.StringValue(actualValue))
	case float64:
		tflog.Info(ctx, "Converting float64 to Float64Value")
		if _, ok := plan.(types.Number); ok {
			dest = types.DynamicValue(types.NumberValue(big.NewFloat(actualValue)))
		} else {
			dest = types.DynamicValue(types.Float64Value(actualValue))
		}
	case int64:
		tflog.Info(ctx, "Converting int64 to Int64Value")
		if _, ok := plan.(types.Number); ok {
			dest = types.DynamicValue(types.NumberValue(new(big.Float).SetInt64(actualValue)))
		} else {
			dest = types.DynamicValue(types.Int64Value(actualValue))
		}
	case nil:
		tflog.Info(ctx, "Skipping nil value")
	default:
		generic, err := decodeSDKValue(source)
		if err != nil {
			diags.AddError("Failed to convert Element", "Unhandled type for "+path.String()+": "+err.Error())
			return
		}
		if generic != nil {
			return ToDynamic(ctx, path, generic, plan)
		}
	}

	return
}

// decodeSDKValue turns Go SDK objects and typed lists, such as []files_sdk.FormField, into generic JSON values.
func decodeSDKValue(source any) (any, error) {
	if source == nil {
		return nil, nil
	}
	if kind := reflect.TypeOf(source).Kind(); kind != reflect.Struct && kind != reflect.Slice {
		return nil, fmt.Errorf("%T", source)
	}
	encoded, err := json.Marshal(source)
	if err != nil {
		return nil, fmt.Errorf("%T: %w", source, err)
	}
	var generic any
	if err := json.Unmarshal(encoded, &generic); err != nil {
		return nil, fmt.Errorf("%T: %w", source, err)
	}
	switch generic.(type) {
	case nil, map[string]any, []any:
		return generic, nil
	}
	return nil, fmt.Errorf("%T is not a JSON object or array", source)
}

func ListValueToString(ctx context.Context, path path.Path, list types.List, delim string) (string, diag.Diagnostics) {
	if list.IsNull() || list.IsUnknown() {
		return "", nil
	}

	ctx = setAttributePath(ctx, path)
	length := len(list.Elements())
	strs := make([]string, 0, length)

	switch list.ElementType(ctx) {
	case types.StringType:
		tflog.Info(ctx, "Converting StringType List to string slice")
		diags := list.ElementsAs(ctx, &strs, false)
		if diags.HasError() {
			return "", diags
		}
	case types.Int64Type:
		tflog.Info(ctx, "Converting Int64Type List to string slice")
		elements := make([]int64, 0, length)
		diags := list.ElementsAs(ctx, &elements, false)
		if diags.HasError() {
			return "", diags
		}

		for _, element := range elements {
			strs = append(strs, fmt.Sprint(element))
		}
	default:
		return "", diag.Diagnostics{
			diag.NewAttributeErrorDiagnostic(
				path,
				"Failed to convert List elements",
				"Unhandled type: "+list.ElementType(ctx).String(),
			),
		}
	}

	return strings.Join(strs, delim), nil
}

func setAttributePath(ctx context.Context, path path.Path) context.Context {
	return tflog.SetField(ctx, "attribute", path.String())
}
