resource "files_lock" "example_lock" {
  path                     = "path"
  allow_access_by_any_user = false
  exclusive                = false
  recursive                = true
  timeout                  = 1
}

