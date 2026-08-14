package lib

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type objectAttributesValidator struct {
	names   []string
	exactly bool
}

func ExactlyOneOfAttributes(names ...string) validator.Object {
	return objectAttributesValidator{names: names, exactly: true}
}

func AtLeastOneOfAttributes(names ...string) validator.Object {
	return objectAttributesValidator{names: names}
}

func (v objectAttributesValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v objectAttributesValidator) MarkdownDescription(_ context.Context) string {
	requirement := "At least one"
	if v.exactly {
		requirement = "Exactly one"
	}
	return fmt.Sprintf("%s of %s must be specified", requirement, strings.Join(v.names, ", "))
}

func (v objectAttributesValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	count := 0
	attributes := req.ConfigValue.Attributes()
	for _, name := range v.names {
		value := attributes[name]
		if value.IsUnknown() {
			return
		}
		if !value.IsNull() {
			count++
		}
	}

	if count == 0 || (v.exactly && count > 1) {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid Attribute Combination", v.MarkdownDescription(ctx))
	}
}
