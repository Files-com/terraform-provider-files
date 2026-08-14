package lib

import (
	"context"
	"fmt"
	"maps"
	"math"
	"slices"
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
	Name     string
	Value    string
	Required []string
	Allowed  []string
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
	return attributesToMap(ctx, path, attrs, false)
}

func attributesToMap(ctx context.Context, path path.Path, attrs map[string]attr.Value, omitNulls bool) (dest map[string]interface{}, diags diag.Diagnostics) {
	dest = make(map[string]interface{})
	for key, value := range attrs {
		if omitNulls && (value.IsNull() || value.IsUnknown()) {
			continue
		}
		attrValue, attrDiags := attributeToInterface(ctx, path.AtMapKey(key), value, omitNulls)
		if attrDiags.HasError() {
			diags.Append(attrDiags...)
		} else {
			dest[key] = attrValue
		}
	}

	return
}

func AttributeToInterface(ctx context.Context, path path.Path, source attr.Value) (dest interface{}, diags diag.Diagnostics) {
	return attributeToInterface(ctx, path, source, false)
}

func SchemaAttributeToInterface(ctx context.Context, path path.Path, source attr.Value) (dest interface{}, diags diag.Diagnostics) {
	return attributeToInterface(ctx, path, source, true)
}

func attributeToInterface(ctx context.Context, path path.Path, source attr.Value, omitNulls bool) (dest interface{}, diags diag.Diagnostics) {
	if omitNulls && (source.IsNull() || source.IsUnknown()) {
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
		dest, diags = attributesToMap(ctx, path, actualValue.Attributes(), omitNulls)
	case types.Map:
		tflog.Info(ctx, "Converting MapValue to map")
		dest, diags = attributesToMap(ctx, path, actualValue.Elements(), omitNulls)
	case types.List:
		tflog.Info(ctx, "Converting ListValue to interface slice")
		dest, diags = elementsToInterfaces(ctx, path, actualValue.Elements(), omitNulls)
	case types.Tuple:
		tflog.Info(ctx, "Converting TupleValue to interface slice")
		dest, diags = elementsToInterfaces(ctx, path, actualValue.Elements(), false)
	case types.Dynamic:
		dest, diags = attributeToInterface(ctx, path, actualValue.UnderlyingValue(), false)
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
	return elementsToInterfaces(ctx, path, elements, false)
}

func elementsToInterfaces(ctx context.Context, path path.Path, elements []attr.Value, omitNulls bool) (dest []interface{}, diags diag.Diagnostics) {
	dest = make([]interface{}, 0, len(elements))
	for i, element := range elements {
		value, valueDiags := attributeToInterface(ctx, path.AtListIndex(i), element, omitNulls)
		if valueDiags.HasError() {
			diags.Append(valueDiags...)
		} else {
			dest = append(dest, value)
		}
	}
	return
}

func UnwrapDiscriminatedUnionAtPath(ctx context.Context, attributePath path.Path, source any, names []string, discriminator string, variants []JSONSchemaVariant) (any, diag.Diagnostics) {
	return transformJSONSchemaAtPath(attributePath, source, names, func(valuePath path.Path, value any) (any, diag.Diagnostics) {
		items, ok := value.([]interface{})
		if !ok {
			return jsonSchemaTransformError(valuePath, "expected an array")
		}
		result := make([]interface{}, 0, len(items))
		for index, item := range items {
			itemPath := valuePath.AtListIndex(index)
			wrapper, ok := item.(map[string]interface{})
			if !ok {
				return jsonSchemaTransformError(itemPath, "expected an object")
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
			result = append(result, unwrapped)
		}
		return result, nil
	})
}

func WrapDiscriminatedUnionAtPath(ctx context.Context, attributePath path.Path, source any, names []string, discriminator string, variants []JSONSchemaVariant) (any, diag.Diagnostics) {
	return transformJSONSchemaAtPath(attributePath, source, names, func(valuePath path.Path, value any) (any, diag.Diagnostics) {
		items, ok := value.([]interface{})
		if !ok {
			return jsonSchemaTransformError(valuePath, "expected an array")
		}
		result := make([]interface{}, 0, len(items))
		for index, item := range items {
			itemPath := valuePath.AtListIndex(index)
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
			result = append(result, map[string]interface{}{variant.Name: wrapped})
		}
		return result, nil
	})
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

	switch targetType := target.(type) {
	case basetypes.ObjectType:
		sourceMap, ok := source.(map[string]any)
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
		value, ok := source.(bool)
		if !ok {
			return schemaConversionError(path, target, source)
		}
		return types.BoolValue(value), nil
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
		dest = types.DynamicValue(types.Float64Value(actualValue))
	case nil:
		tflog.Info(ctx, "Skipping nil value")
	default:
		diags.AddError(
			"Failed to convert Element",
			"Unhandled type for "+path.String()+": "+fmt.Sprintf("%T", source),
		)
	}

	return
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
