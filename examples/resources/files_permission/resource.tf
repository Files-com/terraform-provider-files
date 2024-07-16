resource "files_permission" "example_permission" {
  group_id   = 1
  path       = "example"
  permission = "full"
  recursive  = true
  user_id    = 1
  username   = "Sser"
}

