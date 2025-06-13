resource "files_remote_mount_backend" "example_remote_mount_backend" {
  canary_file_path       = "backend1.txt"
  remote_server_mount_id = 1
  remote_server_id       = 1
  enabled                = true
  fall                   = 1
  health_check_enabled   = true
  health_check_type      = "active"
  interval               = 60
  min_free_cpu           = 1.0
  min_free_mem           = 1.0
  priority               = 1
  remote_path            = "/path/on/remote"
  rise                   = 1
}

