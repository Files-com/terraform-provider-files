package lib

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Rails rejects an empty members list on create and ignores it on update, so a configured [] would never be read back.
func TestRejectEmptyList(t *testing.T) {
	ctx := context.Background()
	memberType := types.ObjectType{AttrTypes: map[string]attr.Type{"email": types.StringType}}
	listErrors := func(value types.List) int {
		resp := &validator.ListResponse{}
		RejectEmptyList("detail").ValidateList(ctx, validator.ListRequest{Path: path.Root("members"), ConfigValue: value}, resp)
		return resp.Diagnostics.ErrorsCount()
	}
	assert.Equal(t, 1, listErrors(types.ListValueMust(memberType, []attr.Value{})))
	assert.Zero(t, listErrors(types.ListValueMust(memberType, []attr.Value{types.ObjectValueMust(memberType.AttrTypes, map[string]attr.Value{"email": types.StringValue("a@example.com")})})))
	assert.Zero(t, listErrors(types.ListNull(memberType)))
	assert.Zero(t, listErrors(types.ListUnknown(memberType)))

	dynamicErrors := func(value types.Dynamic) int {
		resp := &validator.DynamicResponse{}
		RejectEmptyList("detail").ValidateDynamic(ctx, validator.DynamicRequest{Path: path.Root("members"), ConfigValue: value}, resp)
		return resp.Diagnostics.ErrorsCount()
	}
	assert.Equal(t, 1, dynamicErrors(types.DynamicValue(types.TupleValueMust([]attr.Type{}, []attr.Value{}))))
	assert.Equal(t, 1, dynamicErrors(types.DynamicValue(types.StringValue("[]"))))
	assert.Zero(t, dynamicErrors(types.DynamicValue(types.StringValue(`[{"email":"a@example.com"}]`))))
	assert.Zero(t, dynamicErrors(types.DynamicNull()))
}

// Rails removes keys that are not policy settings when it saves a policy.
func TestKnownKeys(t *testing.T) {
	ctx := context.Background()
	errorsFor := func(value types.Dynamic) []string {
		resp := &validator.DynamicResponse{}
		KnownKeys("Files.com removes them.", "welcome_email_enabled", "folder_behaviors").ValidateDynamic(ctx, validator.DynamicRequest{Path: path.Root("value"), ConfigValue: value}, resp)
		var details []string
		for _, d := range resp.Diagnostics.Errors() {
			details = append(details, d.Detail())
		}
		return details
	}
	assert.Equal(t, []string{"value sets not_a_setting, which is not documented. Files.com removes them."}, errorsFor(types.DynamicValue(types.StringValue(`{"welcome_email_enabled":true,"not_a_setting":1}`))))
	assert.Equal(t, []string{"value sets a, b, which are not documented. Files.com removes them."}, errorsFor(types.DynamicValue(types.StringValue(`{"b":1,"a":2}`))))
	assert.Empty(t, errorsFor(types.DynamicValue(types.StringValue(`{"welcome_email_enabled":true}`))))
	assert.Equal(t, []string{"value must be an object of documented fields."}, errorsFor(types.DynamicValue(types.StringValue(`[1]`))))
	assert.Empty(t, errorsFor(types.DynamicUnknown()))
	assert.Empty(t, errorsFor(types.DynamicValue(types.StringValue("null"))), "Rails stores a JSON null policy value as {}")
}

// A typed Behavior value must be set under the name in behavior, or the request cannot be built at apply.
func TestSiblingDiscriminatorMatches(t *testing.T) {
	ctx := context.Background()
	valueType := map[string]attr.Type{"webhook": types.ObjectType{AttrTypes: map[string]attr.Type{"urls": types.ListType{ElemType: types.StringType}}}, "storage_region": types.StringType}
	schema := resourceschema.Schema{Attributes: map[string]resourceschema.Attribute{
		"behavior": resourceschema.StringAttribute{Required: true},
		"value": resourceschema.SingleNestedAttribute{Optional: true, Attributes: map[string]resourceschema.Attribute{
			"webhook":        resourceschema.SingleNestedAttribute{Optional: true, Attributes: map[string]resourceschema.Attribute{"urls": resourceschema.ListAttribute{Optional: true, ElementType: types.StringType}}},
			"storage_region": resourceschema.StringAttribute{Optional: true},
		}},
	}}
	webhook := types.ObjectValueMust(valueType["webhook"].(types.ObjectType).AttrTypes, map[string]attr.Value{"urls": types.ListValueMust(types.StringType, []attr.Value{types.StringValue("https://example.com")})})
	value := types.ObjectValueMust(valueType, map[string]attr.Value{"webhook": webhook, "storage_region": types.StringNull()})
	check := func(behavior string) int {
		raw, err := value.ToTerraformValue(ctx)
		require.NoError(t, err)
		config := tfsdk.Config{Schema: schema, Raw: tftypes.NewValue(schema.Type().TerraformType(ctx), map[string]tftypes.Value{"behavior": tftypes.NewValue(tftypes.String, behavior), "value": raw})}
		resp := &validator.ObjectResponse{}
		SiblingDiscriminatorMatches("behavior").ValidateObject(ctx, validator.ObjectRequest{Path: path.Root("value"), Config: config, ConfigValue: value}, resp)
		return resp.Diagnostics.ErrorsCount()
	}
	assert.Zero(t, check("webhook"))
	assert.Equal(t, 1, check("storage_region"))
}
