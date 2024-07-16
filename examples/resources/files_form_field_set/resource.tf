resource "files_form_field_set" "example_form_field_set" {
  user_id      = 1
  title        = "Sample Form Title"
  skip_email   = true
  skip_name    = true
  skip_company = true
  form_fields  = [
    {
      label              = "Sample Label"
      required           = true
      help_text          = "Help Text"
      field_type         = "text"
      options_for_select = ["red", "green", "blue"]
      default_option     = "red"
      form_field_set_id  = 1
    }
  ]
}

