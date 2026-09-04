package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestBehaviorDataSource(t *testing.T) {
	VcrTest(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(t.Name()),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `data "files_behavior" "foo_file_expiration" { id = 272834 }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.files_behavior.foo_file_expiration", "id", "272834"),
					resource.TestCheckResourceAttr("data.files_behavior.foo_file_expiration", "behavior", "file_expiration"),
					resource.TestCheckResourceAttr("data.files_behavior.foo_file_expiration", "path", "Foo"),
					resource.TestCheckResourceAttr("data.files_behavior.foo_file_expiration", "value", "14"),
				),
			},
			{
				Config: providerConfig + `
data "files_behavior" "foo_file_expiration" {
  id           = 272834
  value_format = "typed"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.files_behavior.foo_file_expiration", "id", "272834"),
					resource.TestCheckResourceAttr("data.files_behavior.foo_file_expiration", "behavior", "file_expiration"),
					resource.TestCheckResourceAttr("data.files_behavior.foo_file_expiration", "path", "Foo"),
					resource.TestCheckResourceAttr("data.files_behavior.foo_file_expiration", "value.file_expiration.days_to_retain", "14"),
				),
			},
			{
				Config: providerConfig + `
data "files_behavior" "bar_serve_publicly" {
  id           = 272843
  value_format = "typed"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.files_behavior.bar_serve_publicly", "behavior", "serve_publicly"),
					resource.TestCheckResourceAttr("data.files_behavior.bar_serve_publicly", "value.serve_publicly.password_required", "true"),
				),
			},
		},
	})
}

func TestBehaviorResource(t *testing.T) {
	VcrTest(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(t.Name()),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "files_behavior" "bar_serve_publicly" {
  path     = "Bar"
  behavior = "serve_publicly"
  value    = {
    key = "Bar"
    show_index = true
    force_download = false
    password = "secret"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("files_behavior.bar_serve_publicly", "id", "272843"),
					resource.TestCheckResourceAttr("files_behavior.bar_serve_publicly", "behavior", "serve_publicly"),
					resource.TestCheckResourceAttr("files_behavior.bar_serve_publicly", "path", "Bar"),
					resource.TestCheckResourceAttr("files_behavior.bar_serve_publicly", "value.key", "Bar"),
					resource.TestCheckResourceAttr("files_behavior.bar_serve_publicly", "value.show_index", "true"),
					resource.TestCheckResourceAttr("files_behavior.bar_serve_publicly", "value.password", "secret"),
				),
			},
			{
				Config: providerConfig + `
resource "files_behavior" "bar_serve_publicly" {
  path     = "Bar"
  behavior = "serve_publicly"
  value    = {
    serve_publicly = {
      key = "Bar"
      show_index = true
      force_download = false
      password = "secret"
    }
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("files_behavior.bar_serve_publicly", "id", "272843"),
					resource.TestCheckResourceAttr("files_behavior.bar_serve_publicly", "behavior", "serve_publicly"),
					resource.TestCheckResourceAttr("files_behavior.bar_serve_publicly", "path", "Bar"),
					resource.TestCheckResourceAttr("files_behavior.bar_serve_publicly", "value.serve_publicly.key", "Bar"),
					resource.TestCheckResourceAttr("files_behavior.bar_serve_publicly", "value.serve_publicly.show_index", "true"),
					resource.TestCheckResourceAttr("files_behavior.bar_serve_publicly", "value.serve_publicly.password", "secret"),
				),
			},
		},
	})
}

func TestBehaviorFileExpirationResource(t *testing.T) {
	VcrTest(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(t.Name()),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "files_behavior" "primitive_file_expiration" {
  behavior = "file_expiration"
  path     = "Bar"
  value    = 14
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("files_behavior.primitive_file_expiration", "behavior", "file_expiration"),
					resource.TestCheckResourceAttr("files_behavior.primitive_file_expiration", "path", "Bar"),
					resource.TestCheckResourceAttr("files_behavior.primitive_file_expiration", "value", "14"),
				),
			},
		},
	})
}

func TestBehaviorGooglePubSubTransition(t *testing.T) {
	VcrTest(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(t.Name()),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "files_behavior" "google_pub_sub" {
  behavior = "google_pub_sub"
  path     = "Google"
  value = {
    projects_topics = [{ project_id = "my-project", topic_id = "files-events" }]
    google_credentials = {
      type             = "service_account"
      project_id       = "my-project"
      private_key      = "secret"
      universe_domain  = "googleapis.com"
      quota_project_id = "billing-project"
    }
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("files_behavior.google_pub_sub", "value.google_credentials.private_key", "secret"),
					resource.TestCheckResourceAttr("files_behavior.google_pub_sub", "value.google_credentials.quota_project_id", "billing-project"),
				),
			},
			{
				Config: providerConfig + `
resource "files_behavior" "google_pub_sub" {
  behavior = "google_pub_sub"
  path     = "Google"
  value = {
    google_pub_sub = {
      projects_topics = [{ project_id = "my-project", topic_id = "files-events" }]
      google_credentials = {
        type             = "service_account"
        project_id       = "my-project"
        private_key      = "secret"
        universe_domain  = "googleapis.com"
        quota_project_id = "billing-project"
      }
    }
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("files_behavior.google_pub_sub", "value.google_pub_sub.google_credentials.private_key", "secret"),
					resource.TestCheckResourceAttr("files_behavior.google_pub_sub", "value.google_pub_sub.google_credentials.quota_project_id", "billing-project"),
				),
			},
		},
	})
}

func TestBehaviorStorageRegionTransition(t *testing.T) {
	VcrTest(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(t.Name()),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "files_behavior" "storage_region" {
  behavior = "storage_region"
  path     = "Storage"
  value    = "us-east-1"
}
`,
				Check: resource.TestCheckResourceAttr("files_behavior.storage_region", "value", "us-east-1"),
			},
			{
				Config: providerConfig + `
resource "files_behavior" "storage_region" {
  behavior = "storage_region"
  path     = "Storage"
  value    = { storage_region = "us-east-1" }
}
`,
				Check: resource.TestCheckResourceAttr("files_behavior.storage_region", "value.storage_region", "us-east-1"),
			},
		},
	})
}

func TestBehaviorLimitFileRegexTransition(t *testing.T) {
	VcrTest(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(t.Name()),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "files_behavior" "limit_file_regex" {
  behavior = "limit_file_regex"
  path     = "Regex"
  value    = ["/Document-.*/"]
}
`,
				Check: resource.TestCheckResourceAttr("files_behavior.limit_file_regex", "value.0", "/Document-.*/"),
			},
			{
				Config: providerConfig + `
resource "files_behavior" "limit_file_regex" {
  behavior = "limit_file_regex"
  path     = "Regex"
  value    = { limit_file_regex = ["/Document-.*/"] }
}
`,
				Check: resource.TestCheckResourceAttr("files_behavior.limit_file_regex", "value.limit_file_regex.0", "/Document-.*/"),
			},
		},
	})
}

func TestBehaviorMalwareScanningTransition(t *testing.T) {
	VcrTest(t, resource.TestCase{
		ProtoV6ProviderFactories: ProviderFactories(t.Name()),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "files_behavior" "malware_scanning" {
  behavior = "malware_scanning"
  path     = "Malware"
  value    = {}
}
`,
				Check: resource.TestCheckResourceAttr("files_behavior.malware_scanning", "behavior", "malware_scanning"),
			},
			{
				Config: providerConfig + `
resource "files_behavior" "malware_scanning" {
  behavior = "malware_scanning"
  path     = "Malware"
  value    = { malware_scanning = {} }
}
`,
				Check: resource.TestCheckResourceAttr("files_behavior.malware_scanning", "behavior", "malware_scanning"),
			},
		},
	})
}
