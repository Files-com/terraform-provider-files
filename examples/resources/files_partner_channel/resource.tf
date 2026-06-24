resource "files_partner_channel" "example_partner_channel" {
  from_partner_folder_name = "incoming"
  from_partner_route_path  = "processing/from-partner"
  to_partner_folder_name   = "outgoing"
  to_partner_route_path    = "delivery/to-partner"
  partner_id               = 1
  path                     = "claims/medical"
  workspace_id             = 0
}

