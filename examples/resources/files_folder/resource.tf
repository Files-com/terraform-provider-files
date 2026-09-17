resource "files_folder" "example_folder" {
  path            = "path"
  mkdir_parents   = false
  provided_mtime  = "2000-01-01T01:00:00Z"
  custom_metadata = {
    department = "finance"
  }
  priority_color  = "red"
}

