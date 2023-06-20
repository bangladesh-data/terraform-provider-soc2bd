package datasource

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"testing"

	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/test"
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/test/acctests"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatasourceSoc2bdConnector_basic(t *testing.T) {
	t.Run("Test Soc2bd Datasource : Acc Connector Basic", func(t *testing.T) {
		networkName := test.RandomName()
		connectorName := test.RandomConnectorName()

		resource.Test(t, resource.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			CheckDestroy:      acctests.CheckSoc2bdConnectorDestroy,
			Steps: []resource.TestStep{
				{
					Config: testDatasourceSoc2bdConnector(networkName, connectorName),
					Check: acctests.ComposeTestCheckFunc(
						resource.TestCheckOutput("my_connector", connectorName),
						resource.TestCheckOutput("my_connector_notification_status", "true"),
					),
				},
			},
		})
	})
}

func testDatasourceSoc2bdConnector(remoteNetworkName, connectorName string) string {
	return fmt.Sprintf(`
	resource "soc2bd_remote_network" "test_dc1" {
	  name = "%s"
	}
	resource "soc2bd_connector" "test_dc1" {
	  remote_network_id = soc2bd_remote_network.test_dc1.id
	  name  = "%s"
	}

	data "soc2bd_connector" "out_dc1" {
	  id = soc2bd_connector.test_dc1.id
	}

	output "my_connector" {
	  value = data.soc2bd_connector.out_dc1.name
	}

	output "my_connector_notification_status" {
	  value = data.soc2bd_connector.out_dc1.status_updates_enabled
	}
	`, remoteNetworkName, connectorName)
}

func TestAccDatasourceSoc2bdConnector_negative(t *testing.T) {
	t.Run("Test Soc2bd Datasource : Acc Connector - does not exists", func(t *testing.T) {
		connectorID := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("Connector:%d", acctest.RandInt())))

		resource.Test(t, resource.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck: func() {
				acctests.PreCheck(t)
			},
			Steps: []resource.TestStep{
				{
					Config:      testSoc2bdConnectorDoesNotExists(connectorID),
					ExpectError: regexp.MustCompile("Error: failed to read connector with id"),
				},
			},
		})
	})
}

func testSoc2bdConnectorDoesNotExists(id string) string {
	return fmt.Sprintf(`
	data "soc2bd_connector" "test_dc2" {
	  id = "%s"
	}

	output "my_connector" {
	  value = data.soc2bd_connector.test_dc2.name
	}
	`, id)
}

func TestAccDatasourceSoc2bdConnector_invalidID(t *testing.T) {
	t.Run("Test Soc2bd Datasource : Acc Connector - failed parse ID", func(t *testing.T) {
		connectorID := acctest.RandString(10)

		resource.Test(t, resource.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck: func() {
				acctests.PreCheck(t)
			},
			Steps: []resource.TestStep{
				{
					Config:      testSoc2bdConnectorDoesNotExists(connectorID),
					ExpectError: regexp.MustCompile("Unable to parse global ID"),
				},
			},
		})
	})
}
