package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestUserDataSourceByUsername(t *testing.T) {
	VcrTest(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(t.Name()),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `data "files_user" "test_user" { username = "test_user" }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.files_user.test_user", "id", "1105124"),
					resource.TestCheckResourceAttr("data.files_user.test_user", "username", "test_user"),
					resource.TestCheckResourceAttr("data.files_user.test_user", "name", "Test User"),
				),
			},
		},
	})
}
