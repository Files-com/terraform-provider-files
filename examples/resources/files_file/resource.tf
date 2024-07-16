resource "files_file" "example_file" {
  source          = "path"
  md5             = "17c54824e9931a4688ca032d03f6663c"
  path            = "path"
  custom_metadata = {
    key = "value"
  }
  provided_mtime  = "2000-01-01T01:00:00Z"
  priority_color  = "red"
}

