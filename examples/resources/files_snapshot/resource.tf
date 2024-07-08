resource "files_snapshot" "example_snapshot" {
  expires_at = "2000-01-01T01:00:00Z"
  name       = "My Snapshot"
}
