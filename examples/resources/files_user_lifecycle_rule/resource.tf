resource "files_user_lifecycle_rule" "example_user_lifecycle_rule" {
  authentication_method = "password"
  group_ids             = [1, 2, 3]
  inactivity_days       = 12
  include_site_admins   = true
  include_folder_admins = true
  user_state            = "inactive"
  name                  = "password specific rules"
}

