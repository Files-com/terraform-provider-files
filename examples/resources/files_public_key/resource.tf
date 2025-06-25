resource "files_public_key" "example_public_key" {
  user_id                       = 1
  title                         = "My Main Key"
  public_key                    = "example"
  generate_keypair              = false
  generate_private_key_password = "[your private key password]"
  generate_algorithm            = "rsa"
  generate_length               = 4096
}

