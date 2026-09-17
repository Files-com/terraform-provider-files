package provider

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Files-com/terraform-provider-files/lib"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAutomationJSONStateUpgrades(t *testing.T) {
	ctx := context.Background()
	r := &automationResource{}
	t.Run("v0/definition", func(t *testing.T) {
		example := "{\"schema_version\":1,\"nodes\":[{\"id\":\"trigger\",\"type\":\"trigger_manual\"},{\"id\":\"create_reports\",\"type\":\"create_folder\",\"config\":{\"destinations\":[\"reports/\"]}}],\"edges\":[{\"from\":\"trigger\",\"to\":\"create_reports\"}]}"
		var apiValue any
		require.NoError(t, json.Unmarshal([]byte(example), &apiValue))
		native, diags := lib.ToDynamic(ctx, path.Root("definition"), apiValue, nil)
		require.False(t, diags.HasError(), diags)
		for _, value := range []types.Dynamic{native, types.DynamicValue(types.StringValue(example)), types.DynamicNull()} {
			upgrader := r.UpgradeState(ctx)[0]
			priorType := upgrader.PriorSchema.Type().TerraformType(ctx).(tftypes.Object)
			attributes := make(map[string]tftypes.Value, len(priorType.AttributeTypes))
			for name, attributeType := range priorType.AttributeTypes {
				attributes[name] = tftypes.NewValue(attributeType, nil)
			}
			prior := tfsdk.State{Schema: *upgrader.PriorSchema, Raw: tftypes.NewValue(priorType, attributes)}
			require.False(t, prior.SetAttribute(ctx, path.Root("definition"), value).HasError())
			response := resource.UpgradeStateResponse{State: tfsdk.State{Schema: r.resourceSchema()}}
			upgrader.StateUpgrader(ctx, resource.UpgradeStateRequest{State: &prior}, &response)
			require.False(t, response.Diagnostics.HasError(), response.Diagnostics)
			var upgraded types.Object
			require.False(t, response.State.GetAttribute(ctx, path.Root("definition"), &upgraded).HasError())
			if value.IsNull() {
				assert.True(t, upgraded.IsNull())
				continue
			}
			actual, diags := lib.SchemaAttributeToInterface(ctx, path.Root("definition"), upgraded)
			require.False(t, diags.HasError(), diags)
			actual, transformDiags0 := lib.UnwrapDiscriminatedUnionAtPath(ctx, path.Root("definition"), actual, []string{"nodes"}, "type", []lib.JSONSchemaVariant{{Name: "trigger_scheduled", Value: "trigger_scheduled"}, {Name: "trigger_manual", Value: "trigger_manual"}, {Name: "trigger_action", Value: "trigger_action"}, {Name: "trigger_webhook", Value: "trigger_webhook"}, {Name: "trigger_email", Value: "trigger_email"}, {Name: "create_folder", Value: "create_folder"}, {Name: "copy_file", Value: "copy_file"}, {Name: "move_file", Value: "move_file"}, {Name: "delete_file", Value: "delete_file"}, {Name: "import_file", Value: "import_file"}, {Name: "run_sync", Value: "run_sync"}, {Name: "as2_send", Value: "as2_send"}, {Name: "send_email", Value: "send_email"}, {Name: "agent_compute", Value: "agent_compute"}, {Name: "set_metadata", Value: "set_metadata"}, {Name: "extract", Value: "extract"}, {Name: "document_convert", Value: "document_convert"}, {Name: "image_convert", Value: "image_convert"}, {Name: "zip", Value: "zip"}, {Name: "unzip", Value: "unzip"}, {Name: "gpg_encrypt", Value: "gpg_encrypt"}, {Name: "gpg_decrypt", Value: "gpg_decrypt"}, {Name: "if", Value: "if"}, {Name: "switch", Value: "switch"}, {Name: "filter", Value: "filter"}, {Name: "join", Value: "join"}, {Name: "aggregate", Value: "aggregate"}, {Name: "wait", Value: "wait"}, {Name: "transform", Value: "transform"}, {Name: "run_automation", Value: "run_automation"}})
			diags.Append(transformDiags0...)
			require.False(t, diags.HasError(), diags)
			encoded, err := json.Marshal(actual)
			require.NoError(t, err)
			assert.JSONEq(t, example, string(encoded))
		}
	})
}
