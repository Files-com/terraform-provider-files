resource "files_sync" "example_sync" {
  name                  = "example"
  description           = "example"
  src_path              = "example"
  dest_path             = "example"
  src_remote_server_id  = 1
  dest_remote_server_id = 1
  keep_after_copy       = false
  delete_empty_folders  = false
  disabled              = false
  interval              = "week"
  trigger               = "example"
  trigger_file          = "example"
  holiday_region        = "us_dc"
  sync_interval_minutes = 1
  recurring_day         = 25
  schedule_time_zone    = "Eastern Time (US & Canada)"
  schedule_days_of_week = [0, 2, 4]
  schedule_times_of_day = ["06:30", "14:30"]
}

