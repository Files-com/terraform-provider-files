resource "files_behavior" "example_behavior" {
  value                          = {
    method = "GET"
  }
  disable_parent_folder_behavior = true
  recursive                      = true
  name                           = "example"
  description                    = "example"
  path                           = "path"
  behavior                       = "webhook"
}
