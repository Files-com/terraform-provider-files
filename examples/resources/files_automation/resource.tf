resource "files_automation" "example_automation" {
  source                               = "example"
  destinations                         = [
    "folder_a/file_a.txt",
    {
      folder_path = "folder_b"
      file_path   = "file_b.txt"
    },
    {
      folder_path = "folder_c"
    }
  ]
  destination_replace_from             = "example"
  destination_replace_to               = "example"
  interval                             = "year"
  path                                 = "example"
  sync_ids                             = [1, 2]
  user_ids                             = [1, 2]
  group_ids                            = [1, 2]
  schedule_days_of_week                = [0, 1, 3]
  schedule_times_of_day                = ["7:30", "11:30"]
  schedule_time_zone                   = "Eastern Time (US & Canada)"
  holiday_region                       = "us_dc"
  always_overwrite_size_matching_files = true
  always_serialize_jobs                = true
  description                          = "example"
  disabled                             = true
  exclude_pattern                      = "path/to/exclude/*"
  import_urls                          = [
    {
      name    = "users.json"
      url     = "http://example.com/users"
      method  = "POST"
      headers = {
        Content-Type = "application/json"
      }
      content = {
        group = "support"
      }
    }
  ]
  flatten_destination_structure        = true
  ignore_locked_folders                = true
  legacy_folder_matching               = false
  name                                 = "example"
  overwrite_files                      = true
  path_time_zone                       = "Eastern Time (US & Canada)"
  retry_on_failure_interval_in_minutes = 60
  retry_on_failure_number_of_attempts  = 10
  trigger                              = "daily"
  trigger_actions                      = ["create"]
  value                                = {
    limit = 1
  }
  recurring_day                        = 25
  automation                           = "create_folder"
}

