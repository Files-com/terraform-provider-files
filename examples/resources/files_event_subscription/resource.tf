resource "files_event_subscription" "example_event_subscription" {
  event_channel_id        = 1
  workspace_id            = 1
  apply_to_all_workspaces = true
  name                    = "example"
  subject                 = "example"
  message                 = "example"
  message_only            = true
  enabled                 = true
  event_types             = ["sync_run.failure.v1"]
  filter                  = {
    "payload.sync_name"    = {
      prefix = "Nightly"
    }
    "payload.event_errors" = {
      exists = true
    }
  }
  delivery_policy         = {
    max_attempts        = 3
    retry_delay_minutes = [5, 30]
  }
  event_target_ids        = [1]
}

