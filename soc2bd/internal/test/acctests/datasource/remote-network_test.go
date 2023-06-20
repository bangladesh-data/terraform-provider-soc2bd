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

func TestAccDatasourceSoc2bdRemoteNetwork_basic(t *testing.T) {
	t.Run("Test Soc2bd Datasource : Acc Remote Network Basic", func(t *testing.T) {

		networkName := test.RandomName()

		resource.Test(t, resource.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			CheckDestroy:      acctests.CheckSoc2bdRemoteNetworkDestroy,
			Steps: []resource.TestStep{
				{
					Config: testDatasourceSoc2bdRemoteNetwork(networkName),
					Check: acctests.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.soc2bd_remote_network.test_dn1_2", attr.Name, networkName),
					),
				},
			},
		})
	})
}

func testDatasourceSoc2bdRemoteNetwork(name string) string {
	return fmt.Sprintf(`
	resource "soc2bd_remote_network" "test_dn1_1" {
	  name = "%s"
	}

	data "soc2bd_remote_network" "test_dn1_2" {
	  id = soc2bd_remote_network.test_dn1_1.id
	}

	output "my_network_dn1_" {
	  value = data.soc2bd_remote_network.test_dn1_2.name
	}
	`, name)
}

func TestAccDatasourceSoc2bdRemoteNetworkByName_basic(t *testing.T) {
	t.Run("Test Soc2bd Datasource : Acc Remote Network Basic", func(t *testing.T) {

		networkName := test.RandomName()

		resource.Test(t, resource.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			CheckDestroy:      acctests.CheckSoc2bdRemoteNetworkDestroy,
			Steps: []resource.TestStep{
				{
					Config: testDatasourceSoc2bdRemoteNetworkByName(networkName),
					Check: acctests.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.soc2bd_remote_network.test_dn2_2", attr.Name, networkName),
					),
				},
			},
		})
	})
}

func testDatasourceSoc2bdRemoteNetworkByName(name string) string {
	return fmt.Sprintf(`
	resource "soc2bd_remote_network" "test_dn2_1" {
	  name = "%s"
	}

	data "soc2bd_remote_network" "test_dn2_2" {
	  name = "%s"
	  depends_on = [resource.soc2bd_remote_network.test_dn2_1]
	}

	output "my_network_dn2" {
	  value = data.soc2bd_remote_network.test_dn2_2.name
	}
	`, name, name)
}

func TestAccDatasourceSoc2bdRemoteNetwork_negative(t *testing.T) {
	t.Run("Test Soc2bd Datasource : Acc Remote Network - does not exists", func(t *testing.T) {
		networkID := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("RemoteNetwork:%d", acctest.RandInt())))

		resource.Test(t, resource.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			Steps: []resource.TestStep{
				{
					Config:      testSoc2bdRemoteNetworkDoesNotExists(networkID),
					ExpectError: regexp.MustCompile("Error: failed to read remote network with id"),
				},
			},
		})
	})
}

func testSoc2bdRemoteNetworkDoesNotExists(id string) string {
	return fmt.Sprintf(`
	data "soc2bd_remote_network" "test_dn3" {
	  id = "%s"
	}

	output "my_network_dn3" {
	  value = data.soc2bd_remote_network.test_dn3.name
	}
	`, id)
}

func TestAccDatasourceSoc2bdRemoteNetwork_invalidNetworkID(t *testing.T) {
	t.Run("Test Soc2bd Datasource : Acc Remote Network - failed parse network ID", func(t *testing.T) {
		networkID := acctest.RandString(10)

		resource.Test(t, resource.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			Steps: []resource.TestStep{
				{
					Config:      testSoc2bdRemoteNetworkDoesNotExists(networkID),
					ExpectError: regexp.MustCompile("Unable to parse global ID"),
				},
			},
		})
	})
}

func TestAccDatasourceSoc2bdRemoteNetwork_bothNetworkIDAndName(t *testing.T) {
	t.Run("Test Soc2bd Datasource : Acc Remote Network - failed passing both network ID and name", func(t *testing.T) {
		networkID := acctest.RandString(10)
		networkName := acctest.RandString(10)

		resource.Test(t, resource.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck: func() {
				acctests.PreCheck(t)
			},
			Steps: []resource.TestStep{
				{
					Config:      testSoc2bdRemoteNetworkValidationFailed(networkID, networkName),
					ExpectError: regexp.MustCompile("Invalid combination of arguments"),
				},
			},
		})
	})
}

func testSoc2bdRemoteNetworkValidationFailed(id, name string) string {
	return fmt.Sprintf(`
	data "soc2bd_remote_network" "test_dn4" {
	  id = "%s"
	  name = "%s"
	}

	output "my_network_dn4" {
	  value = data.soc2bd_remote_network.test_dn4.name
	}
	`, id, name)
}
