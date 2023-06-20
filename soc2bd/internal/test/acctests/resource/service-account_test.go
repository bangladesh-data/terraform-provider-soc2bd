package resource

import (
	"fmt"
	"testing"

	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/attr"
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/provider/resource"
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/test"
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/test/acctests"
	sdk "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func createServiceAccount(resourceName, serviceAccountName string) string {
	return fmt.Sprintf(`
	resource "soc2bd_service_account" "%s" {
	  name = "%s"
	}
	`, resourceName, serviceAccountName)
}

func TestAccSoc2bdServiceAccountCreateUpdate(t *testing.T) {
	t.Run("Test Soc2bd Resource : Acc Service Account Create/Update", func(t *testing.T) {
		const terraformResourceName = "test01"
		theResource := acctests.TerraformServiceAccount(terraformResourceName)
		nameBefore := test.RandomName()
		nameAfter := test.RandomName()

		sdk.Test(t, sdk.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			CheckDestroy:      acctests.CheckSoc2bdServiceAccountDestroy,
			Steps: []sdk.TestStep{
				{
					Config: createServiceAccount(terraformResourceName, nameBefore),
					Check: acctests.ComposeTestCheckFunc(
						acctests.CheckSoc2bdResourceExists(theResource),
						sdk.TestCheckResourceAttr(theResource, attr.Name, nameBefore),
					),
				},
				{
					Config: createServiceAccount(terraformResourceName, nameAfter),
					Check: acctests.ComposeTestCheckFunc(
						acctests.CheckSoc2bdResourceExists(theResource),
						sdk.TestCheckResourceAttr(theResource, attr.Name, nameAfter),
					),
				},
			},
		})
	})
}

func TestAccSoc2bdServiceAccountDeleteNonExisting(t *testing.T) {
	t.Run("Test Soc2bd Resource : Acc Service Account Delete NonExisting", func(t *testing.T) {
		const terraformResourceName = "test02"
		theResource := acctests.TerraformServiceAccount(terraformResourceName)
		name := test.RandomName()

		sdk.Test(t, sdk.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			CheckDestroy:      acctests.CheckSoc2bdServiceAccountDestroy,
			Steps: []sdk.TestStep{
				{
					Config:  createServiceAccount(terraformResourceName, name),
					Destroy: true,
					Check: acctests.ComposeTestCheckFunc(
						acctests.CheckSoc2bdResourceDoesNotExists(theResource),
					),
				},
			},
		})
	})
}

func TestAccSoc2bdServiceAccountReCreateAfterDeletion(t *testing.T) {
	t.Run("Test Soc2bd Resource : Acc Service Account Create After Deletion", func(t *testing.T) {
		const terraformResourceName = "test03"
		theResource := acctests.TerraformServiceAccount(terraformResourceName)
		name := test.RandomName()

		sdk.Test(t, sdk.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			CheckDestroy:      acctests.CheckSoc2bdServiceAccountDestroy,
			Steps: []sdk.TestStep{
				{
					Config: createServiceAccount(terraformResourceName, name),
					Check: acctests.ComposeTestCheckFunc(
						acctests.CheckSoc2bdResourceExists(theResource),
						acctests.DeleteSoc2bdResource(theResource, resource.Soc2bdServiceAccount),
						acctests.WaitTestFunc(),
					),
					ExpectNonEmptyPlan: true,
				},
				{
					Config: createServiceAccount(terraformResourceName, name),
					Check: acctests.ComposeTestCheckFunc(
						acctests.CheckSoc2bdResourceExists(theResource),
					),
				},
			},
		})
	})
}
