resource "files_gpg_key" "example_gpg_key" {
  user_id              = 1
  public_key           = "7f8bc1210b09b9ddf469e6b6b8920e76"
  private_key          = "ab236cfe4a195f0226bc2e674afdd6b0"
  private_key_password = "[your GPG private key password]"
  name                 = "key name"
  generate_expires_at  = "2025-06-19 12:00:00"
  generate_keypair     = false
  generate_full_name   = "John Doe"
  generate_email       = "jdoe@example.com"
}

