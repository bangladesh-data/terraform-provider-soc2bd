package datasource

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"testing"

	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/test/acctests"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatasourceSoc2bdSecurityPolicyInvalidID(t *testing.T) {
	t.Run("Test Soc2bd Datasource : Acc Security Policy - failed parse ID", func(t *testing.T) {
		randStr := acctest.RandString(10)

		resource.Test(t, resource.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck: func() {
				acctests.PreCheck(t)
			},
			Steps: []resource.TestStep{
				{
					Config:      testDatasourceSoc2bdSecurityPolicy(randStr),
					ExpectError: regexp.MustCompile("Unable to parse global ID"),
				},
			},
		})
	})
}

func testDatasourceSoc2bdSecurityPolicy(id string) string {
	return fmt.Sprintf(`
	data "soc2bd_security_policy" "test_1" {
	  id = "%s"
	}

	output "security_policy_name" {
	  value = data.soc2bd_security_policy.test_1.name
	}
	`, id)
}

func TestAccDatasourceSoc2bdSecurityPolicyReadWithNameAndID(t *testing.T) {
	t.Run("Test Soc2bd Datasource : Acc Security Policy - read with name and id", func(t *testing.T) {
		randStr := acctest.RandString(10)

		resource.Test(t, resource.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck: func() {
				acctests.PreCheck(t)
			},
			Steps: []resource.TestStep{
				{
					Config:      testDatasourceSoc2bdSecurityPolicyWithNameAndID(randStr, randStr),
					ExpectError: regexp.MustCompile("Error: Invalid combination of arguments"),
				},
			},
		})
	})
}

func testDatasourceSoc2bdSecurityPolicyWithNameAndID(id, name string) string {
	return fmt.Sprintf(`
	data "soc2bd_security_policy" "test_2" {
	  id = "%s"
	  name = "%s"
	}

	output "security_policy_name" {
	  value = data.soc2bd_security_policy.test_2.name
	}
	`, id, name)
}

func TestAccDatasourceSoc2bdSecurityPolicyDoesNotExists(t *testing.T) {
	t.Run("Test Soc2bd Datasource : Acc Security Policy - does not exists", func(t *testing.T) {
		securityPolicyID := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("SecurityPolicy:%d", acctest.RandInt())))

		resource.Test(t, resource.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck: func() {
				acctests.PreCheck(t)
			},
			Steps: []resource.TestStep{
				{
					Config:      testDatasourceSoc2bdSecurityPolicy(securityPolicyID),
					ExpectError: regexp.MustCompile("Error: failed to read security policy with id"),
				},
			},
		})
	})
}

func TestAccDatasourceSoc2bdSecurityPolicyReadOkByID(t *testing.T) {
	t.Run("Test Soc2bd Datasource : Acc Security Policy - read Ok By ID", func(t *testing.T) {
		securityPolicies, err := acctests.ListSecurityPolicies()
		if err != nil {
			t.Skip("can't run test:", err)
		}

		testPolicy := securityPolicies[0]

		resource.Test(t, resource.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck: func() {
				acctests.PreCheck(t)
			},
			Steps: []resource.TestStep{
				{
					Config: testDatasourceSoc2bdSecurityPolicyByID(testPolicy.ID),
					Check: acctests.ComposeTestCheckFunc(
						resource.TestCheckOutput("security_policy_name", testPolicy.Name),
					),
				},
			},
		})
	})
}

func testDatasourceSoc2bdSecurityPolicyByID(id string) string {
	return fmt.Sprintf(`
	data "soc2bd_security_policy" "test" {
	  id = "%s"
	}

	output "security_policy_name" {
	  value = data.soc2bd_security_policy.test.name
	}
	`, id)
}

func TestAccDatasourceSoc2bdSecurityPolicyReadOkByName(t *testing.T) {
	t.Run("Test Soc2bd Datasource : Acc Security Policy - read Ok By Name", func(t *testing.T) {
		securityPolicies, err := acctests.ListSecurityPolicies()
		if err != nil {
			t.Skip("can't run test:", err)
		}

		testPolicy := securityPolicies[0]

		resource.Test(t, resource.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck: func() {
				acctests.PreCheck(t)
			},
			Steps: []resource.TestStep{
				{
					Config: testDatasourceSoc2bdSecurityPolicyByName(testPolicy.Name),
					Check: acctests.ComposeTestCheckFunc(
						resource.TestCheckOutput("security_policy_id", testPolicy.ID),
					),
				},
			},
		})
	})
}

func testDatasourceSoc2bdSecurityPolicyByName(name string) string {
	return fmt.Sprintf(`
	data "soc2bd_security_policy" "test" {
	  name = "%s"
	}

	output "security_policy_id" {
	  value = data.soc2bd_security_policy.test.id
	}
	`, name)
}
