package lib

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestObjectAttributesValidators(t *testing.T) {
	attributeTypes := map[string]attr.Type{"first": types.StringType, "second": types.StringType}
	tests := []struct {
		message   string
		validator validator.Object
		value     types.Object
		expectErr bool
	}{
		{message: "optional null object", validator: ExactlyOneOfAttributes("first", "second"), value: types.ObjectNull(attributeTypes)},
		{message: "exactly one", validator: ExactlyOneOfAttributes("first", "second"), value: types.ObjectValueMust(attributeTypes, map[string]attr.Value{"first": types.StringValue("value"), "second": types.StringNull()})},
		{message: "exactly one missing", validator: ExactlyOneOfAttributes("first", "second"), value: types.ObjectValueMust(attributeTypes, map[string]attr.Value{"first": types.StringNull(), "second": types.StringNull()}), expectErr: true},
		{message: "exactly one exceeded", validator: ExactlyOneOfAttributes("first", "second"), value: types.ObjectValueMust(attributeTypes, map[string]attr.Value{"first": types.StringValue("value"), "second": types.StringValue("value")}), expectErr: true},
		{message: "at least one allows both", validator: AtLeastOneOfAttributes("first", "second"), value: types.ObjectValueMust(attributeTypes, map[string]attr.Value{"first": types.StringValue("value"), "second": types.StringValue("value")})},
	}

	for _, test := range tests {
		t.Run(test.message, func(t *testing.T) {
			request := validator.ObjectRequest{ConfigValue: test.value, Path: path.Root("test")}
			response := &validator.ObjectResponse{}
			test.validator.ValidateObject(context.Background(), request, response)
			assert.Equal(t, test.expectErr, response.Diagnostics.HasError())
		})
	}
}
