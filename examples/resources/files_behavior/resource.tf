resource "files_behavior" "example_behavior" {
  value                          = {
    method = "GET"
  }
  disable_parent_folder_behavior = false
  recursive                      = false
  name                           = "example"
  description                    = "example"
  path                           = "path"
  behavior                       = "webhook"
}

resource "files_behavior" "example_webhook_behavior" {
  path     = "path"
  behavior = "webhook"
  value    = {
    urls                 = ["https://mysite.com/url..."]
    method               = "POST"
    triggers             = ["create", "read", "update", "destroy", "move", "copy"]
    triggering_filenames = ["*.pdf", "*so*.jpg"]
    exclude_filenames    = ["*.txt", "*wo*.png"]
    encoding             = "RAW"
    headers              = {
      MY-HEADER = "foo"
    }
    body                 = {
      MY_BODY_PARAM = "bar"
    }
    verification_token   = "tok12345"
    file_form_field      = "my_form_field"
    file_as_body         = "my_file_body"
    use_dedicated_ips    = false
  }
}

resource "files_behavior" "example_file_expiration_behavior" {
  path     = "path"
  behavior = "file_expiration"
  value    = 30
}

resource "files_behavior" "example_auto_encrypt_behavior" {
  path     = "path"
  behavior = "auto_encrypt"
  value    = {
    gpg_key_id         = 1
    gpg_key_ids        = [1]
    algorithm          = "PGP/GPG"
    signing_key_id     = 1
    suffix             = ".gpg"
    armor              = false
    gpg_key_partner_id = 1
  }
}

resource "files_behavior" "example_lock_subfolders_behavior" {
  path     = "path"
  behavior = "lock_subfolders"
  value    = {
    level = "children_recursive"
  }
}

resource "files_behavior" "example_storage_region_behavior" {
  path     = "path"
  behavior = "storage_region"
  value    = "us-east-1"
}

resource "files_behavior" "example_serve_publicly_behavior" {
  path     = "path"
  behavior = "serve_publicly"
  value    = {
    key            = "public-photos"
    show_index     = true
    force_download = true
    cors_enabled   = false
  }
}

resource "files_behavior" "example_create_user_folders_behavior" {
  path     = "path"
  behavior = "create_user_folders"
  value    = {
    permission            = "full"
    additional_permission = "bundle"
    existing_users        = true
    group_id              = 1
    new_folder_name       = "username"
    subfolders            = ["in", "out"]
  }
}

resource "files_behavior" "example_inbox_behavior" {
  path     = "path"
  behavior = "inbox"
  value    = {
    key                                            = "application-forms"
    dont_separate_submissions_by_folder            = true
    dont_allow_folders_in_uploads                  = false
    require_inbox_recipient                        = false
    show_on_login_page                             = true
    title                                          = "Submit Your Job Applications Here"
    description                                    = "Thanks for coming to the Files.com Job Application Page"
    help_text                                      = "If you have trouble here, please contact your recruiter."
    require_registration                           = true
    password                                       = "foobar"
    path_template                                  = "{{name}}_{{ip}}"
    path_template_time_zone                        = "Eastern Time (US & Canada)"
    enable_inbound_email_address                   = true
    notify_senders_on_successful_uploads_via_email = true
    notify_senders_on_successful_uploads_via_web   = true
    allow_whitelisting                             = true
    whitelist                                      = ["john@test.com", "mydomain.com"]
    disable_web_upload                             = true
    capture_email_body_filename                    = "_body.txt"
  }
}

resource "files_behavior" "example_limit_file_extensions_behavior" {
  path     = "path"
  behavior = "limit_file_extensions"
  value    = {
    extensions = ["xls", "csv"]
    mode       = "whitelist"
  }
}

resource "files_behavior" "example_limit_file_regex_behavior" {
  path     = "path"
  behavior = "limit_file_regex"
  value    = ["/Document-.*/"]
}

resource "files_behavior" "example_amazon_sns_behavior" {
  path     = "path"
  behavior = "amazon_sns"
  value    = {
    arns            = ["ARN"]
    triggers        = ["create", "read", "update", "destroy", "move", "copy"]
    aws_credentials = {
      access_key_id     = "ACCESS_KEY_ID"
      region            = "us-east-1"
      secret_access_key = "SECRET_ACCESS_KEY"
    }
  }
}

resource "files_behavior" "example_watermark_behavior" {
  path     = "path"
  behavior = "watermark"
  value    = {
    gravity             = "SouthWest"
    max_height_or_width = 20
    transparency        = 25
    dynamic_text        = "Confidential: For use by {{user}} only."
  }
}

resource "files_behavior" "example_remote_server_mount_behavior" {
  path     = "path"
  behavior = "remote_server_mount"
  value    = {
    remote_server_id = 1
    remote_path      = ""
  }
}

resource "files_behavior" "example_slack_webhook_behavior" {
  path     = "path"
  behavior = "slack_webhook"
  value    = {
    url        = "https://mysite.com/url..."
    username   = "Files.com"
    channel    = "alerts"
    icon_emoji = ":robot_face:"
    triggers   = ["create", "read", "update", "destroy", "move", "copy"]
  }
}

resource "files_behavior" "example_auto_decrypt_behavior" {
  path     = "path"
  behavior = "auto_decrypt"
  value    = {
    gpg_key_id         = 1
    gpg_key_ids        = [1]
    algorithm          = "PGP/GPG"
    suffix             = ".gpg"
    ignore_mdc_error   = true
    gpg_key_partner_id = 1
  }
}

resource "files_behavior" "example_override_upload_filename_behavior" {
  path     = "path"
  behavior = "override_upload_filename"
  value    = {
    filename_override_pattern = "%Fb_addition5%Fe"
    time_zone                 = "Eastern Time (US & Canada)"
  }
}

resource "files_behavior" "example_permission_fence_behavior" {
  path     = "path"
  behavior = "permission_fence"
  value    = {
    fenced_permissions = "all"
  }
}

resource "files_behavior" "example_limit_filename_length_behavior" {
  path     = "path"
  behavior = "limit_filename_length"
  value    = {
    max_length = 30
    shorten    = true
  }
}

resource "files_behavior" "example_organize_files_into_subfolders_behavior" {
  path     = "path"
  behavior = "organize_files_into_subfolders"
  value    = {
    subfolder_name_type = "regex, extension, created_at, :provided_modified_at"
    regex               = "(?<=\\-)(.*?)(?=\\.)"
    strftime_format     = "%Y-%m-%d"
    time_zone           = "Eastern Time (US & Canada)"
    apply_behavior      = true
  }
}

resource "files_behavior" "example_teams_webhook_behavior" {
  path     = "path"
  behavior = "teams_webhook"
  value    = {
    url      = "https://mysite.com/url..."
    triggers = ["create", "read", "update", "destroy", "move", "copy"]
  }
}

resource "files_behavior" "example_google_pub_sub_behavior" {
  path     = "path"
  behavior = "google_pub_sub"
  value    = {
    projects_topics    = [
      {
        project_id = "my-project-id"
        topic_id   = "my-topic-id"
      }
    ]
    triggers           = ["create", "read", "update", "destroy", "move", "copy"]
    google_credentials = {
      type                        = "service_account"
      project_id                  = "your-project-id"
      private_key_id              = "your-private-key-id"
      private_key                 = "-----BEGIN PRIVATE KEY-----\\nMIIC..."
      client_email                = "your-service-account@your-project-id.iam.gserviceaccount.com"
      client_id                   = "your-client-id"
      auth_uri                    = "https=>//accounts.google.com/o/oauth2/auth"
      token_uri                   = "https=>//oauth2.googleapis.com/token"
      auth_provider_x509_cert_url = "https://www.googleapis.com/oauth2/v1/certs"
      client_x509_cert_url        = "https://www.googleapis.com/robot/v1/metadata/x509/your-service-account%40your-project-id.iam.gserviceaccount.com"
    }
  }
}

resource "files_behavior" "example_archive_overwritten_or_deleted_files_behavior" {
  path     = "path"
  behavior = "archive_overwritten_or_deleted_files"
  value    = {
    filename_override_pattern = "%Fb_addition5%Fe"
    time_zone                 = "Eastern Time (US & Canada)"
    archive_path              = "/Archive"
  }
}

resource "files_behavior" "example_auto_recrypt_behavior" {
  path     = "path"
  behavior = "auto_recrypt"
  value    = {
    decrypt_gpg_key_ids        = [1]
    encrypt_gpg_key_ids        = [1]
    decrypt_gpg_key_partner_id = 1
    encrypt_gpg_key_partner_id = 1
    ignore_mdc_error           = true
    signing_key_id             = 1
    armor                      = false
  }
}
