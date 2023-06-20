package datasource

import (
	"fmt"
	"testing"

	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/attr"
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/test"
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/test/acctests"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	resourcesLen     = attr.Len(attr.Resources)
	resourceNamePath = attr.Path(attr.Resources, attr.Name)
)

func TestAccDatasourceSoc2bdResources_basic(t *testing.T) {
	t.Run("Test Soc2bd Datasource : Acc Resources Basic", func(t *testing.T) {
		acctests.SetPageLimit(1)
		networkName := test.RandomName()
		resourceName := test.RandomResourceName()
		const theDatasource = "data.soc2bd_resources.out_drs1"

		resource.Test(t, resource.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			CheckDestroy:      acctests.CheckSoc2bdResourceDestroy,
			Steps: []resource.TestStep{
				{
					Config: testDatasourceSoc2bdResources(networkName, resourceName),
					Check: acctests.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(theDatasource, resourcesLen, "2"),
						resource.TestCheckResourceAttr(theDatasource, resourceNamePath, resourceName),
					),
				},
			},
		})
	})
}

func testDatasourceSoc2bdResources(networkName, resourceName string) string {
	return fmt.Sprintf(`
	resource "soc2bd_remote_network" "test_drs1" {
	  name = "%s"
	}

	resource "soc2bd_resource" "test_drs1_1" {
	  name = "%s"
	  address = "acc-test.com"
	  remote_network_id = soc2bd_remote_network.test_drs1.id
	  protocols {
	    allow_icmp = true
	    tcp {
	      policy = "RESTRICTED"
	      ports = ["80-83", "85"]
	    }
	    udp {
	      policy = "ALLOW_ALL"
	      ports = []
	    }
	  }
	}

	resource "soc2bd_resource" "test_drs1_2" {
	  name = "%s"
	  address = "acc-test.com"
	  remote_network_id = soc2bd_remote_network.test_drs1.id
	  protocols {
	    allow_icmp = true
	    tcp {
	      policy = "ALLOW_ALL"
	      ports = []
	    }
	    udp {
	      policy = "ALLOW_ALL"
	      ports = []
	    }
	  }
	}

	data "soc2bd_resources" "out_drs1" {
	  name = "%s"

	  depends_on = [soc2bd_resource.test_drs1_1, soc2bd_resource.test_drs1_2]
	}
	`, networkName, resourceName, resourceName, resourceName)
}

func TestAccDatasourceSoc2bdResources_emptyResult(t *testing.T) {
	t.Run("Test Soc2bd Datasource : Acc Resources - empty result", func(t *testing.T) {
		resourceName := test.RandomResourceName()

		resource.Test(t, resource.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck: func() {
				acctests.PreCheck(t)
			},
			Steps: []resource.TestStep{
				{
					Config: testSoc2bdResourcesDoesNotExists(resourceName),
					Check: acctests.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.soc2bd_resources.out_drs2", resourcesLen, "0"),
					),
				},
			},
		})
	})
}

func testSoc2bdResourcesDoesNotExists(name string) string {
	return fmt.Sprintf(`
	data "soc2bd_resources" "out_drs2" {
	  name = "%s"
	}

	output "my_resources_drs2" {
	  value = data.soc2bd_resources.out_drs2.resources
	}
	`, name)
}
