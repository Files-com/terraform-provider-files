resource "files_event_subscription" "example_event_subscription" {
  event_channel_id        = 1
  workspace_id            = 1
  apply_to_all_workspaces = true
  name                    = "example"
  enabled                 = true
  event_types             = ["example"]
  delivery_policy         = "example"
  event_target_ids        = [1]
}

