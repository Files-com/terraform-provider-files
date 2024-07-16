resource "files_notification" "example_notification" {
  user_id                     = 1
  notify_on_copy              = true
  notify_on_delete            = true
  notify_on_download          = true
  notify_on_move              = true
  notify_on_upload            = true
  notify_user_actions         = true
  recursive                   = true
  send_interval               = "daily"
  message                     = "custom notification email message"
  triggering_filenames        = ["*.jpg", "notify_file.txt"]
  triggering_group_ids        = [1]
  triggering_user_ids         = [1]
  trigger_by_share_recipients = true
  group_id                    = 1
  username                    = "User"
}

