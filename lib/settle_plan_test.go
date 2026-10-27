package lib

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A resource with a typed Optional and Computed delivery_policy, as the typed Event Target has.
func TestSettleUnchangedPlan(t *testing.T) {
	policyType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"max_attempts": tftypes.Number}}
	resourceType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"name": tftypes.String, "updated_at": tftypes.String, "delivery_policy": policyType}}
	policy := tftypes.NewValue(policyType, map[string]tftypes.Value{"max_attempts": tftypes.NewValue(tftypes.Number, 3)})
	resourceValue := func(name, updatedAt, deliveryPolicy tftypes.Value) tftypes.Value {
		return tftypes.NewValue(resourceType, map[string]tftypes.Value{"name": name, "updated_at": updatedAt, "delivery_policy": deliveryPolicy})
	}
	name := tftypes.NewValue(tftypes.String, "example")
	unknownTime := tftypes.NewValue(tftypes.String, tftypes.UnknownValue)
	nullString := tftypes.NewValue(tftypes.String, nil)
	state := tfsdk.State{Raw: resourceValue(name, tftypes.NewValue(tftypes.String, "2026-10-15T00:00:00Z"), policy)}

	// delivery_policy omitted: its plan modifier kept the prior value, and the framework marked updated_at unknown.
	config := tfsdk.Config{Raw: resourceValue(name, nullString, tftypes.NewValue(policyType, nil))}
	plan := tfsdk.Plan{Raw: resourceValue(name, unknownTime, policy)}
	require.NoError(t, SettleUnchangedPlan(config, state, &plan))
	assert.True(t, plan.Raw.Equal(state.Raw), "a plan that changes nothing keeps prior state")

	// A changed name is a real update, so updated_at stays unknown.
	renamed := resourceValue(tftypes.NewValue(tftypes.String, "renamed"), unknownTime, policy)
	plan = tfsdk.Plan{Raw: renamed}
	require.NoError(t, SettleUnchangedPlan(tfsdk.Config{Raw: resourceValue(tftypes.NewValue(tftypes.String, "renamed"), nullString, tftypes.NewValue(policyType, nil))}, state, &plan))
	assert.True(t, plan.Raw.Equal(renamed))

	// An unknown that comes from configuration is never replaced with prior state.
	pending := resourceValue(tftypes.NewValue(tftypes.String, tftypes.UnknownValue), tftypes.NewValue(tftypes.String, "2026-10-15T00:00:00Z"), policy)
	plan = tfsdk.Plan{Raw: pending}
	require.NoError(t, SettleUnchangedPlan(tfsdk.Config{Raw: resourceValue(tftypes.NewValue(tftypes.String, tftypes.UnknownValue), nullString, tftypes.NewValue(policyType, nil))}, state, &plan))
	assert.True(t, plan.Raw.Equal(pending))

	// Create and destroy plans are left alone.
	created := resourceValue(name, unknownTime, policy)
	plan = tfsdk.Plan{Raw: created}
	require.NoError(t, SettleUnchangedPlan(config, tfsdk.State{Raw: tftypes.NewValue(resourceType, nil)}, &plan))
	assert.True(t, plan.Raw.Equal(created))
	plan = tfsdk.Plan{Raw: tftypes.NewValue(resourceType, nil)}
	require.NoError(t, SettleUnchangedPlan(config, state, &plan))
	assert.True(t, plan.Raw.IsNull())
}
