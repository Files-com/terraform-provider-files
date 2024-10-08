---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "files_user_request Data Source - files"
subcategory: ""
description: |-
  A UserRequest is an operation that allows anonymous users to place a request for access on the login screen to the site administrator.
---

# files_user_request (Data Source)

A UserRequest is an operation that allows anonymous users to place a request for access on the login screen to the site administrator.

## Example Usage

```terraform
data "files_user_request" "example_user_request" {
  id = 1
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (Number) ID

### Read-Only

- `company` (String) User's company name
- `details` (String) Details of the user's request
- `email` (String) User email address
- `name` (String) User's full name
