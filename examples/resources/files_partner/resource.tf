resource "files_partner" "example_partner" {
  allow_bypassing_2fa_policies = false
  allow_credential_changes     = false
  allow_providing_gpg_keys     = false
  allow_user_creation          = false
  notes                        = "This is a note about the partner."
  root_folder                  = "/AcmeCorp"
  tags                         = "example"
  name                         = "Acme Corp"
}

