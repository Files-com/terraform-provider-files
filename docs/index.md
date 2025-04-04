---
page_title: "Files.com Provider"
description: |-
  The Files.com Terraform Provider provides convenient access to the Files.com API for managing your Files.com account via Terraform.
---

# Files.com Provider

The [Files.com](https://www.files.com/) Terraform Provider provides access to the Files.com API for managing your Files.com account via Terraform or OpenTofu.

You must configure the provider with the proper credentials before you can use it.

Use the navigation to the left to read about the available Resources and Data Sources.

## Example Usage

```terraform
terraform {
  required_providers {
    files = {
      source = "Files-com/files"
      version = "X.Y.Z"
    }
  }
}

provider "files" {
  api_key = var.files_api_key
}

resource "files_folder" "example_folder" {
  path            = "public/photos"
  mkdir_parents   = false
  provided_mtime  = "2000-01-01T01:00:00Z"
  custom_metadata = {
    key = "value"
  }
  priority_color  = "red"
}

resource "files_behavior" "example_serve_publicly_behavior" {
  path     = files_folder.example_folder.path
  behavior = "serve_publicly"
  value    = {
    key            = "public-photos"
    show_index     = true
    force_download = true
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `api_key` (String, Sensitive) The API key used to authenticate with Files.com. It can also be sourced from the `FILES_API_KEY` environment variable.
- `endpoint_override` (String) Required if your site is configured to disable global acceleration. This can also be set to use a mock server in development or CI.
- `environment` (String)
- `feature_flags` (List of String)
