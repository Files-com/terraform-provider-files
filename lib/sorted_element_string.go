package lib

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ basetypes.StringTypable                    = SortedElementStringType{}
	_ basetypes.StringValuableWithSemanticEquals = SortedElementString{}
)

type SortedElementStringType struct {
	basetypes.StringType
}

func (t SortedElementStringType) ValueType(ctx context.Context) attr.Value {
	return SortedElementString{}
}

func (t SortedElementStringType) String() string {
	return "SortedElementStringType"
}

func (t SortedElementStringType) Equal(o attr.Type) bool {
	if other, ok := o.(SortedElementStringType); ok {
		return t.StringType.Equal(other.StringType)
	}

	return false
}

func (t SortedElementStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	return SortedElementString{StringValue: in}, nil
}

func (t SortedElementStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.StringType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}

	stringValue, ok := attrValue.(basetypes.StringValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	stringValuable, diags := t.ValueFromString(ctx, stringValue)
	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting StringValue to StringValuable: %v", diags)
	}

	return stringValuable, nil
}

func SortedElementStringValue(value string) SortedElementString {
	return SortedElementString{basetypes.NewStringValue(value)}
}

type SortedElementString struct {
	basetypes.StringValue
}

func (v SortedElementString) Type(ctx context.Context) attr.Type {
	return SortedElementStringType{}
}

func (v SortedElementString) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	// The framework should always pass the correct value type, but always check
	newValue, ok := newValuable.(SortedElementString)

	if !ok {
		diags.AddError(
			"Semantic Equality Check Error",
			"An unexpected value type was received while performing semantic equality checks. "+
				"Please report this to the provider developers.\n\n"+
				"Expected Value Type: "+fmt.Sprintf("%T", v)+"\n"+
				"Got Value Type: "+fmt.Sprintf("%T", newValuable),
		)

		return false, diags
	}

	return normalizeString(v.StringValue.ValueString()) == normalizeString(newValue.StringValue.ValueString()), diags
}

func normalizeString(s string) string {
	items := strings.Split(s, ",")
	for i, item := range items {
		items[i] = strings.TrimSpace(item)
	}
	slices.Sort(items)
	return strings.Join(items, ",")
}
