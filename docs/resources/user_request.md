---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "files_user_request Resource - files"
subcategory: ""
description: |-
  A UserRequest is an operation that allows anonymous users to place a request for access on the login screen to the site administrator.
---

# files_user_request (Resource)

A UserRequest is an operation that allows anonymous users to place a request for access on the login screen to the site administrator.

## Example Usage

```terraform
resource "files_user_request" "example_user_request" {
  name    = "name"
  email   = "email"
  details = "details"
  company = "Acme Inc."
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `details` (String) Details of the user's request
- `email` (String) User email address
- `name` (String) User's full name

### Optional

- `company` (String) User's company name

### Read-Only

- `id` (Number) ID

## Import

Import is supported using the following syntax:

```shell
# User Requests can be imported by specifying the id.
terraform import files_user_request.example_user_request 1
```
