resource "files_group" "example_group" {
  notes              = "example"
  user_ids           = 1
  admin_ids          = 1
  ftp_permission     = true
  sftp_permission    = true
  dav_permission     = true
  restapi_permission = true
  allowed_ips        = "10.0.0.0/8\n127.0.0.1"
  name               = "name"
}
