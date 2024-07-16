resource "files_as2_partner" "example_as2_partner" {
  enable_dedicated_ips = true
  http_auth_username   = "username"
  mdn_validation_level = "none"
  server_certificate   = "require_match"
  as2_station_id       = 1
  name                 = "AS2 Partner Name"
  uri                  = "example"
  public_certificate   = "public_certificate"
}

