resource "files_remote_server_credential" "example_remote_server_credential" {
  name                                          = "My Credential"
  description                                   = "More information or notes about this credential."
  server_type                                   = "s3"
  aws_access_key                                = "example"
  s3_assume_role_arn                            = "example"
  s3_assume_role_duration_seconds               = 1
  cloudflare_access_key                         = "example"
  filebase_access_key                           = "example"
  google_cloud_storage_s3_compatible_access_key = "example"
  linode_access_key                             = "example"
  s3_compatible_access_key                      = "example"
  username                                      = "user"
  wasabi_access_key                             = "example"
  workspace_id                                  = 0
}

