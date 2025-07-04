---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "files_remote_mount_backend Resource - files"
subcategory: ""
description: |-
  A Remote Mount Backend is used to provide high availability for a Remote Server Mount Folder Behavior.
---

# files_remote_mount_backend (Resource)

A Remote Mount Backend is used to provide high availability for a Remote Server Mount Folder Behavior.

## Example Usage

```terraform
resource "files_remote_mount_backend" "example_remote_mount_backend" {
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
  canary_file_path       = "backend1.txt"
  remote_server_mount_id = 1
  remote_server_id       = 1
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `canary_file_path` (String) Path to the canary file used for health checks.
- `remote_server_id` (Number) The remote server that this backend is associated with.
- `remote_server_mount_id` (Number) The mount ID of the Remote Server Mount that this backend is associated with.

### Optional

- `enabled` (Boolean) True if this backend is enabled.
- `fall` (Number) Number of consecutive failures before considering the backend unhealthy.
- `health_check_enabled` (Boolean) True if health checks are enabled for this backend.
- `health_check_type` (String) Type of health check to perform.
- `interval` (Number) Interval in seconds between health checks.
- `min_free_cpu` (String) Minimum free CPU percentage required for this backend to be considered healthy.
- `min_free_mem` (String) Minimum free memory percentage required for this backend to be considered healthy.
- `priority` (Number) Priority of this backend.
- `remote_path` (String) Path on the remote server to treat as the root of this mount.
- `rise` (Number) Number of consecutive successes before considering the backend healthy.

### Read-Only

- `id` (Number) Unique identifier for this backend.
- `status` (String) Status of this backend.
- `undergoing_maintenance` (Boolean) True if this backend is undergoing maintenance.

## Import

Import is supported using the following syntax:

```shell
# Remote Mount Backends can be imported by specifying the id.
terraform import files_remote_mount_backend.example_remote_mount_backend 1
```
