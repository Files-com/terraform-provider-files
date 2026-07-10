resource "files_integration_centric_profile" "example_integration_centric_profile" {
  name                    = "Business Systems Onboarding"
  expected_remote_servers = ["example"]
  workspace_id            = 1
  use_for_all_users       = false
}

