resource "files_event_channel" "example_event_channel" {
  name            = "example"
  description     = "example"
  enabled         = true
  default_channel = true
}

