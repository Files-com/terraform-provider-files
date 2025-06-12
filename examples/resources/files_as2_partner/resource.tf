resource "files_as2_partner" "example_as2_partner" {
  enable_dedicated_ips       = true
  http_auth_username         = "username"
  mdn_validation_level       = "none"
  signature_validation_level = "normal"
  server_certificate         = "require_match"
  default_mime_type          = "application/octet-stream"
  additional_http_headers    = {
    key = "example value"
  }
  as2_station_id             = 1
  name                       = "AS2 Partner Name"
  uri                        = "example"
  public_certificate         = "public_certificate"
}

