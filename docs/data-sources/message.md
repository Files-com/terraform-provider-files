---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "files_message Data Source - files"
subcategory: ""
description: |-
  A Messages is a part of Files.com's project management features and represent a message posted by a user to a project.
---

# files_message (Data Source)

A Messages is a part of Files.com's project management features and represent a message posted by a user to a project.

## Example Usage

```terraform
data "files_message" "example_message" {
  id = 1
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (Number) Message ID

### Read-Only

- `body` (String) Message body.
- `comments` (Dynamic) Comments.
- `subject` (String) Message subject.
