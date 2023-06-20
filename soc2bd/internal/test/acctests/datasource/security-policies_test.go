package datasource

import (
	"fmt"
	"testing"

	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/attr"
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/test/acctests"
	sdk "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatasourceSoc2bdSecurityPoliciesBasic(t *testing.T) {
	t.Run("Test Soc2bd Datasource : Acc Security Policies - basic", func(t *testing.T) {
		acctests.SetPageLimit(1)

		securityPolicies, err := acctests.ListSecurityPolicies()
		if err != nil {
			t.Skip("can't run test:", err)
		}

		sdk.Test(t, sdk.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			Steps: []sdk.TestStep{
				{
					Config: testDatasourceSoc2bdSecurityPolicies(),
					Check: acctests.ComposeTestCheckFunc(
						sdk.TestCheckResourceAttr("data.soc2bd_security_policies.all", attr.Len(attr.SecurityPolicies), fmt.Sprintf("%d", len(securityPolicies))),
					),
				},
			},
		})
	})
}

func testDatasourceSoc2bdSecurityPolicies() string {
	return `
	data "soc2bd_security_policies" "all" {}
	`
}
