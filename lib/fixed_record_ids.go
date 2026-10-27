package lib

import (
	"context"
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type fixedRecordIDs struct {
	key          string
	serverFilled []string
}

// FixedRecordIDs plans records that the API keeps under an ID, such as form fields. Rails keeps a record whose request
// names a current ID: it keeps every stored value that the request leaves out, and it stores a changed value under a
// new ID. key is the ID field. serverFilled fields get an API value for a new record too, such as a default.
//
// For a typed list, a record that names a stored ID plans each left-out value as the stored value, and a new record
// plans each left-out value as null, except serverFilled fields. A configured value that differs from the stored record
// is an error, because the API would replace the ID. An unknown value is checked when Terraform plans again with the
// final values, before apply. A Dynamic value holds only the configured keys, so it is only checked.
func FixedRecordIDs(key string, serverFilled []string) fixedRecordIDs {
	return fixedRecordIDs{key: key, serverFilled: serverFilled}
}

func (m fixedRecordIDs) Description(ctx context.Context) string {
	return m.MarkdownDescription(ctx)
}

func (m fixedRecordIDs) MarkdownDescription(_ context.Context) string {
	return "A record that sets " + m.key + " keeps the stored values it leaves out and must not change the others."
}

func (m fixedRecordIDs) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() || resp.PlanValue.IsNull() || resp.PlanValue.IsUnknown() {
		return
	}
	stored := map[string]types.Object{}
	if !req.StateValue.IsNull() && !req.StateValue.IsUnknown() {
		for _, element := range req.StateValue.Elements() {
			if record, ok := element.(types.Object); ok && knownValue(ctx, record.Attributes()[m.key]) {
				stored[record.Attributes()[m.key].String()] = record
			}
		}
	}
	configured, planned := req.ConfigValue.Elements(), resp.PlanValue.Elements()
	if len(configured) != len(planned) {
		return
	}
	elements := make([]attr.Value, len(planned))
	for index, element := range planned {
		elements[index] = element
		config, configOK := configured[index].(types.Object)
		plan, planOK := element.(types.Object)
		if !configOK || !planOK || config.IsNull() || config.IsUnknown() || plan.IsNull() || plan.IsUnknown() {
			continue
		}
		attributes := plan.Attributes()
		id := config.Attributes()[m.key]
		previous, kept := types.Object{}, false
		if knownValue(ctx, id) && !id.IsNull() {
			previous, kept = stored[id.String()]
		}
		var changed []string
		for name, value := range config.Attributes() {
			if name == m.key {
				continue
			}
			switch {
			case !knownValue(ctx, value):
				// Unknown until apply. Terraform plans again with the final value before applying.
			case kept && value.IsNull():
				attributes[name] = previous.Attributes()[name]
			case kept && !value.Equal(previous.Attributes()[name]):
				changed = append(changed, name)
			case !kept && id.IsNull() && value.IsNull() && !slices.Contains(m.serverFilled, name):
				attributes[name] = nullValue(ctx, attributes[name])
			}
		}
		if len(changed) > 0 {
			resp.Diagnostics.Append(m.changedError(req.Path.AtListIndex(index), id.String(), changed))
			continue
		}
		object, diags := types.ObjectValue(plan.AttributeTypes(ctx), attributes)
		resp.Diagnostics.Append(diags...)
		elements[index] = object
	}
	if resp.Diagnostics.HasError() {
		return
	}
	list, diags := types.ListValue(resp.PlanValue.ElementType(ctx), elements)
	resp.Diagnostics.Append(diags...)
	resp.PlanValue = list
}

// PlanModifyDynamic checks a Dynamic list of records. State holds only the keys configured before, and the API casts
// boolean and number strings, so only text values and lists of text that both sides hold are compared.
func (m fixedRecordIDs) PlanModifyDynamic(ctx context.Context, req planmodifier.DynamicRequest, resp *planmodifier.DynamicResponse) {
	configured := dynamicRecords(ctx, req.Path, req.ConfigValue)
	stored := map[string]map[string]any{}
	for _, record := range dynamicRecords(ctx, req.Path, req.StateValue) {
		if id, ok := record[m.key]; ok && id != nil && id != unknownField {
			stored[fmt.Sprint(id)] = record
		}
	}
	for index, record := range configured {
		id, ok := record[m.key]
		if !ok || id == nil || id == unknownField {
			continue
		}
		previous, kept := stored[fmt.Sprint(id)]
		if !kept {
			continue
		}
		var changed []string
		for name, value := range record {
			if name != m.key && value != unknownField && textValue(value) && textValue(previous[name]) && !reflect.DeepEqual(value, previous[name]) {
				changed = append(changed, name)
			}
		}
		if len(changed) > 0 {
			resp.Diagnostics.Append(m.changedError(req.Path.AtListIndex(index), fmt.Sprint(id), changed))
		}
	}
}

func (m fixedRecordIDs) changedError(recordPath path.Path, id string, changed []string) diag.Diagnostic {
	slices.Sort(changed)
	return diag.NewAttributeErrorDiagnostic(recordPath.AtName(m.key), "Record keeps its "+m.key+" but changes its values",
		fmt.Sprintf("This record sets %s = %s and changes %s. Files.com keeps a record under its %s only while its values stay the same, and stores a changed record under a new %s. Remove %s from this record to let Files.com replace it, or keep its stored values. A left-out value keeps the stored one.",
			m.key, id, strings.Join(changed, ", "), m.key, m.key, m.key))
}

type unknownMarker struct{}

// unknownField stands for a configured field that is unknown until apply.
var unknownField any = unknownMarker{}

// dynamicRecords returns the objects of a Dynamic list as Go values. A field that is not fully known is unknownField.
func dynamicRecords(ctx context.Context, attributePath path.Path, value attr.Value) []map[string]any {
	dynamic, ok := value.(types.Dynamic)
	if !ok || dynamic.IsNull() || dynamic.IsUnknown() || dynamic.IsUnderlyingValueNull() || dynamic.IsUnderlyingValueUnknown() {
		return nil
	}
	var elements []attr.Value
	switch underlying := dynamic.UnderlyingValue().(type) {
	case types.String:
		decoded, diags := JSONValueToAPI(ctx, attributePath, dynamic)
		if diags.HasError() {
			return nil
		}
		var records []map[string]any
		for _, record := range asList(decoded) {
			if object, ok := record.(map[string]any); ok {
				records = append(records, object)
			}
		}
		return records
	case types.Tuple:
		elements = underlying.Elements()
	case types.List:
		elements = underlying.Elements()
	default:
		return nil
	}
	var records []map[string]any
	for _, element := range elements {
		object, ok := element.(types.Object)
		if !ok || object.IsNull() || object.IsUnknown() {
			continue
		}
		record := map[string]any{}
		for name, field := range object.Attributes() {
			if !knownValue(ctx, field) {
				record[name] = unknownField
				continue
			}
			converted, _ := attributeToInterface(ctx, attributePath, field, false, true)
			record[name] = converted
		}
		records = append(records, record)
	}
	return records
}

func knownValue(ctx context.Context, value attr.Value) bool {
	if value == nil {
		return false
	}
	terraformValue, err := value.ToTerraformValue(ctx)
	return err == nil && terraformValue.IsFullyKnown()
}

func nullValue(ctx context.Context, value attr.Value) attr.Value {
	valueType := value.Type(ctx)
	null, err := valueType.ValueFromTerraform(ctx, tftypes.NewValue(valueType.TerraformType(ctx), nil))
	if err != nil {
		return value
	}
	return null
}

// textValue reports whether value is text or a list of text, which the API stores as sent.
func textValue(value any) bool {
	switch value := value.(type) {
	case string:
		return true
	case []any:
		for _, item := range value {
			if _, ok := item.(string); !ok {
				return false
			}
		}
		return true
	}
	return false
}

func asList(value any) []any {
	list, _ := value.([]any)
	return list
}
