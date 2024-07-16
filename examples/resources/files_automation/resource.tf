resource "files_automation" "example_automation" {
  source                               = "source"
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
  schedule                             = "example"
  schedule_days_of_week                = [0, 1, 3]
  schedule_times_of_day                = ["7:30", "11:30"]
  schedule_time_zone                   = "Eastern Time (US & Canada)"
  always_overwrite_size_matching_files = true
  description                          = "example"
  disabled                             = true
  flatten_destination_structure        = true
  ignore_locked_folders                = true
  legacy_folder_matching               = true
  name                                 = "example"
  overwrite_files                      = true
  path_time_zone                       = "Eastern Time (US & Canada)"
  trigger                              = "daily"
  trigger_actions                      = ["create"]
  value                                = {
    limit = 1
  }
  recurring_day                        = 25
  automation                           = "create_folder"
}

