resource "files_gpg_key" "example_gpg_key" {
  user_id             = 1
  partner_id          = 1
  name                = "key name"
  workspace_id        = 0
  generate_expires_at = "2025-06-19 12:00:00"
  generate_keypair    = false
  generate_full_name  = "John Doe"
  generate_email      = "jdoe@example.com"
}

