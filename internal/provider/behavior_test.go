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
				Config: providerConfig + `data "files_behavior" "foo_remote_server_sync" { id = 272834 }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.files_behavior.foo_remote_server_sync", "id", "272834"),
					resource.TestCheckResourceAttr("data.files_behavior.foo_remote_server_sync", "behavior", "remote_server_sync"),
					resource.TestCheckResourceAttr("data.files_behavior.foo_remote_server_sync", "path", "Foo"),
					resource.TestCheckResourceAttr("data.files_behavior.foo_remote_server_sync", "value.direction", "pull_from_server"),
					resource.TestCheckResourceAttr("data.files_behavior.foo_remote_server_sync", "value.schedule.days_of_week.0", "0"),
					resource.TestCheckResourceAttr("data.files_behavior.foo_remote_server_sync", "value.schedule.days_of_week.1", "1"),
					resource.TestCheckResourceAttr("data.files_behavior.foo_remote_server_sync", "value.schedule.times_of_day.0", "12:00"),
					resource.TestCheckResourceAttr("data.files_behavior.foo_remote_server_sync", "value.schedule.times_of_day.1", "13:00"),
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
resource "files_behavior" "bar_remote_server_sync" {
  behavior  = "remote_server_sync"
  path      = "Bar"
  value     = {
    direction        = "pull_from_server"
    keep_after_copy  = "keep"
	remote_server_id = 21226
    trigger          = "custom_schedule"
    schedule         = {
      days_of_week = [0, 1]
      times_of_day = ["12:00", "13:00"]
    }
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("files_behavior.bar_remote_server_sync", "id", "272843"),
					resource.TestCheckResourceAttr("files_behavior.bar_remote_server_sync", "behavior", "remote_server_sync"),
					resource.TestCheckResourceAttr("files_behavior.bar_remote_server_sync", "path", "Bar"),
					resource.TestCheckResourceAttr("files_behavior.bar_remote_server_sync", "recursive", "false"),
					resource.TestCheckResourceAttr("files_behavior.bar_remote_server_sync", "value.direction", "pull_from_server"),
					resource.TestCheckResourceAttr("files_behavior.bar_remote_server_sync", "value.schedule.days_of_week.0", "0"),
					resource.TestCheckResourceAttr("files_behavior.bar_remote_server_sync", "value.schedule.days_of_week.1", "1"),
					resource.TestCheckResourceAttr("files_behavior.bar_remote_server_sync", "value.schedule.times_of_day.0", "12:00"),
					resource.TestCheckResourceAttr("files_behavior.bar_remote_server_sync", "value.schedule.times_of_day.1", "13:00"),
				),
			},
			{
				Config: providerConfig + `
resource "files_behavior" "bar_remote_server_sync" {
  behavior  = "remote_server_sync"
  path      = "Bar"
  recursive = true
  value     = {
    direction        = "pull_from_server"
    keep_after_copy  = "keep"
	remote_server_id = 21226
    trigger          = "custom_schedule"
    schedule         = {
      days_of_week = [4, 5]
      times_of_day = ["02:00", "03:00"]
    }
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("files_behavior.bar_remote_server_sync", "id", "272843"),
					resource.TestCheckResourceAttr("files_behavior.bar_remote_server_sync", "behavior", "remote_server_sync"),
					resource.TestCheckResourceAttr("files_behavior.bar_remote_server_sync", "path", "Bar"),
					resource.TestCheckResourceAttr("files_behavior.bar_remote_server_sync", "recursive", "true"),
					resource.TestCheckResourceAttr("files_behavior.bar_remote_server_sync", "value.direction", "pull_from_server"),
					resource.TestCheckResourceAttr("files_behavior.bar_remote_server_sync", "value.schedule.days_of_week.0", "4"),
					resource.TestCheckResourceAttr("files_behavior.bar_remote_server_sync", "value.schedule.days_of_week.1", "5"),
					resource.TestCheckResourceAttr("files_behavior.bar_remote_server_sync", "value.schedule.times_of_day.0", "02:00"),
					resource.TestCheckResourceAttr("files_behavior.bar_remote_server_sync", "value.schedule.times_of_day.1", "03:00"),
				),
			},
		},
	})
}
