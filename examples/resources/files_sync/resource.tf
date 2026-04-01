resource "files_sync" "example_sync" {
  delete_empty_folders  = true
  description           = "example"
  dest_path             = "example"
  dest_remote_server_id = 1
  dest_site_id          = 1
  disabled              = true
  exclude_patterns      = ["example"]
  holiday_region        = "us_dc"
  include_patterns      = ["example"]
  interval              = "week"
  keep_after_copy       = true
  name                  = "example"
  recurring_day         = 25
  schedule_days_of_week = [0, 2, 4]
  schedule_time_zone    = "Eastern Time (US & Canada)"
  schedule_times_of_day = ["06:30", "14:30"]
  src_path              = "example"
  src_remote_server_id  = 1
  src_site_id           = 1
  sync_interval_minutes = 1
  trigger               = "example"
  trigger_file          = "example"
  workspace_id          = 1
}

