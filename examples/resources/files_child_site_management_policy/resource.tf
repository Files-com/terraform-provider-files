resource "files_child_site_management_policy" "example_child_site_management_policy" {
  value               = {
    color2_left = "#000000"
  }
  skip_child_site_ids = [1, 2]
  child_site_ids      = [1, 2]
  default_policy      = true
  policy_type         = "settings"
  name                = "example"
  description         = "example"
}

