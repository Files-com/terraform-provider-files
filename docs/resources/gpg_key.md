---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "files_gpg_key Resource - files"
subcategory: ""
description: |-
  A GPGKey object on Files.com is used to securely store both the private and public keys associated with a GPG (GNU Privacy Guard) encryption key pair. This object enables the encryption and decryption of data using GPG, allowing you to protect sensitive information.
  The private key is kept confidential and is used for decrypting data or signing messages to prove authenticity, while the public key is used to encrypt messages that only the owner of the private key can decrypt.
  By storing both keys together in a GPGKey object, Files.com makes it easier to understand encryption operations, ensuring secure and efficient handling of encrypted data within the platform.
---

# files_gpg_key (Resource)

A GPGKey object on Files.com is used to securely store both the private and public keys associated with a GPG (GNU Privacy Guard) encryption key pair. This object enables the encryption and decryption of data using GPG, allowing you to protect sensitive information.



The private key is kept confidential and is used for decrypting data or signing messages to prove authenticity, while the public key is used to encrypt messages that only the owner of the private key can decrypt.



By storing both keys together in a GPGKey object, Files.com makes it easier to understand encryption operations, ensuring secure and efficient handling of encrypted data within the platform.

## Example Usage

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Your GPG key name.

### Optional

- `generate_email` (String) Email address of the key owner. Used for the generation of the key. Will be ignored if `generate_keypair` is false.
- `generate_expires_at` (String) Expiration date of the key. Used for the generation of the key. Will be ignored if `generate_keypair` is false.
- `generate_full_name` (String) Full name of the key owner. Used for the generation of the key. Will be ignored if `generate_keypair` is false.
- `generate_keypair` (Boolean) If true, generate a new GPG key pair. Can not be used with `public_key`/`private_key`
- `private_key` (String) Your GPG private key.
- `private_key_password` (String) Your GPG private key password. Only required for password protected keys.
- `public_key` (String) Your GPG public key
- `user_id` (Number) GPG owner's user id

### Read-Only

- `expires_at` (String) Your GPG key expiration date.
- `id` (Number) Your GPG key ID.
- `private_key_hash` (String)
- `private_key_md5` (String) MD5 hash of your GPG private key.
- `private_key_password_hash` (String)
- `public_key_hash` (String)
- `public_key_md5` (String) MD5 hash of your GPG public key

## Import

Import is supported using the following syntax:

```shell
# Gpg Keys can be imported by specifying the id.
terraform import files_gpg_key.example_gpg_key 1
```
