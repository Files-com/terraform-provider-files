resource "files_certificate" "example_certificate" {
  name                      = "Name"
  certificate_domain        = "domain.com"
  certificate_extra_domains = ["www.domain.com"]
  key_type                  = "RSA-4096"
  certificate               = "[certificate]"
  intermediates             = "[certificates]"
}

