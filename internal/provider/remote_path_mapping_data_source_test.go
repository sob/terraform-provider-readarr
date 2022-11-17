package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRemotePathMappingDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccRemotePathMappingDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.readarr_remote_path_mapping.test", "id"),
					resource.TestCheckResourceAttr("data.readarr_remote_path_mapping.test", "host", "transmission")),
			},
		},
	})
}

const testAccRemotePathMappingDataSourceConfig = `
resource "readarr_download_client" "test" {
	enable = false
	priority = 1
	name = "remotepatdstest"
	implementation = "Transmission"
	protocol = "torrent"
	config_contract = "TransmissionSettings"
	host = "transmission"
	url_base = "/transmission/"
	port = 9091
}

resource "readarr_remote_path_mapping" "test" {
	host = "transmission"
	remote_path = "/datatest/"
	local_path = "/config/"
}

data "readarr_remote_path_mapping" "test" {
	id = readarr_remote_path_mapping.test.id
}
`
