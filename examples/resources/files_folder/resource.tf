resource "files_folder" "example_folder" {
  path            = "path"
  mkdir_parents   = true
  provided_mtime  = "2000-01-01T01:00:00Z"
  custom_metadata = {
    key = "value"
  }
  priority_color  = "red"
}
