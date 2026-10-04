resource "files_lock" "example_lock" {
  path                     = "locked_file"
  token                    = "17c54824e9931a4688ca032d03f6663c"
  allow_access_by_any_user = false
  exclusive                = false
  recursive                = true
  timeout                  = 1
}

