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

Files.com preserves the spelling of file and folder paths while comparing them using shared case and Unicode rules. Use the SDK comparison helpers when matching paths locally.
<div></div>

### Capitalization

Files.com uses case-insensitive path matching based on its fixed Unicode comparison map.

For example, the following paths have the same comparison key:

| Path Variant                          | Comparison Key              |
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

Use `/` between folder and file names, without leading or trailing slashes. SDK normalization helpers convert backslashes to `/`, remove duplicate separators, and discard exact `.` and `..` components. Discarding `..` leaves the preceding folder name intact.

| Input | Normalized path |
|-------|-----------------|
| `folder/subfolder/file.txt` | `folder/subfolder/file.txt` |
| `/folder/subfolder/file.txt` | `folder/subfolder/file.txt` |
| `folder/subfolder/file.txt/` | `folder/subfolder/file.txt` |
| `//folder//file.txt` | `folder/file.txt` |
| `folder/../file.txt` | `folder/file.txt` |

<div></div>

### Unicode and Path Comparison

Files.com compares paths using a fixed mapping shared by the server and SDKs. It treats case and many accent differences as equivalent: `Résumé.txt` and `resume.txt` identify the same file, as do `q` followed by a combining acute accent and `q`. The mapping also handles other equivalences, such as Hiragana and Katakana. Lowercasing or applying a standard Unicode normalization form alone does not reproduce these rules.

SDK comparison helpers normalize path separators and dot segments, then apply the bundled [versioned comparison map](https://github.com/Files-com/files-sdk-javascript/blob/master/shared/path_comparison.json). The [shared examples](https://github.com/Files-com/files-sdk-javascript/blob/master/shared/comparison_examples.json) give exact comparison results for integrations that implement their own matching. The map uses hexadecimal Unicode scalar values as keys: a missing entry preserves the character, an empty replacement removes it, and other replacements may contain several characters. Apply each replacement once without normalizing or lowercasing the result again.

Use comparison results only for matching. Send the original path spelling in API requests and preserve it for display and local filenames; comparison results can have a different spelling or length.

Trailing whitespace is significant for comparison. `report.txt` and `report.txt ` are different file paths, and SDK helpers preserve spaces, tabs, and newlines. Folder names cannot end in whitespace. See [Unicode Normalization](https://www.files.com/docs/files-and-folders/file-system-semantics/unicode-normalization) for the complete path rules.

<div></div>

## Workspaces

A Workspace groups files, users, groups, Partners, integrations, and workflows within a Files.com Site. An integration can provision a Workspace for a department or project and delegate its operation to a team without making that team Site Administrators. Every Site has a Default Workspace, with ID `0`; additional Workspaces have their own IDs and root folders.

Account membership, request context, and permission grants serve different purposes. Creating an account in a Workspace determines where it belongs. Selecting a Workspace determines which resources a request operates on. A permission grant determines what the caller can do there. Selecting a Workspace never grants access to it.

### Accounts and Administrative Access

A user's or group's `workspace_id` identifies the Workspace the account belongs to. Accounts belonging to a Custom Workspace stay within it. Default Workspace users and groups can receive permissions in one or more Custom Workspaces while keeping their existing accounts in Workspace `0`.

| Account | Workspace Administrator assignment | Scope |
| --- | --- | --- |
| User belonging to a Custom Workspace | Set the user's `workspace_admin` to `true`. | That user's own Custom Workspace. |
| Default Workspace user | Create an `admin` Permission for the user on a Custom Workspace's root folder. | Each Custom Workspace with a root grant. |
| Default Workspace group | Create an `admin` Permission for the group on a Custom Workspace's root folder. | Every member inherits administration of each Workspace with a root grant. |

`workspace_admin` is not a summary of a user's effective administrative access. A Default Workspace user can administer a Custom Workspace through a direct or group root grant while their `workspace_admin` remains `false`. Groups have no `workspace_admin` field. See [Users](/terraform/resources/user-accounts/users) and [Groups](/terraform/resources/user-accounts/groups) for account fields.

An `admin` grant on the **Custom Workspace root** provides full Workspace Administrator authority over its files, users, groups, Partners, workflows, and integrations. An `admin` grant on a subfolder provides Folder Admin authority over that folder and its descendants; it does not provide Workspace administration. Other permission levels provide their corresponding folder access without Workspace administration. [Permissions](/terraform/resources/user-accounts/permissions) defines the levels.

Site Administrators manage cross-Workspace assignments to Default Workspace accounts. Workspace Administrators manage accounts and permissions within their own scope. Site Administrators retain access to every Workspace; adding a Workspace grant does not narrow Site Administrator authority. The [product documentation](https://www.files.com/docs/workspaces/workspace-administrators) explains the administrator's operational scope and site-wide controls.

### Request Context and API Keys

You can include the `X-Files-Workspace-Id` REST header to select a Workspace for a request. SDK request options and CLI configuration send that same selection. When a Workspace is selected, Workspace-scoped resources are listed, created, and changed within that context, and ordinary paths are relative to its root.

A resource's `workspace_id` request field describes the resource's Workspace membership. It is separate from the SDK's Workspace request option or REST header. Creating a Workspace-scoped resource in a Custom Workspace defaults its `workspace_id` to the selected Workspace; a mismatching membership value is rejected with `not-authorized/insufficient-permission-for-params`.

Selecting another Workspace with an API key requires a **Full Access key created in the Default Workspace**. A user key follows that user's current access, including group permissions. A site-wide Full Access key created in the Default Workspace has Site Administrator authority in every Workspace. A Files Only key stays in its creation Workspace, even if its user has cross-Workspace access. Any key created in a Custom Workspace stays within that Workspace. Selecting another context with these confined keys is rejected with `bad-request/invalid-workspace-id-header`.

An account belonging to a Custom Workspace is scoped there when it authenticates normally. For a Default Workspace user, explicitly select the intended Workspace for an integration rather than relying on an interactive login preference. [API Keys](/terraform/resources/developers/api-keys) and [Authentication](/terraform/overview/authentication) cover credentials.

The adjacent request uses the member's own Default Workspace Full Access user key to list the root of Workspace `123`, after access has been assigned.

```shell title="Request as the group member"
curl 'https://YOUR_SUBDOMAIN.files.com/api/rest/v1/folders/' \
  -H 'X-FilesAPI-Key: YOUR_MEMBER_API_KEY' \
  -H 'X-Files-Workspace-Id: 123'
```

### Delegating a Workspace to an Existing Group

An operations team already represented by a Default Workspace group can administer a Custom Workspace through one root Permission. The group and its members stay in the Default Workspace, so the same team can receive different access in other Workspaces.

First retrieve the target [Workspace](/terraform/resources/settings/workspaces) and [Group](/terraform/resources/user-accounts/groups) IDs as a Site Administrator in Workspace `0`. The examples use Workspace `123`, group `456`, and member user `789`; replace them with your own IDs. Confirm that the group belongs to Workspace `0` and that the intended user is a member.

Create the Permission using a Default Workspace Full Access site-wide key or a Full Access user key belonging to a Site Administrator. Keep the request context at `0` and use the qualified root path `_/Workspaces/123`. Set `group_id` to the group's ID, `permission` to `admin`, and `recursive` to `true`. Save the returned Permission `id` for later removal. For an individual Default Workspace user, use `user_id` instead of `group_id`.

For a Default Workspace group, a Site Administrator can also select Workspace `123` and use an empty `path` to grant access to its root. The qualified path in Workspace `0` works for both Default Workspace users and groups and keeps the account scope and target Workspace explicit. Appending a subfolder to the path would grant Folder Admin access instead of Workspace Administrator authority.

After the grant, run the request-context example above with the member's own credential and Workspace `123` selected. That member can work with the Workspace's files and perform Workspace Administrator operations, such as managing its users, Partners, and integrations. A Site Administrator's successful request does not establish that the member has the intended access.

```shell title="Find account and Workspace IDs"
curl 'https://YOUR_SUBDOMAIN.files.com/api/rest/v1/workspaces' \
  -H 'X-FilesAPI-Key: YOUR_SITE_ADMIN_API_KEY' \
  -H 'X-Files-Workspace-Id: 0'
curl 'https://YOUR_SUBDOMAIN.files.com/api/rest/v1/groups' \
  -H 'X-FilesAPI-Key: YOUR_SITE_ADMIN_API_KEY' \
  -H 'X-Files-Workspace-Id: 0'
```

```shell title="Grant group administration"
curl 'https://YOUR_SUBDOMAIN.files.com/api/rest/v1/permissions' \
  -X POST \
  -H 'X-FilesAPI-Key: YOUR_SITE_ADMIN_API_KEY' \
  -H 'X-Files-Workspace-Id: 0' \
  -H 'Content-Type: application/json' \
  -d '{"path":"_/Workspaces/123","group_id":456,"permission":"admin","recursive":true}'
```

### Permission Inspection and Removal

List the member's Permissions with `user_id` and `include_groups=true` to include grants inherited through group membership. Listing only direct user grants can miss the Permission that provides Workspace administration. In Workspace `0`, the Custom Workspace root appears as `_/Workspaces/123`; in Workspace `123`, paths are relative to that root. Inspect the root path and `permission=admin`, rather than treating the user's `workspace_admin` field as their effective administrative access.

Permission lists show individual grants, rather than a single flag for effective administrative access. Membership in several groups combines their access. A Permission using `group_ids` instead of `group_id` requires membership in all the specified groups; it is not a shorthand for assigning the same grant to several independent groups.

Removing a member ends access received through that group. Deleting the root Permission ends the group's Workspace Administrator grant for every member. These changes leave independent direct and other group grants in place, so review all applicable grants when withdrawing access. Default Workspace user API keys follow those permission changes without being recreated.

Group membership maintained through SCIM follows the same rule. A Group Admin allowed to add members can give those users the group's existing Workspace Administrator access. Choose who manages the group with that authority in mind.

Delete the Permission by its returned `id` as the Site Administrator in Workspace `0`. The removal examples use Permission ID `9001`; replace it with the ID returned by your create request. Permissions are created and deleted, rather than updated in place. If narrower folder access is still needed, assign it explicitly; deleting a broad grant does not restore narrower grants it previously replaced.

```shell title="Inspect member grants and remove the group grant"
curl 'https://YOUR_SUBDOMAIN.files.com/api/rest/v1/permissions?user_id=789&include_groups=true' \
  -H 'X-FilesAPI-Key: YOUR_SITE_ADMIN_API_KEY' \
  -H 'X-Files-Workspace-Id: 0'
curl 'https://YOUR_SUBDOMAIN.files.com/api/rest/v1/permissions/9001' \
  -X DELETE \
  -H 'X-FilesAPI-Key: YOUR_SITE_ADMIN_API_KEY' \
  -H 'X-Files-Workspace-Id: 0'
```

## Foreign Language Support

The Files.com Terraform Provider will soon be updated to support localized responses by using a configuration
method. When available, it can be used to guide the API in selecting a preferred language for applicable response content.

Language support currently applies to select human-facing fields only, such as notification messages
and error descriptions.

If the specified language is not supported or the value is omitted, the API defaults to English.

## Errors

## Mock Server

Files.com publishes a mock Files.com API server, which is useful for testing your use of the Files.com
SDKs and other direct integrations against the Files.com API in an integration test environment.
It never checks credentials: send any placeholder API key, and never use real Files.com credentials
with it.

The server has two modes, chosen when it starts:

* **Legacy mode** (the default) checks required parameters and parameter types, then returns a fixed
  example response for each API endpoint. It does not maintain state and it does not deeply inspect
  your submissions for correctness, which makes it useful for testing basic network operations and
  JSON encoding for your SDK or API client.
* **Simulation mode** keeps records, files and folders in memory, so a test can create, list, update
  and delete resources, upload a file and download the same bytes, and make chosen requests fail,
  stall or lose their connection on purpose. Requests it does not simulate fail with a clear error
  instead of returning an example response.

Start the server from its source with Ruby and Bundler. `FILES_MOCK_MODE` is read once at startup:
leaving it unset or setting it to `legacy` starts legacy mode, `simulation` starts simulation mode,
and any other value stops startup with an error. Legacy mode listens on port 4041 on all IPv4
interfaces, and simulation mode on `127.0.0.1:4041`.

Simulation mode keeps its state only in the server process. Its control endpoints under
`/__files_mock/v1` report when the server is ready, reset it with your fixtures, add fault rules
and return the journal of the requests it received. It refuses work over its limits instead of
truncating it. The README in the source describes all of these, the operations it simulates and
how to configure its limits.

Download the server as a Docker image via [Docker Hub](https://hub.docker.com/r/filescom/files-mock-server).
The image's `latest` tag moves to whichever server was published last, so it need not include
simulation mode; to run exactly the server the README describes, build the image from the source.

The Source Code is also available on [GitHub](https://github.com/Files-com/files-mock-server).

```shell title="Start the Mock Server"
bundle install

## Legacy mode
bundle exec puma

## Simulation mode
FILES_MOCK_MODE=simulation bundle exec puma
```

## Usage

The macOS provider binaries require macOS 13 or later. Building from source
requires Go 1.27.1 or later.

See the [docs](./docs) directory for Resource and Data Source examples and documentation.

## Debugging

This provider uses the standard Terraform debugging methods. For more information, please refer to the [Terraform Debugging](https://www.terraform.io/docs/internals/debugging.html) documentation.

## Getting Support

The Files.com team is happy to help with any SDK Integration challenges you may face.

Just email <support@files.com> and we'll get the process started.
