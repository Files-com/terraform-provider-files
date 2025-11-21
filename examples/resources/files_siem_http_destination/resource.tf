resource "files_siem_http_destination" "example_siem_http_destination" {
  name                                     = "example"
  additional_headers                       = {
    key = "example value"
  }
  sending_active                           = true
  generic_payload_type                     = "example"
  file_destination_path                    = "example"
  file_format                              = "example"
  file_interval_minutes                    = 1
  azure_dcr_immutable_id                   = "example"
  azure_stream_name                        = "example"
  azure_oauth_client_credentials_tenant_id = "example"
  azure_oauth_client_credentials_client_id = "example"
  qradar_username                          = "example"
  sftp_action_send_enabled                 = true
  ftp_action_send_enabled                  = true
  web_dav_action_send_enabled              = true
  sync_send_enabled                        = true
  outbound_connection_send_enabled         = true
  automation_send_enabled                  = true
  api_request_send_enabled                 = true
  public_hosting_request_send_enabled      = true
  email_send_enabled                       = true
  exavault_api_request_send_enabled        = true
  settings_change_send_enabled             = true
  destination_type                         = "example"
  destination_url                          = "example"
}

