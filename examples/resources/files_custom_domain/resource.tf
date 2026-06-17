resource "files_custom_domain" "example_custom_domain" {
  destination        = "site_alias"
  folder_behavior_id = 1
  ssl_certificate_id = 1
  domain             = "files.example.com"
}

