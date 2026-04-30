resource "files_desktop_configuration_profile" "example_desktop_configuration_profile" {
  name                   = "North America Desktop Profile"
  mount_mappings         = {
    key = "example value"
  }
  workspace_id           = 1
  use_for_all_users      = false
  disable_drive_mounting = false
}

