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

func TestAccDatasourceSoc2bdGroup_basic(t *testing.T) {
	t.Run("Test Soc2bd Datasource : Acc Group Basic", func(t *testing.T) {
		groupName := test.RandomName()

		securityPolicies, err := acctests.ListSecurityPolicies()
		if err != nil {
			t.Skip("can't run test:", err)
		}

		testPolicy := securityPolicies[0]

		resource.Test(t, resource.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			CheckDestroy:      acctests.CheckSoc2bdGroupDestroy,
			Steps: []resource.TestStep{
				{
					Config: testDatasourceSoc2bdGroup(groupName, testPolicy.ID),
					Check: acctests.ComposeTestCheckFunc(
						resource.TestCheckOutput("my_group_dg1", groupName),
						resource.TestCheckOutput("my_group_is_active_dg1", "true"),
						resource.TestCheckOutput("my_group_type_dg1", "MANUAL"),
						resource.TestCheckOutput("my_group_policy_dg1", testPolicy.ID),
					),
				},
			},
		})
	})
}

func testDatasourceSoc2bdGroup(name, securityPolicyID string) string {
	return fmt.Sprintf(`
	resource "soc2bd_group" "foo_dg1" {
	  name = "%s"
	  security_policy_id = "%s"
	}

	data "soc2bd_group" "bar_dg1" {
	  id = soc2bd_group.foo_dg1.id
	}

	output "my_group_dg1" {
	  value = data.soc2bd_group.bar_dg1.name
	}

	output "my_group_is_active_dg1" {
	  value = data.soc2bd_group.bar_dg1.is_active
	}

	output "my_group_type_dg1" {
	  value = data.soc2bd_group.bar_dg1.type
	}

	output "my_group_policy_dg1" {
	  value = data.soc2bd_group.bar_dg1.security_policy_id
	}
	`, name, securityPolicyID)
}

func TestAccDatasourceSoc2bdGroup_negative(t *testing.T) {
	t.Run("Test Soc2bd Datasource : Acc Group - does not exists", func(t *testing.T) {
		groupID := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("Group:%d", acctest.RandInt())))

		resource.Test(t, resource.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck: func() {
				acctests.PreCheck(t)
			},
			Steps: []resource.TestStep{
				{
					Config:      testSoc2bdGroupDoesNotExists(groupID),
					ExpectError: regexp.MustCompile("Error: failed to read group with id"),
				},
			},
		})
	})
}

func testSoc2bdGroupDoesNotExists(id string) string {
	return fmt.Sprintf(`
	data "soc2bd_group" "foo_dg2" {
	  id = "%s"
	}
	`, id)
}

func TestAccDatasourceSoc2bdGroup_invalidGroupID(t *testing.T) {
	t.Run("Test Soc2bd Datasource : Acc Group - failed parse group ID", func(t *testing.T) {
		groupID := acctest.RandString(10)

		resource.Test(t, resource.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck: func() {
				acctests.PreCheck(t)
			},
			Steps: []resource.TestStep{
				{
					Config:      testSoc2bdGroupDoesNotExists(groupID),
					ExpectError: regexp.MustCompile("Unable to parse global ID"),
				},
			},
		})
	})
}
