package resource

import (
	"errors"
	"fmt"
	"testing"

	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/attr"
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/model"
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/provider/resource"
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/test"
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/test/acctests"
	sdk "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var ErrEmptyValue = errors.New("empty value")

func createServiceKey(terraformResourceName, serviceAccountName string) string {
	return fmt.Sprintf(`
	%s

	resource "soc2bd_service_account_key" "%s" {
	  service_account_id = soc2bd_service_account.%s.id
	}
	`, createServiceAccount(terraformResourceName, serviceAccountName), terraformResourceName, terraformResourceName)
}

func createServiceKeyWithName(terraformResourceName, serviceAccountName, serviceKeyName string) string {
	return fmt.Sprintf(`
	%s

	resource "soc2bd_service_account_key" "%s" {
	  service_account_id = soc2bd_service_account.%s.id
	  name = "%s"
	}
	`, createServiceAccount(terraformResourceName, serviceAccountName), terraformResourceName, terraformResourceName, serviceKeyName)
}

func nonEmptyValue(value string) error {
	if value != "" {
		return nil
	}

	return ErrEmptyValue
}

func TestAccSoc2bdServiceKeyCreateUpdate(t *testing.T) {
	t.Run("Test Soc2bd Resource : Acc Service Key Create/Update", func(t *testing.T) {
		serviceAccountName := test.RandomName()
		terraformResourceName := test.TerraformRandName("test_key")
		serviceAccount := acctests.TerraformServiceAccount(terraformResourceName)
		serviceKey := acctests.TerraformServiceKey(terraformResourceName)

		sdk.Test(t, sdk.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			CheckDestroy:      acctests.CheckSoc2bdServiceAccountDestroy,
			Steps: []sdk.TestStep{
				{
					Config: createServiceKey(terraformResourceName, serviceAccountName),
					Check: acctests.ComposeTestCheckFunc(
						acctests.CheckSoc2bdResourceExists(serviceAccount),
						sdk.TestCheckResourceAttr(serviceAccount, attr.Name, serviceAccountName),
						acctests.CheckSoc2bdResourceExists(serviceKey),
						sdk.TestCheckResourceAttrWith(serviceKey, attr.Token, nonEmptyValue),
					),
				},
				{
					Config: createServiceKey(terraformResourceName, serviceAccountName),
					Check: acctests.ComposeTestCheckFunc(
						acctests.CheckSoc2bdResourceExists(serviceAccount),
						sdk.TestCheckResourceAttr(serviceAccount, attr.Name, serviceAccountName),
						acctests.CheckSoc2bdResourceExists(serviceKey),
						sdk.TestCheckResourceAttrWith(serviceKey, attr.Token, nonEmptyValue),
					),
				},
			},
		})
	})
}

func TestAccSoc2bdServiceKeyCreateUpdateWithName(t *testing.T) {
	t.Run("Test Soc2bd Resource : Acc Service Key Create/Update With Name", func(t *testing.T) {
		serviceAccountName := test.RandomName()
		terraformResourceName := test.TerraformRandName("test_key")
		serviceAccount := acctests.TerraformServiceAccount(terraformResourceName)
		serviceKey := acctests.TerraformServiceKey(terraformResourceName)
		beforeName := test.RandomName()
		afterName := test.RandomName()

		sdk.Test(t, sdk.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			CheckDestroy:      acctests.CheckSoc2bdServiceAccountDestroy,
			Steps: []sdk.TestStep{
				{
					Config: createServiceKeyWithName(terraformResourceName, serviceAccountName, beforeName),
					Check: acctests.ComposeTestCheckFunc(
						acctests.CheckSoc2bdResourceExists(serviceAccount),
						sdk.TestCheckResourceAttr(serviceAccount, attr.Name, serviceAccountName),
						acctests.CheckSoc2bdResourceExists(serviceKey),
						sdk.TestCheckResourceAttr(serviceKey, attr.Name, beforeName),
						sdk.TestCheckResourceAttrWith(serviceKey, attr.Token, nonEmptyValue),
					),
				},
				{
					Config: createServiceKeyWithName(terraformResourceName, serviceAccountName, afterName),
					Check: acctests.ComposeTestCheckFunc(
						acctests.CheckSoc2bdResourceExists(serviceAccount),
						sdk.TestCheckResourceAttr(serviceAccount, attr.Name, serviceAccountName),
						acctests.CheckSoc2bdResourceExists(serviceKey),
						sdk.TestCheckResourceAttr(serviceKey, attr.Name, afterName),
						sdk.TestCheckResourceAttrWith(serviceKey, attr.Token, nonEmptyValue),
						acctests.WaitTestFunc(),
					),
				},
			},
		})
	})
}

func TestAccSoc2bdServiceKeyReCreateAfterInactive(t *testing.T) {
	t.Run("Test Soc2bd Resource : Acc Service Key ReCreate After Inactive", func(t *testing.T) {
		serviceAccountName := test.RandomName()
		terraformResourceName := test.TerraformRandName("test_key")
		serviceKey := acctests.TerraformServiceKey(terraformResourceName)

		sdk.Test(t, sdk.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			CheckDestroy:      acctests.CheckSoc2bdServiceAccountDestroy,
			Steps: []sdk.TestStep{
				{
					Config: createServiceKey(terraformResourceName, serviceAccountName),
					Check: acctests.ComposeTestCheckFunc(
						acctests.CheckSoc2bdResourceExists(serviceKey),
						sdk.TestCheckResourceAttrWith(serviceKey, attr.Token, nonEmptyValue),
						acctests.RevokeSoc2bdServiceKey(serviceKey),
						acctests.WaitTestFunc(),
						acctests.CheckSoc2bdServiceKeyStatus(serviceKey, model.StatusRevoked),
					),
				},
				{
					Config: createServiceKey(terraformResourceName, serviceAccountName),
					Check: acctests.ComposeTestCheckFunc(
						acctests.CheckSoc2bdResourceExists(serviceKey),
						acctests.CheckSoc2bdServiceKeyStatus(serviceKey, model.StatusActive),
						sdk.TestCheckResourceAttrWith(serviceKey, attr.Token, nonEmptyValue),
					),
				},
			},
		})
	})
}

func TestAccSoc2bdServiceKeyDelete(t *testing.T) {
	t.Run("Test Soc2bd Resource : Acc Service Key Delete", func(t *testing.T) {
		serviceAccountName := test.RandomName()
		terraformResourceName := test.TerraformRandName("test_key")
		serviceKey := acctests.TerraformServiceKey(terraformResourceName)

		sdk.Test(t, sdk.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			CheckDestroy:      acctests.CheckSoc2bdServiceAccountDestroy,
			Steps: []sdk.TestStep{
				{
					Config:  createServiceKey(terraformResourceName, serviceAccountName),
					Destroy: true,
					Check: acctests.ComposeTestCheckFunc(
						acctests.CheckSoc2bdResourceDoesNotExists(serviceKey),
					),
				},
			},
		})
	})
}

func TestAccSoc2bdServiceKeyReCreateAfterDeletion(t *testing.T) {
	t.Run("Test Soc2bd Resource : Acc Service Key ReCreate After Delete", func(t *testing.T) {
		serviceAccountName := test.RandomName()
		terraformResourceName := test.TerraformRandName("test_key")
		serviceKey := acctests.TerraformServiceKey(terraformResourceName)

		sdk.Test(t, sdk.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			CheckDestroy:      acctests.CheckSoc2bdServiceAccountDestroy,
			Steps: []sdk.TestStep{
				{
					Config: createServiceKey(terraformResourceName, serviceAccountName),
					Check: acctests.ComposeTestCheckFunc(
						acctests.CheckSoc2bdResourceExists(serviceKey),
						acctests.RevokeSoc2bdServiceKey(serviceKey),
						acctests.DeleteSoc2bdResource(serviceKey, resource.Soc2bdServiceAccountKey),
					),
					ExpectNonEmptyPlan: true,
				},
				{
					Config: createServiceKey(terraformResourceName, serviceAccountName),
					Check: acctests.ComposeTestCheckFunc(
						acctests.CheckSoc2bdResourceExists(serviceKey),
						sdk.TestCheckResourceAttrWith(serviceKey, attr.Token, nonEmptyValue),
					),
				},
			},
		})
	})
}
