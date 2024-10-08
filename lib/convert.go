package lib

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

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

func ElementsToStringMap(ctx context.Context, path path.Path, attrs map[string]attr.Value) (dest map[string]interface{}, diags diag.Diagnostics) {
	dest = make(map[string]interface{})

	for key, value := range attrs {
		attrValue, attrDiags := AttributeToInterface(ctx, path.AtMapKey(key), value)
		if attrDiags.HasError() {
			diags.Append(attrDiags...)
		} else {
			dest[key] = attrValue
		}
	}

	return
}

func AttributeToInterface(ctx context.Context, path path.Path, source attr.Value) (dest interface{}, diags diag.Diagnostics) {
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
	case types.Object:
		tflog.Info(ctx, "Converting ObjectValue to map")
		dest, diags = ElementsToStringMap(ctx, path, actualValue.Attributes())
	case types.Tuple:
		tflog.Info(ctx, "Converting TupleValue to interface slice")
		dest = []interface{}{}

		for i, element := range actualValue.Elements() {
			attrValue, attrDiags := AttributeToInterface(ctx, path.AtListIndex(i), element)
			if attrDiags.HasError() {
				diags.Append(attrDiags...)
			} else {
				dest = append(dest.([]interface{}), attrValue)
			}
		}
	default:
		diags.AddAttributeError(
			path,
			"Failed to convert Element",
			"Unhandled type: "+actualValue.Type(ctx).String(),
		)
	}

	return
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
