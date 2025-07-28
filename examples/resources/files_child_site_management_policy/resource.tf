resource "files_child_site_management_policy" "example_child_site_management_policy" {
  site_setting_name   = "color2_left"
  managed_value       = "#FF0000"
  skip_child_site_ids = [1, 5]
}

