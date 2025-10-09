resource "files_permission" "example_permission" {
  path       = "path"
  group_id   = 1
  permission = "full"
  recursive  = false
  partner_id = 1
  user_id    = 1
  username   = "user"
  group_name = "example"
  site_id    = 1
}

