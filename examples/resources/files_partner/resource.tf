resource "files_partner" "example_partner" {
  allow_bypassing_2fa_policies = false
  allow_credential_changes     = false
  allow_providing_gpg_keys     = false
  allow_user_creation          = false
  notes                        = "This is a note about the partner."
  tags                         = "example"
  name                         = "Acme Corp"
  root_folder                  = "/AcmeCorp"
  workspace_id                 = 1
}

