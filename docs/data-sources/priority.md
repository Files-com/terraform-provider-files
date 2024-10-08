---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "files_priority Data Source - files"
subcategory: ""
description: |-
  A Priority is a color tag that is attached to the path.
---

# files_priority (Data Source)

A Priority is a color tag that is attached to the path.

## Example Usage

```terraform
data "files_priority" "example_priority" {
  path = "path"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `path` (String) The path corresponding to the priority color. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.

### Read-Only

- `color` (String) The priority color
