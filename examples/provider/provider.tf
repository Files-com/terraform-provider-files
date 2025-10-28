terraform {
  required_providers {
    files = {
      source = "Files-com/files"
      version = "X.Y.Z"
    }
  }
}

provider "files" {
  api_key = var.files_api_key
}

resource "files_folder" "example_folder" {
  path            = "public/photos"
  mkdir_parents   = false
  provided_mtime  = "2000-01-01T01:00:00Z"
  custom_metadata = {
    key = "value"
  }
  priority_color  = "red"
}

resource "files_behavior" "example_serve_publicly_behavior" {
  path     = files_folder.example_folder.path
  behavior = "serve_publicly"
  value    = {
    key            = "public-photos"
    show_index     = true
    force_download = true
    cors_enabled   = false
  }
}

