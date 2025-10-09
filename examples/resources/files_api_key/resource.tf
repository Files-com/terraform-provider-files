resource "files_api_key" "example_api_key" {
  user_id               = 1
  description           = "example"
  expires_at            = "2000-01-01T01:00:00Z"
  permission_set        = "full"
  name                  = "My Main API Key"
  aws_style_credentials = true
  path                  = "shared/docs"
}

