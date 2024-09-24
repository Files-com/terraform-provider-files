package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestSiteResource(t *testing.T) {
	VcrTest(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(t.Name()),
		Steps: []resource.TestStep{
			{
				Config:        providerConfig + `resource "files_site" "my_site" {}`,
				ResourceName:  "files_site.my_site",
				ImportState:   true,
				ImportStateId: "12345",
				ImportStateCheck: func(s []*terraform.InstanceState) error {
					if len(s) != 1 {
						return fmt.Errorf("expected 1 state: %#v", s)
					}

					site := s[0]
					expected := map[string]string{
						"always_mkdir_parents": "false",
						"ftp_enabled":          "true",
					}

					for key, expectedValue := range expected {
						if site.Attributes[key] != expectedValue {
							return fmt.Errorf("expected site `%s` attribute to be %s, got: %s", key, expectedValue, site.Attributes[key])
						}
					}

					return nil
				},
			},
		},
	})
}
