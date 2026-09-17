resource "files_integration_centric_profile" "example_integration_centric_profile" {
  name                    = "Business Systems Onboarding"
  expected_remote_servers = [
    {
      server_type = "dropbox"
      name        = "Dropbox"
    }
  ]
  workspace_id            = 1
  use_for_all_users       = false
}

