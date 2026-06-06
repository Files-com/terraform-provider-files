resource "files_scheduled_export" "example_scheduled_export" {
  name                  = "Monthly access review"
  export_type           = "permission_audit"
  export_options        = {
    group_by = "user"
  }
  user_id               = 1
  disabled              = true
  trigger               = "daily"
  interval              = "month"
  recurring_day         = 1
  schedule_days_of_week = [1, 3, 5]
  schedule_times_of_day = ["06:30"]
  schedule_time_zone    = "Eastern Time (US & Canada)"
  holiday_region        = "us"
}

