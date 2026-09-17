package lib

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type structConversionItem struct {
	Path      string                 `json:"path"`
	Recursive bool                   `json:"recursive"`
	Nested    *structConversionChild `json:"nested,omitempty"`
}

type structConversionChild struct {
	Count int64 `json:"count"`
}

// SDK entity arrays such as []files_sdk.BundlePath arrive as typed structs, not generic JSON values.
func TestToDynamicConvertsStructValues(t *testing.T) {
	ctx := context.Background()
	value, diags := ToDynamic(ctx, path.Root("items"), []structConversionItem{{Path: "reports", Recursive: true, Nested: &structConversionChild{Count: 2}}}, nil)
	require.False(t, diags.HasError(), diags)
	actual, diags := AttributeToInterface(ctx, path.Root("items"), value)
	require.False(t, diags.HasError(), diags)
	assert.Equal(t, []any{map[string]any{"path": "reports", "recursive": true, "nested": map[string]any{"count": float64(2)}}}, actual)

	empty, diags := ToDynamic(ctx, path.Root("items"), []structConversionItem(nil), nil)
	require.False(t, diags.HasError(), diags)
	assert.True(t, empty.IsNull())

	single, diags := ToDynamic(ctx, path.Root("item"), structConversionItem{Path: "reports"}, nil)
	require.False(t, diags.HasError(), diags)
	actual, diags = AttributeToInterface(ctx, path.Root("item"), single)
	require.False(t, diags.HasError(), diags)
	assert.Equal(t, map[string]any{"path": "reports", "recursive": false}, actual)
}

func TestToDynamicRejectsUnsupportedFallbackValues(t *testing.T) {
	for name, source := range map[string]any{
		"integer":   int(5),
		"typed map": map[string]string{"path": "reports"},
		"pointer":   &structConversionItem{Path: "reports"},
		"time":      time.Date(2026, time.September, 16, 0, 0, 0, 0, time.UTC),
		"bytes":     []byte("reports"),
		"channel":   make(chan int),
	} {
		t.Run(name, func(t *testing.T) {
			_, diags := ToDynamic(context.Background(), path.Root("item"), source, nil)
			require.True(t, diags.HasError(), "unsupported value must not silently enter state")
			assert.Contains(t, diags.Errors()[0].Detail(), "item")
		})
	}
}
