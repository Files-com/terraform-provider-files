package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestGroupResource(t *testing.T) {
	VcrTest(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(t.Name()),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "files_group" "test_group" {
  name     = "example_group"
  user_ids = "1105124   , 942660   , 1173788"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("files_group.test_group", "name", "example_group"),
					resource.TestCheckResourceAttr("files_group.test_group", "user_ids", "1105124   , 942660   , 1173788"),
				),
			},
		},
	})
}
