package lib

import (
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// SettleUnchangedPlan keeps prior state when the plan differs from it only in computed attributes that have no
// configuration and are unknown. Terraform proposes null for an omitted Optional and Computed typed object whose prior
// value has fields that are not Computed. The framework then marks every computed attribute without configuration
// unknown, although the object's UseStateForUnknown plan modifier keeps its prior value. Without this, every plan would
// show an update that changes nothing.
func SettleUnchangedPlan(config tfsdk.Config, state tfsdk.State, plan *tfsdk.Plan) error {
	if state.Raw.IsNull() || plan.Raw.IsNull() {
		return nil
	}
	settled, err := tftypes.Transform(plan.Raw, func(attributePath *tftypes.AttributePath, value tftypes.Value) (tftypes.Value, error) {
		if value.IsKnown() {
			return value, nil
		}
		// An unknown value with configuration comes from the configuration, such as another resource's pending output.
		if configured, _, err := tftypes.WalkAttributePath(config.Raw, attributePath); err == nil {
			if configuredValue, ok := configured.(tftypes.Value); ok && !configuredValue.IsNull() {
				return value, nil
			}
		}
		prior, _, err := tftypes.WalkAttributePath(state.Raw, attributePath)
		if err != nil {
			return value, nil
		}
		if priorValue, ok := prior.(tftypes.Value); ok {
			return priorValue, nil
		}
		return value, nil
	})
	if err != nil {
		return err
	}
	if settled.Equal(state.Raw) {
		plan.Raw = state.Raw
	}
	return nil
}
