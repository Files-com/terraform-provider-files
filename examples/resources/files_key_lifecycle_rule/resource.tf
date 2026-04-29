resource "files_key_lifecycle_rule" "example_key_lifecycle_rule" {
  apply_to_all_workspaces = true
  expiration_days         = 365
  key_type                = "gpg"
  inactivity_days         = 12
  name                    = "inactive gpg keys"
  workspace_id            = 12
}

