package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestFileDataSource(t *testing.T) {
	VcrTest(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(t.Name()),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `data "files_file" "test_file" { path = "Test Folder/Test File.txt" }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.files_file.test_file", "path", "Test Folder/Test File.txt"),
					resource.TestCheckResourceAttr("data.files_file.test_file", "provided_mtime", "2010-01-01T01:02:03Z"),
					resource.TestCheckResourceAttr("data.files_file.test_file", "priority_color", "blue"),
					resource.TestCheckResourceAttr("data.files_file.test_file", "custom_metadata.foo", "bar"),
					resource.TestCheckResourceAttr("data.files_file.test_file", "custom_metadata.baz", "qux"),
				),
			},
		},
	})
}

func TestFileResource(t *testing.T) {
	VcrTest(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(t.Name()),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "files_file" "bar_file" {
  source          = "fixtures/test_file.txt"
  path            = "Test Folder/Subfolder/Bar.txt"
  provided_mtime  = "2024-06-01T01:02:03Z"
  priority_color  = "orange"
  custom_metadata = {
    key = "value"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("files_file.bar_file", "path", "Test Folder/Subfolder/Bar.txt"),
					resource.TestCheckResourceAttr("files_file.bar_file", "provided_mtime", "2024-06-01T01:02:03Z"),
					resource.TestCheckResourceAttr("files_file.bar_file", "priority_color", "orange"),
					resource.TestCheckResourceAttr("files_file.bar_file", "custom_metadata.key", "value"),
				),
			},
			{
				Config: providerConfig + `
resource "files_file" "bar_file" {
	source          = "fixtures/test_file.txt"
	path            = "Test Folder/Subfolder/Bar.txt"
	provided_mtime  = ""
	priority_color  = "red"
	custom_metadata = {
		custom = "metadata"
	}
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("files_file.bar_file", "path", "Test Folder/Subfolder/Bar.txt"),
					resource.TestCheckResourceAttr("files_file.bar_file", "provided_mtime", ""),
					resource.TestCheckResourceAttr("files_file.bar_file", "priority_color", "red"),
					resource.TestCheckResourceAttr("files_file.bar_file", "custom_metadata.custom", "metadata"),
				),
			},
			{
				Config: providerConfig + `
resource "files_file" "bar_file" {
	source          = "fixtures/test_file.txt"
	path            = "Bar Moved.txt"
	provided_mtime  = "2020-06-01T01:02:03Z"
	priority_color  = "red"
	custom_metadata = {
		custom = "metadata"
	}
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("files_file.bar_file", "path", "Bar Moved.txt"),
					resource.TestCheckResourceAttr("files_file.bar_file", "provided_mtime", "2020-06-01T01:02:03Z"),
					resource.TestCheckResourceAttr("files_file.bar_file", "priority_color", "red"),
					resource.TestCheckResourceAttr("files_file.bar_file", "custom_metadata.custom", "metadata"),
				),
			},
		},
	})
}
