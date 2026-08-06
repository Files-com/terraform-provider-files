resource "files_schedule" "example_schedule" {
  name                  = "Weekday overnight"
  schedule_days_of_week = [1, 2, 3, 4, 5]
  schedule_times_of_day = ["01:00"]
  schedule_time_zone    = "Eastern Time (US & Canada)"
  holiday_region        = "us"
}

