---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "files_snapshot Resource - files"
subcategory: ""
description: |-
  Snapshots allow you to create a read-only archive of files at a specific point in time. You can define a snapshot, add files to it, and then finalize it. Once finalized, the snapshot’s contents are immutable.
  Each snapshot may have an expiration date. When the expiration date is reached, the snapshot is automatically deleted from the Files.com platform.
---

# files_snapshot (Resource)

Snapshots allow you to create a read-only archive of files at a specific point in time. You can define a snapshot, add files to it, and then finalize it. Once finalized, the snapshot’s contents are immutable.



Each snapshot may have an expiration date. When the expiration date is reached, the snapshot is automatically deleted from the Files.com platform.

## Example Usage

```terraform
resource "files_snapshot" "example_snapshot" {
  expires_at = "2000-01-01T01:00:00Z"
  name       = "My Snapshot"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `expires_at` (String) When the snapshot expires.
- `name` (String) A name for the snapshot.
- `paths` (List of String) An array of paths to add to the snapshot.

### Read-Only

- `bundle_id` (Number) The bundle using this snapshot, if applicable.
- `finalized_at` (String) When the snapshot was finalized.
- `id` (Number) The snapshot's unique ID.
- `user_id` (Number) The user that created this snapshot, if applicable.

## Import

Import is supported using the following syntax:

```shell
# Snapshots can be imported by specifying the id.
terraform import files_snapshot.example_snapshot 1
```
