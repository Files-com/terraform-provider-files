---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "files_request Data Source - files"
subcategory: ""
description: |-
  A Request is a file that should be uploaded by a specific user or group.
  Requests can either be manually created and managed, or managed automatically by an Automation.
---

# files_request (Data Source)

A Request is a file that *should* be uploaded by a specific user or group.



Requests can either be manually created and managed, or managed automatically by an Automation.

## Example Usage

```terraform
data "files_request" "example_request" {
  id = 1
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (Number) Request ID

### Read-Only

- `automation_id` (Number) ID of automation that created request
- `destination` (String) Destination filename
- `path` (String) Folder path. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.
- `source` (String) Source filename, if applicable
- `user_display_name` (String) User making the request (if applicable)
