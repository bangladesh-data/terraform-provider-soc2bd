package datasource

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/attr"
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/test"
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/test/acctests"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	groupsLen         = attr.Len(attr.Groups)
	groupNamePath     = attr.Path(attr.Groups, attr.Name)
	groupPolicyIDPath = attr.Path(attr.Groups, attr.SecurityPolicyID)
)

func TestAccDatasourceSoc2bdGroups_basic(t *testing.T) {
	t.Run("Test Soc2bd Datasource : Acc Groups Basic", func(t *testing.T) {
		groupName := test.RandomName()

		const theDatasource = "data.soc2bd_groups.out_dgs1"

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
					Config: testDatasourceSoc2bdGroups(groupName, testPolicy.ID),
					Check: acctests.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(theDatasource, groupsLen, "2"),
						resource.TestCheckResourceAttr(theDatasource, groupNamePath, groupName),
						resource.TestCheckResourceAttr(theDatasource, groupPolicyIDPath, testPolicy.ID),
					),
				},
			},
		})
	})
}

func testDatasourceSoc2bdGroups(name, securityPolicyID string) string {
	return fmt.Sprintf(`
	resource "soc2bd_group" "test_dgs1_1" {
	  name = "%s"
	  security_policy_id = "%s"
	}

	resource "soc2bd_group" "test_dgs1_2" {
	  name = "%s"
	  security_policy_id = "%s"
	}

	data "soc2bd_groups" "out_dgs1" {
	  name = "%s"

	  depends_on = [soc2bd_group.test_dgs1_1, soc2bd_group.test_dgs1_2]
	}
	`, name, securityPolicyID, name, securityPolicyID, name)
}

func TestAccDatasourceSoc2bdGroups_emptyResult(t *testing.T) {
	t.Run("Test Soc2bd Datasource : Acc Groups - empty result", func(t *testing.T) {
		groupName := test.RandomName()

		resource.Test(t, resource.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			Steps: []resource.TestStep{
				{
					Config: testSoc2bdGroupsDoesNotExists(groupName),
					Check: acctests.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.soc2bd_groups.out_dgs2", groupsLen, "0"),
					),
				},
			},
		})
	})
}

func testSoc2bdGroupsDoesNotExists(name string) string {
	return fmt.Sprintf(`
	data "soc2bd_groups" "out_dgs2" {
	  name = "%s"
	}
	`, name)
}

func TestAccDatasourceSoc2bdGroupsWithFilters_basic(t *testing.T) {
	acctests.SetPageLimit(1)
	groupName := test.RandomName()

	const theDatasource = "data.soc2bd_groups.out_dgs2"

	t.Run("Test Soc2bd Datasource : Acc Groups with filters - basic", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			Steps: []resource.TestStep{
				{
					Config: testDatasourceSoc2bdGroupsWithFilters(groupName),
					Check: acctests.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(theDatasource, groupsLen, "2"),
						resource.TestCheckResourceAttr(theDatasource, groupNamePath, groupName),
					),
				},
			},
		})
	})
}

func testDatasourceSoc2bdGroupsWithFilters(name string) string {
	return fmt.Sprintf(`
	resource "soc2bd_group" "test_dgs2_1" {
	  name = "%s"
	}

	resource "soc2bd_group" "test_dgs2_2" {
	  name = "%s"
	}

	data "soc2bd_groups" "out_dgs2" {
	  name = "%s"
	  type = "MANUAL"
	  is_active = true

	  depends_on = [soc2bd_group.test_dgs2_1, soc2bd_group.test_dgs2_2]
	}
	`, name, name, name)
}

func TestAccDatasourceSoc2bdGroupsWithFilters_ErrorNotSupportedTypes(t *testing.T) {
	t.Run("Test Soc2bd Datasource : Acc Groups with filters - error not supported types", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck: func() {
				acctests.PreCheck(t)
			},
			Steps: []resource.TestStep{
				{
					Config:      testSoc2bdGroupsWithFilterNotSupportedType(),
					ExpectError: regexp.MustCompile("Error: expected type to be one of"),
				},
			},
		})
	})
}

func testSoc2bdGroupsWithFilterNotSupportedType() string {
	return `
	data "soc2bd_groups" "test" {
	  type = "OTHER"
	}

	output "my_groups" {
	  value = data.soc2bd_groups.test.groups
	}
	`
}

func TestAccDatasourceSoc2bdGroups_WithEmptyFilters(t *testing.T) {
	t.Run("Test Soc2bd Datasource : Acc Groups - with empty filters", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck: func() {
				acctests.PreCheck(t)
			},
			Steps: []resource.TestStep{
				{
					Config: testSoc2bdGroupsWithEmptyFilter(),
				},
			},
		})
	})
}

func testSoc2bdGroupsWithEmptyFilter() string {
	return `
	data "soc2bd_groups" "all" {}

	output "my_groups" {
	  value = data.soc2bd_groups.all.groups
	}
	`
}

func TestAccDatasourceSoc2bdGroups_withTwoDatasource(t *testing.T) {
	t.Run("Test Soc2bd Datasource : Acc Groups with two datasource", func(t *testing.T) {

		groupName := test.RandomName()

		resource.Test(t, resource.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			CheckDestroy:      acctests.CheckSoc2bdGroupDestroy,
			Steps: []resource.TestStep{
				{
					Config: testDatasourceSoc2bdGroupsWithDatasource(groupName),
					Check: acctests.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.soc2bd_groups.two_dgs3", groupNamePath, groupName),
						resource.TestCheckResourceAttr("data.soc2bd_groups.one_dgs3", groupsLen, "1"),
						resource.TestCheckResourceAttr("data.soc2bd_groups.two_dgs3", groupsLen, "2"),
					),
				},
			},
		})
	})
}

func testDatasourceSoc2bdGroupsWithDatasource(name string) string {
	return fmt.Sprintf(`
	resource "soc2bd_group" "test_dgs3_1" {
	  name = "%s"
	}

	resource "soc2bd_group" "test_dgs3_2" {
	  name = "%s"
	}

	resource "soc2bd_group" "test_dgs3_3" {
	  name = "%s-1"
	}

	data "soc2bd_groups" "two_dgs3" {
	  name = "%s"

	  depends_on = [soc2bd_group.test_dgs3_1, soc2bd_group.test_dgs3_2, soc2bd_group.test_dgs3_3]
	}

	data "soc2bd_groups" "one_dgs3" {
	  name = "%s-1"

	  depends_on = [soc2bd_group.test_dgs3_1, soc2bd_group.test_dgs3_2, soc2bd_group.test_dgs3_3]
	}
	`, name, name, name, name, name)
}
