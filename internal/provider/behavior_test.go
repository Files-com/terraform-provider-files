package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestBehaviorDataSource(t *testing.T) {
	VcrTest(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(t.Name()),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `data "files_behavior" "foo_file_expiration" { id = 272834 }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.files_behavior.foo_file_expiration", "id", "272834"),
					resource.TestCheckResourceAttr("data.files_behavior.foo_file_expiration", "behavior", "file_expiration"),
					resource.TestCheckResourceAttr("data.files_behavior.foo_file_expiration", "path", "Foo"),
					resource.TestCheckResourceAttr("data.files_behavior.foo_file_expiration", "value", "14"),
				),
			},
		},
	})
}

func TestBehaviorResource(t *testing.T) {
	VcrTest(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(t.Name()),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "files_behavior" "bar_serve_publicly" {
  path     = "Bar"
  behavior = "serve_publicly"
  value    = {
    key = "Bar",
    show_index = true
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("files_behavior.bar_serve_publicly", "id", "272843"),
					resource.TestCheckResourceAttr("files_behavior.bar_serve_publicly", "behavior", "serve_publicly"),
					resource.TestCheckResourceAttr("files_behavior.bar_serve_publicly", "path", "Bar"),
					resource.TestCheckResourceAttr("files_behavior.bar_serve_publicly", "value.key", "Bar"),
					resource.TestCheckResourceAttr("files_behavior.bar_serve_publicly", "value.show_index", "true"),
				),
			},
			{
				Config: providerConfig + `
resource "files_behavior" "bar_serve_publicly" {
  path     = "Bar"
  behavior = "serve_publicly"
  value    = {
    key = "Bar",
    show_index = false
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("files_behavior.bar_serve_publicly", "id", "272843"),
					resource.TestCheckResourceAttr("files_behavior.bar_serve_publicly", "behavior", "serve_publicly"),
					resource.TestCheckResourceAttr("files_behavior.bar_serve_publicly", "path", "Bar"),
					resource.TestCheckResourceAttr("files_behavior.bar_serve_publicly", "value.key", "Bar"),
					resource.TestCheckResourceAttr("files_behavior.bar_serve_publicly", "value.show_index", "false"),
				),
			},
			{
				Config: providerConfig + `
resource "files_behavior" "primitive_file_expiration" {
  behavior  = "file_expiration"
  path      = "Bar"
  value     = 14
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("files_behavior.primitive_file_expiration", "behavior", "file_expiration"),
					resource.TestCheckResourceAttr("files_behavior.primitive_file_expiration", "path", "Bar"),
					resource.TestCheckResourceAttr("files_behavior.primitive_file_expiration", "value", "14"),
				),
			},
		},
	})
}
