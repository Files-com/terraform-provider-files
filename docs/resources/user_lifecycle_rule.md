---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "files_user_lifecycle_rule Resource - files"
subcategory: ""
description: |-
  A UserLifecycleRule represents a rule that applies to users based on their inactivity, state and authentication method.
  The rule either disable or delete users who have been inactive or disabled for a specified number of days.
  The authentication_method property specifies the authentication method for the rule, which can be set to "all" or other specific methods.
  The rule can also include or exclude site and folder admins from the action.
---

# files_user_lifecycle_rule (Resource)

A UserLifecycleRule represents a rule that applies to users based on their inactivity, state and authentication method.



The rule either disable or delete users who have been inactive or disabled for a specified number of days.



The authentication_method property specifies the authentication method for the rule, which can be set to "all" or other specific methods.



The rule can also include or exclude site and folder admins from the action.

## Example Usage

```terraform
resource "files_user_lifecycle_rule" "example_user_lifecycle_rule" {
  authentication_method = "password"
  inactivity_days       = 12
  include_site_admins   = true
  include_folder_admins = true
  user_state            = "inactive"
  name                  = "password specific rules"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `action` (String) Action to take on inactive users (disable or delete)
- `authentication_method` (String) User authentication method for the rule
- `inactivity_days` (Number) Number of days of inactivity before the rule applies
- `include_folder_admins` (Boolean) Include folder admins in the rule
- `include_site_admins` (Boolean) Include site admins in the rule
- `name` (String) User Lifecycle Rule name
- `user_state` (String) State of the users to apply the rule to (inactive or disabled)

### Read-Only

- `id` (Number) User Lifecycle Rule ID
- `site_id` (Number) Site ID

## Import

Import is supported using the following syntax:

```shell
# User Lifecycle Rules can be imported by specifying the id.
terraform import files_user_lifecycle_rule.example_user_lifecycle_rule 1
```
