resource "files_form_field_set" "example_form_field_set" {
  user_id                           = 1
  title                             = "Sample Form Title"
  require_complete_metadata_capture = false
  workspace_id                      = 1
  skip_email                        = false
  skip_name                         = false
  skip_company                      = false
  form_fields                       = [
    {
      label              = "Department"
      field_type         = "dropdown"
      required           = true
      options_for_select = ["Finance", "Sales"]
    }
  ]
}

