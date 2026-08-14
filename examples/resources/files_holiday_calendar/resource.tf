resource "files_holiday_calendar" "example_holiday_calendar" {
  definition = {
    months = {
      "0"  = {
        calculated_rules = [
          {
            name              = "Good Friday"
            function          = "easter(year)"
            function_modifier = -2
          }
        ]
      }
      "1"  = {
        fixed_rules   = [
          {
            name     = "New Year's Day"
            mday     = 1
            observed = "to_weekday_if_weekend(date)"
          }
        ]
        weekday_rules = [
          {
            name = "Third Monday"
            week = 3
            wday = 1
          }
        ]
      }
      "11" = {
        weekday_rules = [
          {
            name = "Thanksgiving"
            week = 4
            wday = 4
          }
        ]
      }
      "12" = {
        fixed_rules = [
          {
            name        = "Christmas Eve Early Close"
            mday        = 24
            start_time  = "13:00"
            end_time    = "17:00"
            year_ranges = {
              from = 2026
            }
          }
        ]
      }
    }
  }
  name       = "Company Holidays"
}

