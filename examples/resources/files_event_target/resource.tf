resource "files_event_target" "example_event_target" {
  name                    = "example"
  workspace_id            = 1
  apply_to_all_workspaces = true
  target_type             = "example"
  enabled                 = true
  config                  = "example"
  delivery_policy         = "example"
}

