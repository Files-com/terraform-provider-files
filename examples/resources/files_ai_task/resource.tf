resource "files_ai_task" "example_ai_task" {
  description           = "Summarizes files uploaded by the accounting team."
  disabled              = true
  holiday_region        = "us"
  interval              = "day"
  name                  = "Summarize daily reports"
  path                  = "incoming/reports"
  permission_set        = "files_only"
  prompt                = "Summarize the uploaded file and identify follow-up actions."
  recurring_day         = 1
  schedule_days_of_week = [1, 3, 5]
  schedule_time_zone    = "Eastern Time (US & Canada)"
  schedule_times_of_day = ["06:30"]
  source                = "*.pdf"
  trigger               = "daily"
  trigger_actions       = ["create"]
  workspace_id          = 0
}

