resource "files_user_lifecycle_rule" "example_user_lifecycle_rule" {
  apply_to_all_workspaces = true
  authentication_method   = "password"
  group_ids               = [1, 2, 3]
  inactivity_days         = 12
  include_site_admins     = true
  include_folder_admins   = true
  name                    = "password specific rules"
  partner_tag             = "guest"
  user_state              = "inactive"
  user_tag                = "guest"
  workspace_id            = 12
}

