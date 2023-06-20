package resource

import (
	"context"
	"fmt"
	"testing"

	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/attr"
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/client"
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/provider/resource"
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/test"
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/test/acctests"
	sdk "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccRemoteConnectorWithTokens(t *testing.T) {
	t.Run("Test Soc2bd Resource : Acc Remote Connector With Tokens", func(t *testing.T) {
		const terraformResourceName = "test_t1"
		theResource := acctests.TerraformConnectorTokens(terraformResourceName)
		remoteNetworkName := test.RandomName()

		sdk.Test(t, sdk.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			CheckDestroy:      checkSoc2bdConnectorTokensInvalidated,
			Steps: []sdk.TestStep{
				{
					Config: terraformResourceSoc2bdConnectorTokens(terraformResourceName, remoteNetworkName),
					Check: acctests.ComposeTestCheckFunc(
						checkSoc2bdConnectorTokensSet(theResource),
					),
				},
			},
		})
	})
}

func terraformResourceSoc2bdConnectorTokens(terraformResourceName, remoteNetworkName string) string {
	return fmt.Sprintf(`
	%s

	resource "soc2bd_connector_tokens" "%s" {
	  connector_id = soc2bd_connector.%s.id
      keepers = {
         foo = "bar"
      }
	}
	`, terraformResourceSoc2bdConnector(terraformResourceName, terraformResourceName, remoteNetworkName), terraformResourceName, terraformResourceName)
}

func checkSoc2bdConnectorTokensInvalidated(s *terraform.State) error {
	c := acctests.Provider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resource.Soc2bdConnectorTokens {
			continue
		}

		connectorId := rs.Primary.ID
		accessToken := rs.Primary.Attributes[attr.AccessToken]
		refreshToken := rs.Primary.Attributes[attr.RefreshToken]

		err := c.VerifyConnectorTokens(context.Background(), refreshToken, accessToken)
		// expecting error here , Since tokens invalidated
		if err == nil {
			return fmt.Errorf("connector with ID %s tokens that should be inactive are still active", connectorId)
		}
	}

	return nil
}

func checkSoc2bdConnectorTokensSet(connectorNameTokens string) sdk.TestCheckFunc {
	return func(s *terraform.State) error {
		connectorTokens, ok := s.RootModule().Resources[connectorNameTokens]

		if !ok {
			return fmt.Errorf("not found: %s", connectorNameTokens)
		}

		if connectorTokens.Primary.ID == "" {
			return fmt.Errorf("no connectorTokensID set")
		}

		if connectorTokens.Primary.Attributes[attr.AccessToken] == "" {
			return fmt.Errorf("no access token set")
		}

		if connectorTokens.Primary.Attributes[attr.RefreshToken] == "" {
			return fmt.Errorf("no refresh token set")
		}

		return nil
	}
}
