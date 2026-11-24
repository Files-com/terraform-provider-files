resource "files_database_destination" "example_database_destination" {
  name         = "Orders"
  provider_    = "snowflake"
  config       = {
    account_identifier = "value"
    database_name      = "value"
    schema_name        = "value"
    table_name         = "value"
  }
  secret_id    = 1
  workspace_id = 1
}

