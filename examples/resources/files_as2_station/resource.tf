resource "files_as2_station" "example_as2_station" {
  name               = "AS2 Station Name"
  workspace_id       = 1
  public_certificate = "public_certificate"
  private_key        = "private_key"
}

