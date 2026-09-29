resource "files_custom_domain" "example_custom_domain" {
  available_to_all_workspaces = false
  workspace_id                = 0
  destination                 = "site_alias"
  folder_behavior_id          = 1
  ssl_certificate_id          = 1
  domain                      = "files.example.com"
}

