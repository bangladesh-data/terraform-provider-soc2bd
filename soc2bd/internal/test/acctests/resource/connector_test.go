package resource

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/attr"
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/provider/resource"
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/test"
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/test/acctests"
	sdk "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var testRegexp = regexp.MustCompile(test.Prefix() + ".*")

func TestAccRemoteConnectorCreate(t *testing.T) {
	t.Run("Test Soc2bd Resource : Acc Remote Connector", func(t *testing.T) {
		const terraformResourceName = "test_c1"
		theResource := acctests.TerraformConnector(terraformResourceName)
		remoteNetworkName := test.RandomName()

		sdk.Test(t, sdk.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			CheckDestroy:      acctests.CheckSoc2bdConnectorDestroy,
			Steps: []sdk.TestStep{
				{
					Config: terraformResourceSoc2bdConnector(terraformResourceName, terraformResourceName, remoteNetworkName),
					Check: acctests.ComposeTestCheckFunc(
						checkSoc2bdConnectorSetWithRemoteNetwork(theResource, acctests.TerraformRemoteNetwork(terraformResourceName)),
						sdk.TestCheckResourceAttrSet(theResource, attr.Name),
					),
				},
			},
		})
	})
}

func TestAccRemoteConnectorWithCustomName(t *testing.T) {
	t.Run("Test Soc2bd Resource : Acc Remote Connector With Custom Name", func(t *testing.T) {
		const terraformResourceName = "test_c2"
		theResource := acctests.TerraformConnector(terraformResourceName)
		remoteNetworkName := test.RandomName()
		connectorName := test.RandomConnectorName()

		sdk.Test(t, sdk.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			CheckDestroy:      acctests.CheckSoc2bdConnectorDestroy,
			Steps: []sdk.TestStep{
				{
					Config: terraformResourceSoc2bdConnectorWithName(terraformResourceName, remoteNetworkName, connectorName),
					Check: acctests.ComposeTestCheckFunc(
						checkSoc2bdConnectorSetWithRemoteNetwork(theResource, acctests.TerraformRemoteNetwork(terraformResourceName)),
						sdk.TestMatchResourceAttr(theResource, attr.Name, regexp.MustCompile(connectorName)),
					),
				},
			},
		})
	})
}

func TestAccRemoteConnectorImport(t *testing.T) {
	t.Run("Test Soc2bd Resource : Acc Remote Connector - Import", func(t *testing.T) {
		const terraformResourceName = "test_c3"
		theResource := acctests.TerraformConnector(terraformResourceName)
		remoteNetworkName := test.RandomName()
		connectorName := test.RandomConnectorName()

		sdk.Test(t, sdk.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			CheckDestroy:      acctests.CheckSoc2bdConnectorDestroy,
			Steps: []sdk.TestStep{
				{
					Config: terraformResourceSoc2bdConnectorWithName(terraformResourceName, remoteNetworkName, connectorName),
					Check: acctests.ComposeTestCheckFunc(
						checkSoc2bdConnectorSetWithRemoteNetwork(theResource, acctests.TerraformRemoteNetwork(terraformResourceName)),
						sdk.TestMatchResourceAttr(theResource, attr.Name, testRegexp),
					),
				},
				{
					ResourceName:      theResource,
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})
}

func TestAccRemoteConnectorNotAllowedToChangeRemoteNetworkId(t *testing.T) {
	t.Run("Test Soc2bd Resource : Acc Remote Connector - should fail on remote_network_id update", func(t *testing.T) {
		const (
			terraformConnectorName      = "test_c4"
			terraformRemoteNetworkName1 = "test_c4_1"
			terraformRemoteNetworkName2 = "test_c4_2"
		)
		theResource := acctests.TerraformConnector(terraformConnectorName)
		remoteNetworkName1 := test.RandomName()
		remoteNetworkName2 := test.RandomName()

		sdk.Test(t, sdk.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			CheckDestroy:      acctests.CheckSoc2bdConnectorDestroy,
			Steps: []sdk.TestStep{
				{
					Config: terraformResourceSoc2bdConnector(terraformRemoteNetworkName1, terraformConnectorName, remoteNetworkName1),
					Check: acctests.ComposeTestCheckFunc(
						checkSoc2bdConnectorSetWithRemoteNetwork(theResource, acctests.TerraformRemoteNetwork(terraformRemoteNetworkName1)),
					),
				},
				{
					Config:      terraformResourceSoc2bdConnector(terraformRemoteNetworkName2, terraformConnectorName, remoteNetworkName2),
					ExpectError: regexp.MustCompile(resource.ErrNotAllowChangeRemoteNetworkID.Error()),
				},
			},
		})
	})
}

func TestAccSoc2bdConnectorReCreateAfterDeletion(t *testing.T) {
	t.Run("Test Soc2bd Resource : Acc Remote Connector ReCreate After Deletion", func(t *testing.T) {
		const terraformResourceName = "test_c5"
		theResource := acctests.TerraformConnector(terraformResourceName)
		remoteNetworkName := test.RandomName()

		sdk.Test(t, sdk.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			CheckDestroy:      acctests.CheckSoc2bdConnectorDestroy,
			Steps: []sdk.TestStep{
				{
					Config: terraformResourceSoc2bdConnector(terraformResourceName, terraformResourceName, remoteNetworkName),
					Check: acctests.ComposeTestCheckFunc(
						checkSoc2bdConnectorSetWithRemoteNetwork(theResource, acctests.TerraformRemoteNetwork(terraformResourceName)),
						acctests.DeleteSoc2bdResource(theResource, resource.Soc2bdConnector),
					),
					ExpectNonEmptyPlan: true,
				},
				{
					Config: terraformResourceSoc2bdConnector(terraformResourceName, terraformResourceName, remoteNetworkName),
					Check: acctests.ComposeTestCheckFunc(
						checkSoc2bdConnectorSetWithRemoteNetwork(theResource, acctests.TerraformRemoteNetwork(terraformResourceName)),
					),
				},
			},
		})
	})
}

func terraformResourceSoc2bdConnector(terraformRemoteNetworkName, terraformConnectorName, remoteNetworkName string) string {
	return fmt.Sprintf(`
	%s

	resource "soc2bd_connector" "%s" {
	  remote_network_id = soc2bd_remote_network.%s.id
	}
	`, terraformResourceRemoteNetwork(terraformRemoteNetworkName, remoteNetworkName), terraformConnectorName, terraformRemoteNetworkName)
}

func terraformResourceSoc2bdConnectorWithName(terraformResourceName, remoteNetworkName, connectorName string) string {
	return fmt.Sprintf(`
	%s

	resource "soc2bd_connector" "%s" {
	  remote_network_id = soc2bd_remote_network.%s.id
      name  = "%s"
	}
	`, terraformResourceRemoteNetwork(terraformResourceName, remoteNetworkName), terraformResourceName, terraformResourceName, connectorName)
}

func checkSoc2bdConnectorSetWithRemoteNetwork(connectorResource, remoteNetworkResource string) sdk.TestCheckFunc {
	return func(s *terraform.State) error {
		connector, ok := s.RootModule().Resources[connectorResource]
		if !ok {
			return fmt.Errorf("Not found: %s ", connectorResource)
		}

		if connector.Primary.ID == "" {
			return fmt.Errorf("No connectorID set ")
		}

		remoteNetwork, ok := s.RootModule().Resources[remoteNetworkResource]
		if !ok {
			return fmt.Errorf("Not found: %s ", remoteNetworkResource)
		}

		if connector.Primary.Attributes[attr.RemoteNetworkID] != remoteNetwork.Primary.ID {
			return fmt.Errorf("Remote Network ID not set properly in the connector ")
		}

		return nil
	}
}

func TestAccRemoteConnectorUpdateName(t *testing.T) {
	t.Run("Test Soc2bd Resource : Acc Remote Connector Update Name", func(t *testing.T) {
		const terraformResourceName = "test_c6"
		theResource := acctests.TerraformConnector(terraformResourceName)
		remoteNetworkName := test.RandomName()
		connectorName := test.RandomConnectorName()

		sdk.Test(t, sdk.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			CheckDestroy:      acctests.CheckSoc2bdConnectorDestroy,
			Steps: []sdk.TestStep{
				{
					Config: terraformResourceSoc2bdConnector(terraformResourceName, terraformResourceName, remoteNetworkName),
					Check: acctests.ComposeTestCheckFunc(
						checkSoc2bdConnectorSetWithRemoteNetwork(theResource, acctests.TerraformRemoteNetwork(terraformResourceName)),
						sdk.TestCheckResourceAttrSet(theResource, attr.Name),
					),
				},
				{
					Config: terraformResourceSoc2bdConnectorWithName(terraformResourceName, remoteNetworkName, connectorName),
					Check: acctests.ComposeTestCheckFunc(
						sdk.TestCheckResourceAttr(theResource, attr.Name, connectorName),
					),
				},
			},
		})
	})
}

func TestAccRemoteConnectorCreateWithNotificationStatus(t *testing.T) {
	t.Run("Test Soc2bd Resource : Acc Remote Connector With Notification Status", func(t *testing.T) {
		const terraformResourceName = "test_c7"
		theResource := acctests.TerraformConnector(terraformResourceName)
		remoteNetworkName := test.RandomName()

		sdk.Test(t, sdk.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			CheckDestroy:      acctests.CheckSoc2bdConnectorDestroy,
			Steps: []sdk.TestStep{
				{
					Config: terraformResourceSoc2bdConnector(terraformResourceName, terraformResourceName, remoteNetworkName),
					Check: acctests.ComposeTestCheckFunc(
						checkSoc2bdConnectorSetWithRemoteNetwork(theResource, acctests.TerraformRemoteNetwork(terraformResourceName)),
						sdk.TestCheckResourceAttrSet(theResource, attr.Name),
					),
				},
				{
					// expecting no changes, as by default notifications enabled
					PlanOnly: true,
					Config:   terraformResourceSoc2bdConnectorWithNotificationStatus(terraformResourceName, terraformResourceName, remoteNetworkName, true),
					Check: acctests.ComposeTestCheckFunc(
						sdk.TestCheckResourceAttr(theResource, attr.StatusUpdatesEnabled, "true"),
					),
				},
				{
					Config: terraformResourceSoc2bdConnectorWithNotificationStatus(terraformResourceName, terraformResourceName, remoteNetworkName, false),
					Check: acctests.ComposeTestCheckFunc(
						sdk.TestCheckResourceAttr(theResource, attr.StatusUpdatesEnabled, "false"),
					),
				},
				{
					// expecting no changes, when user removes `status_updates_enabled` field from terraform
					PlanOnly: true,
					Config:   terraformResourceSoc2bdConnector(terraformResourceName, terraformResourceName, remoteNetworkName),
					Check: acctests.ComposeTestCheckFunc(
						sdk.TestCheckResourceAttr(theResource, attr.StatusUpdatesEnabled, "false"),
					),
				},
			},
		})
	})
}

func terraformResourceSoc2bdConnectorWithNotificationStatus(terraformRemoteNetworkName, terraformConnectorName, remoteNetworkName string, notificationStatus bool) string {
	return fmt.Sprintf(`
	%s

	resource "soc2bd_connector" "%s" {
	  remote_network_id = soc2bd_remote_network.%s.id
	  status_updates_enabled = %v
	}
	`, terraformResourceRemoteNetwork(terraformRemoteNetworkName, remoteNetworkName), terraformConnectorName, terraformRemoteNetworkName, notificationStatus)
}
