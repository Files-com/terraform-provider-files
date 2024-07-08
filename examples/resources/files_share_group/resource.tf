resource "files_share_group" "example_share_group" {
  user_id = 1
  notes   = "This group is defined for testing purposes"
  name    = "Test group 1"
  members = [
    {
      name    = "John Doe"
      company = "Acme Ltd"
      email   = "johndoe@gmail.com"
    }
  ]
}
