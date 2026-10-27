package lib

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

type siblingDiscriminatorMatchesValidator struct{ discriminator string }

// SiblingDiscriminatorMatches requires a typed value to set the attribute named by its sibling discriminator, such as
// files_behavior.value under the name in behavior. Without it, a mismatch fails only at apply, when the request is built.
func SiblingDiscriminatorMatches(discriminator string) validator.Object {
	return siblingDiscriminatorMatchesValidator{discriminator: discriminator}
}

func (v siblingDiscriminatorMatchesValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v siblingDiscriminatorMatchesValidator) MarkdownDescription(_ context.Context) string {
	return "The value must be set under the name in " + v.discriminator + "."
}

func (v siblingDiscriminatorMatchesValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	var selector types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root(v.discriminator), &selector)...)
	if resp.Diagnostics.HasError() || selector.IsNull() || selector.IsUnknown() {
		return
	}
	for name, value := range req.ConfigValue.Attributes() {
		if !value.IsNull() && !value.IsUnknown() && name != selector.ValueString() {
			resp.Diagnostics.AddAttributeError(req.Path.AtName(name), "Value does not match "+v.discriminator,
				fmt.Sprintf("The value is set under %s, but %s is %q. Move the settings under %s, or set %s = %q.", name, v.discriminator, selector.ValueString(), selector.ValueString(), v.discriminator, name))
		}
	}
}
