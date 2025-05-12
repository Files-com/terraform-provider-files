resource "files_bundle_notification" "example_bundle_notification" {
  user_id                = 1
  bundle_id              = 1
  notify_user_id         = 1
  notify_on_registration = true
  notify_on_upload       = true
}

