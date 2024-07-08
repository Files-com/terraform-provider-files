resource "files_bundle" "example_bundle" {
  user_id                             = 1
  paths                               = ["file.txt"]
  password                            = "Password"
  form_field_set_id                   = 1
  create_snapshot                     = true
  dont_separate_submissions_by_folder = true
  expires_at                          = "2000-01-01T01:00:00Z"
  finalize_snapshot                   = true
  max_uses                            = 1
  description                         = "The public description of the bundle."
  note                                = "The internal note on the bundle."
  code                                = "abc123"
  path_template                       = "{{name}}_{{ip}}"
  path_template_time_zone             = "Eastern Time (US & Canada)"
  permissions                         = "read"
  require_registration                = true
  clickwrap_id                        = 1
  inbox_id                            = 1
  require_share_recipient             = true
  send_email_receipt_to_uploader      = true
  skip_email                          = true
  skip_name                           = true
  skip_company                        = true
  start_access_on_date                = "2000-01-01T01:00:00Z"
  snapshot_id                         = 1
}
