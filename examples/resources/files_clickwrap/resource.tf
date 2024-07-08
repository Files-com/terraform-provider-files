resource "files_clickwrap" "example_clickwrap" {
  name             = "Example Site NDA for Files.com Use"
  body             = "[Legal body text]"
  use_with_bundles = "example"
  use_with_inboxes = "example"
  use_with_users   = "example"
}
