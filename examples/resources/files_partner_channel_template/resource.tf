resource "files_partner_channel_template" "example_partner_channel_template" {
  from_partner_folder_name          = "incoming"
  from_partner_managed_folder_paths = ["claims/received"]
  from_partner_route_path_pattern   = "processing/{{partner_name}}/from-partner"
  to_partner_folder_name            = "outgoing"
  to_partner_managed_folder_paths   = ["reports/monthly"]
  to_partner_route_path_pattern     = "delivery/{{partner_name}}/to-partner"
  name                              = "Claims Template"
  path                              = "claims/medical"
  workspace_id                      = 0
}

