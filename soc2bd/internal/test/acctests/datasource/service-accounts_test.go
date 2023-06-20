package datasource

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/attr"
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/test"
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/test/acctests"
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	serviceAccountsLen = attr.Len(attr.ServiceAccounts)
	keyIDsLen          = attr.Len(attr.ServiceAccounts, attr.KeyIDs)
)

func TestAccDatasourceSoc2bdServicesFilterByName(t *testing.T) {
	t.Run("Test Soc2bd Datasource : Acc Services - Filter By Name", func(t *testing.T) {

		name := test.Prefix("orange")
		const (
			terraformResourceName = "dts_service"
			theDatasource         = "data.soc2bd_service_accounts.out"
		)

		config := []terraformServiceConfig{
			{
				serviceName:           name,
				terraformResourceName: test.TerraformRandName(terraformResourceName),
			},
			{
				serviceName:           test.Prefix("lemon"),
				terraformResourceName: test.TerraformRandName(terraformResourceName),
			},
		}

		resource.Test(t, resource.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			CheckDestroy:      acctests.CheckSoc2bdServiceAccountDestroy,
			Steps: []resource.TestStep{
				{
					Config: terraformConfig(
						createServices(config),
						datasourceServices(name, config),
					),
					Check: acctests.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(theDatasource, serviceAccountsLen, "1"),
						resource.TestCheckResourceAttr(theDatasource, keyIDsLen, "1"),
						resource.TestCheckResourceAttr(theDatasource, attr.ID, "service-by-name-"+name),
					),
				},
			},
		})
	})
}

func TestAccDatasourceSoc2bdServicesAll(t *testing.T) {
	t.Run("Test Soc2bd Datasource : Acc Services - All", func(t *testing.T) {
		prefix := test.Prefix() + acctest.RandString(4)
		const (
			terraformResourceName = "dts_service"
			theDatasource         = "data.soc2bd_service_accounts.out"
		)

		config := []terraformServiceConfig{
			{
				serviceName:           prefix + "_orange",
				terraformResourceName: test.TerraformRandName(terraformResourceName),
			},
			{
				serviceName:           prefix + "_lemon",
				terraformResourceName: test.TerraformRandName(terraformResourceName),
			},
		}

		resource.Test(t, resource.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			CheckDestroy:      acctests.CheckSoc2bdServiceAccountDestroy,
			Steps: []resource.TestStep{
				{
					Config: filterDatasourceServices(prefix, config),
					Check: acctests.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(theDatasource, attr.ID, "all-services"),
					),
				},
				{
					Config: filterDatasourceServices(prefix, config),
					Check: acctests.ComposeTestCheckFunc(
						testCheckOutputLength("my_services", 2),
					),
				},
			},
		})
	})
}

func TestAccDatasourceSoc2bdServicesEmptyResult(t *testing.T) {
	t.Run("Test Soc2bd Datasource : Acc Services - Empty Result", func(t *testing.T) {

		const theDatasource = "data.soc2bd_service_accounts.out"

		resource.Test(t, resource.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			CheckDestroy:      acctests.CheckSoc2bdServiceAccountDestroy,
			Steps: []resource.TestStep{
				{
					Config: datasourceServices(test.RandomName(), nil),
					Check: acctests.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(theDatasource, serviceAccountsLen, "0"),
					),
				},
			},
		})
	})
}

type terraformServiceConfig struct {
	terraformResourceName, serviceName string
}

func terraformConfig(resources ...string) string {
	return strings.Join(resources, "\n")
}

func datasourceServices(name string, configs []terraformServiceConfig) string {
	var dependsOn string
	ids := getTerraformServiceKeys(configs)

	if ids != "" {
		dependsOn = fmt.Sprintf("depends_on = [%s]", ids)
	}

	return fmt.Sprintf(`
	data "soc2bd_service_accounts" "out" {
	  name = "%s"

	  %s
	}
	`, name, dependsOn)
}

func createServices(configs []terraformServiceConfig) string {
	return strings.Join(
		utils.Map[terraformServiceConfig, string](configs, func(cfg terraformServiceConfig) string {
			return createServiceKey(cfg.terraformResourceName, cfg.serviceName)
		}),
		"\n",
	)
}

func getTerraformServiceKeys(configs []terraformServiceConfig) string {
	return strings.Join(
		utils.Map[terraformServiceConfig, string](configs, func(cfg terraformServiceConfig) string {
			return acctests.TerraformServiceKey(cfg.terraformResourceName)
		}),
		", ",
	)
}

func createServiceKey(terraformResourceName, serviceName string) string {
	return fmt.Sprintf(`
	%s

	resource "soc2bd_service_account_key" "%s" {
	  service_account_id = soc2bd_service_account.%s.id
	}
	`, createServiceAccount(terraformResourceName, serviceName), terraformResourceName, terraformResourceName)
}

func createServiceAccount(terraformResourceName, serviceName string) string {
	return fmt.Sprintf(`
	resource "soc2bd_service_account" "%s" {
	  name = "%s"
	}
	`, terraformResourceName, serviceName)
}

func filterDatasourceServices(prefix string, configs []terraformServiceConfig) string {
	return fmt.Sprintf(`
	%s

	data "soc2bd_service_accounts" "out" {

	}

	output "my_services" {
	  	value = [for c in data.soc2bd_service_accounts.out.service_accounts : c if length(regexall("^%s", c.name)) > 0]
	}
	`, createServices(configs), prefix)
}

func TestAccDatasourceSoc2bdServicesAllCursors(t *testing.T) {
	t.Run("Test Soc2bd Datasource : Acc Services - All Cursors", func(t *testing.T) {
		acctests.SetPageLimit(1)
		prefix := test.Prefix() + acctest.RandString(4)
		const (
			theDatasource = "data.soc2bd_service_accounts.out"
		)

		resource.Test(t, resource.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			CheckDestroy:      acctests.CheckSoc2bdServiceAccountDestroy,
			Steps: []resource.TestStep{
				{
					Config: datasourceServicesConfig(prefix),
					Check: acctests.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(theDatasource, attr.ID, "all-services"),
					),
				},
				{
					Config: datasourceServicesConfig(prefix),
					Check: acctests.ComposeTestCheckFunc(
						testCheckOutputLength("my_services", 3),
						testCheckOutputNestedLen("my_services", 0, attr.ResourceIDs, 1),
						testCheckOutputNestedLen("my_services", 0, attr.KeyIDs, 2),
					),
				},
			},
		})
	})
}

func datasourceServicesConfig(prefix string) string {
	return fmt.Sprintf(`
    resource "soc2bd_service_account" "%s_1" {
      name = "%s-1"
    }
    
    resource "soc2bd_service_account" "%s_2" {
      name = "%s-2"
    }

    resource "soc2bd_service_account" "%s_3" {
      name = "%s-3"
    }
    
    resource "soc2bd_remote_network" "%s_1" {
      name = "%s-1"
    }
    
    resource "soc2bd_remote_network" "%s_2" {
      name = "%s-2"
    }
    
    resource "soc2bd_resource" "%s_1" {
      name = "%s-1"
      address = "acc-test.com"
      remote_network_id = soc2bd_remote_network.%s_1.id
    
      access {
        service_account_ids = [soc2bd_service_account.%s_1.id, soc2bd_service_account.%s_2.id]
      }
    }
    
    resource "soc2bd_resource" "%s_2" {
      name = "%s-2"
      address = "acc-test.com"
      remote_network_id = soc2bd_remote_network.%s_2.id
    
      access {
        service_account_ids = [soc2bd_service_account.%s_3.id]
      }
    }
    
    resource "soc2bd_service_account_key" "%s_1_1" {
      service_account_id = soc2bd_service_account.%s_1.id
    }
    
    resource "soc2bd_service_account_key" "%s_1_2" {
      service_account_id = soc2bd_service_account.%s_1.id
    }
    
    resource "soc2bd_service_account_key" "%s_2_1" {
      service_account_id = soc2bd_service_account.%s_2.id
    }
    
    resource "soc2bd_service_account_key" "%s_2_2" {
      service_account_id = soc2bd_service_account.%s_2.id
    }
    
    resource "soc2bd_service_account_key" "%s_3_1" {
      service_account_id = soc2bd_service_account.%s_3.id
    }

    resource "soc2bd_service_account_key" "%s_3_2" {
      service_account_id = soc2bd_service_account.%s_3.id
    }
    
    data "soc2bd_service_accounts" "out" {
    
    }
    
    output "my_services" {
      value = [for c in data.soc2bd_service_accounts.out.service_accounts : c if length(regexall("^%s", c.name)) > 0]
    }
`, duplicate(prefix, 32)...)
}

func duplicate(val string, n int) []any {
	result := make([]any, 0, n)
	for i := 0; i < n; i++ {
		result = append(result, val)
	}

	return result
}
