resource "files_permission" "example_permission" {
  path       = "path"
  group_id   = 1
  permission = "full"
  recursive  = true
  user_id    = 1
  username   = "user"
  group_name = "example"
}

