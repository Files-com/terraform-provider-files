resource "files_expectation" "example_expectation" {
  name                     = "Daily Vendor Feed"
  description              = "Wait for the vendor CSV every morning."
  path                     = "incoming/vendor_a"
  source                   = "*.csv"
  exclude_pattern          = "*.tmp"
  disabled                 = true
  trigger                  = "manual"
  interval                 = "day"
  recurring_day            = 3
  schedule_days_of_week    = [1, 3, 5]
  schedule_times_of_day    = ["06:00"]
  schedule_time_zone       = "UTC"
  holiday_region           = "us"
  lookback_interval        = 3600
  late_acceptance_interval = 900
  inactivity_interval      = 300
  max_open_interval        = 43200
  criteria                 = {
    count      = {
      exact = 1
    }
    extensions = ["csv"]
  }
  workspace_id             = 0
}

