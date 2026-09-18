package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func (r *remoteServerCredentialResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: {
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				// Attribute types are unchanged. Returning through a state upgrader lets
				// the framework clear write-only values before sending state to Terraform.
				value, err := req.RawState.UnmarshalWithOpts(resp.State.Schema.Type().TerraformType(ctx), tfprotov6.UnmarshalOpts{
					ValueFromJSONOpts: tftypes.ValueFromJSONOpts{IgnoreUndefinedAttributes: true},
				})
				if err != nil {
					resp.Diagnostics.AddError("Unable to Upgrade Resource State", err.Error())
					return
				}
				resp.State.Raw = value
			},
		},
	}
}
