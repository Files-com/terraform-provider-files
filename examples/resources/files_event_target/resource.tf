resource "files_event_target" "example_event_target" {
  name                    = "example"
  workspace_id            = 1
  apply_to_all_workspaces = true
  enabled                 = true
  config                  = {
    url               = "https://example.com/files-events"
    headers           = {
      X-Example = "value"
    }
    use_dedicated_ips = false
  }
  delivery_policy         = {
    max_attempts        = 3
    retry_delay_minutes = [1, 10]
  }
  target_type             = "webhook"
}

