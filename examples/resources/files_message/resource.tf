resource "files_message" "example_message" {
  user_id    = 1
  project_id = 1
  subject    = "Files.com Account Upgrade"
  body       = "We should upgrade our Files.com account!"
}
