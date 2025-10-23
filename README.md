# Files.com Terraform Provider

The Files.com Terraform Provider makes it easy to manage your Files.com account using Terraform configuration files.

Files.com is the cloud-native, next-gen MFT, SFTP, and secure file-sharing platform that replaces brittle legacy servers with one always-on, secure fabric. Automate mission-critical file flows—across any cloud, protocol, or partner—while supporting human collaboration and eliminating manual work.

With universal SFTP, AS2, HTTPS, and 50+ native connectors backed by military-grade encryption, Files.com unifies governance, visibility, and compliance in a single pane of glass.

## Introduction

Terraform lets you manage your Files.com resources using Infrastructure-as-Code (IaC) processes.

Terraform users define and configure data center infrastructure using a declarative configuration language known as Terraform Configuration Language, or optionally JSON.

The Files.com Provider for Terraform allows you to automate the management of your Files.com site, including resources such as users, groups, folders, share links, inboxes, automations, security, encryption, access protocols, permissions, remote servers, remote syncs, data governance, and more.

### Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.13+

### Installation

Require the provider in your Terraform configuration:

```hcl title="Getting started"
## Specify that we want to use the Files.com Provider

terraform {
  required_providers {
    files = {
      source = "Files-com/files"
    }
  }
}

## Configure the Files.com Provider with your API Key and site URL

provider "files" {
  api_key = "YOUR_API_KEY"
  endpoint_override = "https://MYSITE.files.com"
}
```

Explore the [terraform-provider-files](https://github.com/Files-com/terraform-provider-files) code on GitHub.

## Authentication

### Authenticate with an API Key

Authenticating with an API key is the recommended authentication method for most scenarios, and is
the method used in the examples on this site.

To use the Terraform provider with an API Key, first generate an API key from the [web
interface](https://www.files.com/docs/sdk-and-apis/api-keys) or [via the API or an
SDK](/rest/resources/developers/api-keys).

Note that when using a user-specific API key, if the user is an administrator, you will have full
access to the entire API. If the user is not an administrator, you will only be able to access files
that user can access, and no access will be granted to site administration functions in the API.

```hcl title="Example Configuration"
provider "files" {
  api_key = "YOUR_API_KEY"
}
```

```shell title="Environment Variable"
## Alternatively, you can set the API key as an environment variable.
export FILES_API_KEY="YOUR_API_KEY"
```

Don't forget to replace the placeholder, `YOUR_API_KEY`, with your actual API key.

## Configuration

### Configuration Options

#### Endpoint Override

Setting the endpoint override for the API is required if your site is configured to disable global acceleration.
This can also be set to use a mock server in development or CI.

```hcl title="Example Configuration"
provider "files" {
  endpoint_override = "https://SUBDOMAIN.files.com"
}
```

## Sort and Filter

Several of the Files.com API resources have list operations that return multiple instances of the
resource. The List operations can be sorted and filtered.

### Sorting

To sort the returned data, pass in the ```sort_by``` method argument.

Each resource supports a unique set of valid sort fields and can only be sorted by one field at a
time.

#### Special note about the List Folder Endpoint

For historical reasons, and to maintain compatibility
with a variety of other cloud-based MFT and EFSS services, Folders will always be listed before Files
when listing a Folder.  This applies regardless of the sorting parameters you provide.  These *will* be
used, after the initial sort application of Folders before Files.

### Filtering

Filters apply selection criteria to the underlying query that returns the results. They can be
applied individually or combined with other filters, and the resulting data can be sorted by a
single field.

Each resource supports a unique set of valid filter fields, filter combinations, and combinations of
filters and sort fields.

#### Filter Types

| Filter | Type | Description |
| --------- | --------- | --------- |
| `filter` | Exact | Find resources that have an exact field value match to a passed in value. (i.e., FIELD_VALUE = PASS_IN_VALUE). |
| `filter_prefix` | Pattern | Find resources where the specified field is prefixed by the supplied value. This is applicable to values that are strings. |
| `filter_gt` | Range | Find resources that have a field value that is greater than the passed in value.  (i.e., FIELD_VALUE > PASS_IN_VALUE). |
| `filter_gteq` | Range | Find resources that have a field value that is greater than or equal to the passed in value.  (i.e., FIELD_VALUE >=  PASS_IN_VALUE). |
| `filter_lt` | Range | Find resources that have a field value that is less than the passed in value.  (i.e., FIELD_VALUE < PASS_IN_VALUE). |
| `filter_lteq` | Range | Find resources that have a field value that is less than or equal to the passed in value.  (i.e., FIELD_VALUE \<= PASS_IN_VALUE). |

## Paths

Working with paths in Files.com involves several important considerations. Understanding how path comparisons are applied helps developers ensure consistency and accuracy across all interactions with the platform.
<div></div>

### Capitalization

Files.com compares paths in a **case-insensitive** manner. This means path segments are treated as equivalent regardless of letter casing.

For example, all of the following resolve to the same internal path:

| Path Variant                          | Interpreted As              |
|---------------------------------------|------------------------------|
| `Documents/Reports/Q1.pdf`            | `documents/reports/q1.pdf`  |
| `documents/reports/q1.PDF`            | `documents/reports/q1.pdf`  |
| `DOCUMENTS/REPORTS/Q1.PDF`            | `documents/reports/q1.pdf`  |

This behavior applies across:
- API requests
- Folder and file lookup operations
- Automations and workflows

See also: [Case Sensitivity Documentation](https://www.files.com/docs/files-and-folders/case-sensitivity/)

### Slashes

All path parameters in Files.com (API, SDKs, CLI, automations, integrations) must **omit leading and trailing slashes**. Paths are always treated as **absolute and slash-delimited**, so only internal `/` separators are used and never at the start or end of the string.

####  Path Slash Examples
| Path                              | Valid? | Notes                         |
|-----------------------------------|--------|-------------------------------|
| `folder/subfolder/file.txt`       |   ✅   | Correct, internal separators only |
| `/folder/subfolder/file.txt`      |   ❌   | Leading slash not allowed     |
| `folder/subfolder/file.txt/`      |   ❌   | Trailing slash not allowed    |
| `//folder//file.txt`              |   ❌   | Duplicate separators not allowed |

<div></div>

### Unicode Normalization

Files.com normalizes all paths using [Unicode NFC (Normalization Form C)](https://www.unicode.org/reports/tr15/#Norm_Forms) before comparison. This ensures consistency across different representations of the same characters.

For example, the following two paths are treated as equivalent after NFC normalization:

| Input                                  | Normalized Form       |
|----------------------------------------|------------------------|
| `uploads/\u0065\u0301.txt`             | `uploads/é.txt`        |
| `docs/Café/Report.txt`                 | `docs/Café/Report.txt` |

- All input must be UTF‑8 encoded.
- Precomposed and decomposed characters are unified.
- This affects search, deduplication, and comparisons across SDKs.

<div></div>

## Foreign Language Support

The Files.com Terraform Provider will soon be updated to support localized responses by using a configuration
method. When available, it can be used to guide the API in selecting a preferred language for applicable response content.

Language support currently applies to select human-facing fields only, such as notification messages
and error descriptions.

If the specified language is not supported or the value is omitted, the API defaults to English.

## Errors

## Mock Server

Files.com publishes a Files.com API server, which is useful for testing your use of the Files.com
SDKs and other direct integrations against the Files.com API in an integration test environment.

It is a Ruby app that operates as a minimal server for the purpose of testing basic network
operations and JSON encoding for your SDK or API client. It does not maintain state and it does not
deeply inspect your submissions for correctness.

Eventually we will add more features intended for integration testing, such as the ability to
intentionally provoke errors.

Download the server as a Docker image via [Docker Hub](https://hub.docker.com/r/filescom/files-mock-server).

The Source Code is also available on [GitHub](https://github.com/Files-com/files-mock-server).

A README is available on the GitHub link.

## Usage

See the [docs](./docs) directory for Resource and Data Source examples and documentation.

## Debugging

This provider uses the standard Terraform debugging methods. For more information, please refer to the [Terraform Debugging](https://www.terraform.io/docs/internals/debugging.html) documentation.

## Getting Support

The Files.com team is happy to help with any SDK Integration challenges you may face.

Just email <support@files.com> and we'll get the process started.
