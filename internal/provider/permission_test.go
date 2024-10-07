package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestPermissionResource(t *testing.T) {
	VcrTest(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(t.Name()),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "files_permission" "root_permission" {
  path       = ""
  permission = "full"
  recursive  = true
  username   = "example_user"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("files_permission.root_permission", "path", ""),
					resource.TestCheckResourceAttr("files_permission.root_permission", "permission", "full"),
				),
			},
			{
				Config: providerConfig + `
resource "files_permission" "root_permission" {
  path       = ""
  permission = "readonly"
  recursive  = true
  username   = "example_user"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("files_permission.root_permission", "path", ""),
					resource.TestCheckResourceAttr("files_permission.root_permission", "permission", "readonly"),
				),
			},
		},
	})
}
