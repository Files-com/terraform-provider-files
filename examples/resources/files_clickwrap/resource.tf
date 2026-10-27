resource "files_clickwrap" "example_clickwrap" {
  name             = "Example Site NDA for Files.com Use"
  body             = "[Legal body text]"
  use_with_bundles = "none"
  use_with_inboxes = "none"
  use_with_users   = "none"
}

