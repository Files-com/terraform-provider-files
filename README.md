# Files.com Terraform Provider

The Files.com Terraform Provider provides convenient access to the Files.com API for managing your Files.com account.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.13+

## Installation

Require the provider in your Terraform configuration:

```terraform
terraform {
  required_providers {
    files = {
      source = "Files-com/files"
      version = "0.1.73"
    }
  }
}
```

## Setting API Key

### Setting by ENV

```sh
export FILES_API_KEY="XXXX-XXXX..."
```

### Set by Provider Configuration

```terraform
provider "files" {
  api_key = "XXXX-XXXX..."
}
```

## Provider Configuration Options

### Endpoint Override

Set client to use your site subdomain if your site is configured to disable global acceleration.
Otherwise, don't change this setting for production. For dev/CI, you can point this to the mock server.

```terraform
provider "files" {
  endpoint_override = "https://SUBDOMAIN.files.com"
}
```

## Usage

See the [docs](./docs) directory for Resource and Data Source examples and documentation.

## Debugging

This provider uses the standard Terraform debugging methods. For more information, please refer to the [Terraform Debugging](https://www.terraform.io/docs/internals/debugging.html) documentation.

## Getting Support

The Files.com team is happy to help with any SDK Integration challenges you may face.

Just email <support@files.com> and we'll get the process started.
