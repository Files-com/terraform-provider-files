resource "files_partner" "example_partner" {
  allow_bypassing_2fa_policies = false
  allow_credential_changes     = false
  allow_user_creation          = false
  name                         = "Acme Corp"
  notes                        = "This is a note about the partner."
  root_folder                  = "/AcmeCorp"
}

