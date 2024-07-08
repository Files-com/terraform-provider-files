package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestFolderDataSource(t *testing.T) {
	VcrTest(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(t.Name()),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `data "files_folder" "test_folder" { path = "Test Folder" }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.files_folder.test_folder", "path", "Test Folder"),
					resource.TestCheckResourceAttr("data.files_folder.test_folder", "provided_mtime", "2010-01-01T01:02:03Z"),
					resource.TestCheckResourceAttr("data.files_folder.test_folder", "priority_color", "blue"),
					resource.TestCheckResourceAttr("data.files_folder.test_folder", "custom_metadata.foo", "bar"),
					resource.TestCheckResourceAttr("data.files_folder.test_folder", "custom_metadata.baz", "qux"),
				),
			},
		},
	})
}

func TestFolderResource(t *testing.T) {
	VcrTest(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(t.Name()),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "files_folder" "bar_folder" {
  path            = "Test Folder/Subfolder/Bar"
  mkdir_parents   = true
  provided_mtime  = "2024-06-01T01:02:03Z"
  priority_color  = "orange"
  custom_metadata = {
    key = "value"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("files_folder.bar_folder", "path", "Test Folder/Subfolder/Bar"),
					resource.TestCheckResourceAttr("files_folder.bar_folder", "provided_mtime", "2024-06-01T01:02:03Z"),
					resource.TestCheckResourceAttr("files_folder.bar_folder", "priority_color", "orange"),
					resource.TestCheckResourceAttr("files_folder.bar_folder", "custom_metadata.key", "value"),
				),
			},
			{
				Config: providerConfig + `
resource "files_folder" "bar_folder" {
	path            = "Test Folder/Subfolder/Bar"
	mkdir_parents   = true
	provided_mtime  = "2020-06-01T01:02:03Z"
	priority_color  = "red"
	custom_metadata = {
		custom = "metadata"
	}
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("files_folder.bar_folder", "path", "Test Folder/Subfolder/Bar"),
					resource.TestCheckResourceAttr("files_folder.bar_folder", "provided_mtime", "2020-06-01T01:02:03Z"),
					resource.TestCheckResourceAttr("files_folder.bar_folder", "priority_color", "red"),
					resource.TestCheckResourceAttr("files_folder.bar_folder", "custom_metadata.custom", "metadata"),
				),
			},
			{
				Config: providerConfig + `
resource "files_folder" "bar_folder" {
	path            = "Test Folder/Bar Moved"
	mkdir_parents   = true
	provided_mtime  = "2020-06-01T01:02:03Z"
	priority_color  = "red"
	custom_metadata = {
		custom = "metadata"
	}
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("files_folder.bar_folder", "path", "Test Folder/Bar Moved"),
					resource.TestCheckResourceAttr("files_folder.bar_folder", "provided_mtime", "2020-06-01T01:02:03Z"),
					resource.TestCheckResourceAttr("files_folder.bar_folder", "priority_color", "red"),
					resource.TestCheckResourceAttr("files_folder.bar_folder", "custom_metadata.custom", "metadata"),
				),
			},
		},
	})
}
