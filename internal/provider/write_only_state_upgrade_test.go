package provider

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/require"
)

func TestStateUpgradeClearsStoredWriteOnlyValues(t *testing.T) {
	ctx := context.Background()
	server := providerserver.NewProtocol6(New("test")())()
	schemas, err := server.GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})
	require.NoError(t, err)
	require.Empty(t, schemas.Diagnostics)

	for _, test := range []struct {
		resource  string
		preserved map[string]any
		writeOnly map[string]any
	}{
		{
			resource:  "files_remote_server",
			preserved: map[string]any{"id": 123, "name": "test-server", "hostname": "example.invalid"},
			writeOnly: map[string]any{"private_key": "old-private-key", "reset_authentication": true},
		},
		{
			resource:  "files_gpg_key",
			preserved: map[string]any{"id": 456, "name": "test-key", "generated_private_key": "retained-generated-key"},
			writeOnly: map[string]any{"generate_email": "test@example.invalid", "generate_keypair": true},
		},
		{
			resource:  "files_folder",
			preserved: map[string]any{"path": "test-folder"},
			writeOnly: map[string]any{"mkdir_parents": true},
		},
	} {
		t.Run(test.resource, func(t *testing.T) {
			attributes := make(map[string]any)
			for name, value := range test.preserved {
				attributes[name] = value
			}
			for name, value := range test.writeOnly {
				attributes[name] = value
			}
			// Older providers can also have attributes that have since been removed.
			attributes["removed_attribute"] = "obsolete"
			raw, err := json.Marshal(attributes)
			require.NoError(t, err)
			resourceSchema := schemas.ResourceSchemas[test.resource]
			resourceType := resourceSchema.ValueType()

			// The second pass checks that state already written by this provider stays unchanged.
			for _, version := range []int64{0, resourceSchema.Version} {
				response, err := server.UpgradeResourceState(ctx, &tfprotov6.UpgradeResourceStateRequest{
					TypeName: test.resource,
					Version:  version,
					RawState: &tfprotov6.RawState{JSON: raw},
				})
				require.NoError(t, err)
				require.Empty(t, response.Diagnostics)
				require.NotNil(t, response.UpgradedState)
				upgraded, err := response.UpgradedState.Unmarshal(resourceType)
				require.NoError(t, err)
				var actual map[string]tftypes.Value
				require.NoError(t, upgraded.As(&actual))
				for name := range test.writeOnly {
					require.True(t, actual[name].IsNull(), "%s must not be returned in state", name)
					attributes[name] = nil
				}
				delete(attributes, "removed_attribute")
				// Compare all state, so dropping unrelated values also fails the regression.
				raw, err = json.Marshal(attributes)
				require.NoError(t, err)
				expected, err := (&tfprotov6.RawState{JSON: raw}).Unmarshal(resourceType)
				require.NoError(t, err)
				require.True(t, expected.Equal(upgraded), "unrelated state must be preserved")
			}
		})
	}
}
