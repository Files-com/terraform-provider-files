---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "files_gpg_key Data Source - files"
subcategory: ""
description: |-
  GPG keys for decrypt or encrypt behaviors.
---

# files_gpg_key (Data Source)

GPG keys for decrypt or encrypt behaviors.

## Example Usage

```terraform
data "files_gpg_key" "example_gpg_key" {
  id = 1
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (Number) Your GPG key ID.

### Read-Only

- `expires_at` (String) Your GPG key expiration date.
- `name` (String) Your GPG key name.
- `private_key` (String) Your GPG private key.
- `private_key_hash` (String)
- `private_key_password` (String) Your GPG private key password. Only required for password protected keys.
- `private_key_password_hash` (String)
- `public_key` (String) Your GPG public key
- `public_key_hash` (String)
- `user_id` (Number) GPG owner's user id
