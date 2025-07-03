resource "files_user_lifecycle_rule" "example_user_lifecycle_rule" {
  authentication_method = "password"
  inactivity_days       = 12
  include_site_admins   = true
  include_folder_admins = true
  user_state            = "inactive"
  name                  = "password specific rules"
}

