resource "files_secret" "example_secret" {
  name         = "Production API token"
  description  = "Used by production API integrations."
  secret_type  = "token"
  metadata     = {
    key = "example value"
  }
  workspace_id = 0
}

