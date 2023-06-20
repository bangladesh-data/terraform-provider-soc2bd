package datasource

import (
	"fmt"
	"testing"

	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/test"
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/test/acctests"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatasourceSoc2bdRemoteNetworks_read(t *testing.T) {
	t.Run("Test Soc2bd Datasource : Acc Remote Networks Read", func(t *testing.T) {
		acctests.SetPageLimit(1)

		prefix := acctest.RandString(10)
		networkName1 := test.RandomName(prefix)
		networkName2 := test.RandomName(prefix)

		resource.Test(t, resource.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			CheckDestroy:      acctests.CheckSoc2bdRemoteNetworkDestroy,
			Steps: []resource.TestStep{
				{
					Config: testDatasourceSoc2bdRemoteNetworks2(networkName1, networkName2, prefix),
					Check: acctests.ComposeTestCheckFunc(
						testCheckOutputLength("test_networks", 2),
					),
				},
			},
		})
	})
}

func testDatasourceSoc2bdRemoteNetworks2(networkName1, networkName2, prefix string) string {
	return fmt.Sprintf(`
	resource "soc2bd_remote_network" "test_drn1" {
		name = "%s"
	}
	
	resource "soc2bd_remote_network" "test_drn2" {
		name = "%s"
	}
	
	data "soc2bd_remote_networks" "all" {
		depends_on = [soc2bd_remote_network.test_drn1, soc2bd_remote_network.test_drn2]
	}

	output "test_networks" {
	  	value = [for n in [for net in data.soc2bd_remote_networks.all : net if can(net.*.name)][0] : n if length(regexall("%s.*", n.name)) > 0]
	}
		`, networkName1, networkName2, prefix)
}
