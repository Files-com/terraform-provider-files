---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "files_api_key Resource - files"
subcategory: ""
description: |-
  An APIKey is a key that allows programmatic access to your Site.
  API keys confer all the permissions of the user who owns them.
  If an API key is created without a user owner, it is considered a site-wide API key, which has full permissions to do anything on the Site.
  We recommend registering API keys to service users wherever possible and then using User or Group Permissions to restrict that API Key appropriately.
---

# files_api_key (Resource)

An APIKey is a key that allows programmatic access to your Site.



API keys confer all the permissions of the user who owns them.

If an API key is created without a user owner, it is considered a site-wide API key, which has full permissions to do anything on the Site.



We recommend registering API keys to service users wherever possible and then using User or Group Permissions to restrict that API Key appropriately.

## Example Usage

```terraform
resource "files_api_key" "example_api_key" {
  user_id        = 1
  description    = "example"
  expires_at     = "2000-01-01T01:00:00Z"
  permission_set = "full"
  name           = "My Main API Key"
  path           = "shared/docs"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Internal name for the API Key.  For your use.

### Optional

- `description` (String) User-supplied description of API key.
- `expires_at` (String) API Key expiration date
- `path` (String) Folder path restriction for `office_integration` permission set API keys.
- `permission_set` (String) Permissions for this API Key. It must be full for site-wide API Keys.  Keys with the `desktop_app` permission set only have the ability to do the functions provided in our Desktop App (File and Share Link operations). Keys with the `office_integration` permission set are auto generated, and automatically expire, to allow users to interact with office integration platforms. Additional permission sets may become available in the future, such as for a Site Admin to give a key with no administrator privileges.  If you have ideas for permission sets, please let us know.
- `user_id` (Number) User ID for the owner of this API Key.  May be blank for Site-wide API Keys.

### Read-Only

- `created_at` (String) Time which API Key was created
- `descriptive_label` (String) Unique label that describes this API key.  Useful for external systems where you may have API keys from multiple accounts and want a human-readable label for each key.
- `id` (Number) API Key ID
- `key` (String) API Key actual key string
- `last_use_at` (String) API Key last used - note this value is only updated once per 3 hour period, so the 'actual' time of last use may be up to 3 hours later than this timestamp.
- `platform` (String) If this API key represents a Desktop app, what platform was it created on?
- `url` (String) URL for API host.

## Import

Import is supported using the following syntax:

```shell
# Api Keys can be imported by specifying the id.
terraform import files_api_key.example_api_key 1
```
