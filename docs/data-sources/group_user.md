---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "files_group_user Data Source - files"
subcategory: ""
description: |-
  A GroupUser describes the membership of a User within a Group.
  Creating GroupUsers
  GroupUsers can be created via the normal create action. When using the update action, if the
  GroupUser record does not exist for the given user/group IDs it will be created.
---

# files_group_user (Data Source)

A GroupUser describes the membership of a User within a Group.



## Creating GroupUsers

GroupUsers can be created via the normal `create` action. When using the `update` action, if the

GroupUser record does not exist for the given user/group IDs it will be created.

## Example Usage

```terraform
data "files_group_user" "example_group_user" {
  id       = 1
  group_id = 1
  user_id  = 1
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (Number) Group User ID.

### Read-Only

- `admin` (Boolean) Is this user an administrator of this group?
- `group_id` (Number) Group ID
- `group_name` (String) Group name
- `user_id` (Number) User ID
- `usernames` (List of String) A list of usernames for users in this group
