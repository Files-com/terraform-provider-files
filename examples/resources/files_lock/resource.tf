resource "files_lock" "example_lock" {
  path                     = "path"
  allow_access_by_any_user = true
  exclusive                = true
  recursive                = true
  timeout                  = 1
}

