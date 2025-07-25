---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "files_gpg_key Data Source - files"
subcategory: ""
description: |-
  A GPGKey object on Files.com is used to securely store both the private and public keys associated with a GPG (GNU Privacy Guard) encryption key pair. This object enables the encryption and decryption of data using GPG, allowing you to protect sensitive information.
  The private key is kept confidential and is used for decrypting data or signing messages to prove authenticity, while the public key is used to encrypt messages that only the owner of the private key can decrypt.
  By storing both keys together in a GPGKey object, Files.com makes it easier to understand encryption operations, ensuring secure and efficient handling of encrypted data within the platform.
---

# files_gpg_key (Data Source)

A GPGKey object on Files.com is used to securely store both the private and public keys associated with a GPG (GNU Privacy Guard) encryption key pair. This object enables the encryption and decryption of data using GPG, allowing you to protect sensitive information.



The private key is kept confidential and is used for decrypting data or signing messages to prove authenticity, while the public key is used to encrypt messages that only the owner of the private key can decrypt.



By storing both keys together in a GPGKey object, Files.com makes it easier to understand encryption operations, ensuring secure and efficient handling of encrypted data within the platform.

## Example Usage

```terraform
data "files_gpg_key" "example_gpg_key" {
  id = 1
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (Number) Your GPG key ID.

### Read-Only

- `expires_at` (String) Your GPG key expiration date.
- `name` (String) Your GPG key name.
- `private_key` (String) Your GPG private key.
- `private_key_hash` (String)
- `private_key_md5` (String) MD5 hash of your GPG private key.
- `private_key_password` (String) Your GPG private key password. Only required for password protected keys.
- `private_key_password_hash` (String)
- `public_key` (String) Your GPG public key
- `public_key_hash` (String)
- `public_key_md5` (String) MD5 hash of your GPG public key
- `user_id` (Number) GPG owner's user id
