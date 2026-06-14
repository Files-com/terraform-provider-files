resource "files_event_channel" "example_event_channel" {
  name            = "example"
  workspace_id    = 1
  description     = "example"
  enabled         = true
  default_channel = true
}

