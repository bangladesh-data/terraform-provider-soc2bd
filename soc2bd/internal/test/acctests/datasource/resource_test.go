package datasource

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"testing"

	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/attr"
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/test"
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/test/acctests"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatasourceSoc2bdResource_basic(t *testing.T) {
	t.Run("Test Soc2bd Datasource : Acc Resource Basic", func(t *testing.T) {
		networkName := test.RandomName()
		resourceName := test.RandomResourceName()

		resource.Test(t, resource.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			CheckDestroy:      acctests.CheckSoc2bdResourceDestroy,
			Steps: []resource.TestStep{
				{
					Config: testDatasourceSoc2bdResource(networkName, resourceName),
					Check: acctests.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.soc2bd_resource.out_dr1", attr.Name, resourceName),
					),
				},
			},
		})
	})
}

func testDatasourceSoc2bdResource(networkName, resourceName string) string {
	return fmt.Sprintf(`
	resource "soc2bd_remote_network" "test_dr1" {
	  name = "%s"
	}

	resource "soc2bd_resource" "test_dr1" {
	  name = "%s"
	  address = "acc-test.com"
	  remote_network_id = soc2bd_remote_network.test_dr1.id
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

	data "soc2bd_resource" "out_dr1" {
	  id = soc2bd_resource.test_dr1.id
	}
	`, networkName, resourceName)
}

func TestAccDatasourceSoc2bdResource_negative(t *testing.T) {
	t.Run("Test Soc2bd Datasource : Acc Resource - does not exists", func(t *testing.T) {
		resourceID := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("Resource:%d", acctest.RandInt())))

		resource.Test(t, resource.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck: func() {
				acctests.PreCheck(t)
			},
			Steps: []resource.TestStep{
				{
					Config:      testSoc2bdResourceDoesNotExists(resourceID),
					ExpectError: regexp.MustCompile("Error: failed to read resource with id"),
				},
			},
		})
	})
}

func testSoc2bdResourceDoesNotExists(id string) string {
	return fmt.Sprintf(`
	data "soc2bd_resource" "test_dr2" {
	  id = "%s"
	}

	output "my_resource_dr2" {
	  value = data.soc2bd_resource.test_dr2.name
	}
	`, id)
}

func TestAccDatasourceSoc2bdResource_invalidID(t *testing.T) {
	t.Run("Test Soc2bd Datasource : Acc Resource - failed parse resource ID", func(t *testing.T) {
		networkID := acctest.RandString(10)

		resource.Test(t, resource.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck: func() {
				acctests.PreCheck(t)
			},
			Steps: []resource.TestStep{
				{
					Config:      testSoc2bdResourceDoesNotExists(networkID),
					ExpectError: regexp.MustCompile("Unable to parse global ID"),
				},
			},
		})
	})
}
