resource "files_key_lifecycle_rule" "example_key_lifecycle_rule" {
  key_type        = "gpg"
  inactivity_days = 12
  name            = "inactive gpg keys"
}

