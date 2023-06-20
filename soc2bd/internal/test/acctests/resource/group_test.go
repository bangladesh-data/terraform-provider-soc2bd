package resource

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/attr"
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/provider/resource"
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/test"
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/test/acctests"
	sdk "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	userIdsLen = attr.Len(attr.UserIDs)
)

func TestAccSoc2bdGroupCreateUpdate(t *testing.T) {
	t.Run("Test Soc2bd Resource : Acc Group Create/Update", func(t *testing.T) {
		const terraformResourceName = "test001"
		theResource := acctests.TerraformGroup(terraformResourceName)
		nameBefore := test.RandomName()
		nameAfter := test.RandomName()

		sdk.Test(t, sdk.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			CheckDestroy:      acctests.CheckSoc2bdGroupDestroy,
			Steps: []sdk.TestStep{
				{
					Config: terraformResourceSoc2bdGroup(terraformResourceName, nameBefore),
					Check: acctests.ComposeTestCheckFunc(
						acctests.CheckSoc2bdResourceExists(theResource),
						sdk.TestCheckResourceAttr(theResource, attr.Name, nameBefore),
					),
				},
				{
					Config: terraformResourceSoc2bdGroup(terraformResourceName, nameAfter),
					Check: acctests.ComposeTestCheckFunc(
						acctests.CheckSoc2bdResourceExists(theResource),
						sdk.TestCheckResourceAttr(theResource, attr.Name, nameAfter),
					),
				},
			},
		})
	})
}

func terraformResourceSoc2bdGroup(terraformResourceName, name string) string {
	return fmt.Sprintf(`
	resource "soc2bd_group" "%s" {
	  name = "%s"
	}
	`, terraformResourceName, name)
}

func TestAccSoc2bdGroupDeleteNonExisting(t *testing.T) {
	t.Run("Test Soc2bd Resource : Acc Group Delete NonExisting", func(t *testing.T) {
		const terraformResourceName = "test002"
		theResource := acctests.TerraformGroup(terraformResourceName)
		groupName := test.RandomName()

		sdk.Test(t, sdk.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			CheckDestroy:      acctests.CheckSoc2bdGroupDestroy,
			Steps: []sdk.TestStep{
				{
					Config:  terraformResourceSoc2bdGroup(terraformResourceName, groupName),
					Destroy: true,
					Check: acctests.ComposeTestCheckFunc(
						acctests.CheckSoc2bdResourceDoesNotExists(theResource),
					),
				},
			},
		})
	})
}

func TestAccSoc2bdGroupReCreateAfterDeletion(t *testing.T) {
	t.Run("Test Soc2bd Resource : Acc Group Create After Deletion", func(t *testing.T) {
		const terraformResourceName = "test003"
		theResource := acctests.TerraformGroup(terraformResourceName)
		groupName := test.RandomName()

		sdk.Test(t, sdk.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			CheckDestroy:      acctests.CheckSoc2bdGroupDestroy,
			Steps: []sdk.TestStep{
				{
					Config: terraformResourceSoc2bdGroup(terraformResourceName, groupName),
					Check: acctests.ComposeTestCheckFunc(
						acctests.CheckSoc2bdResourceExists(theResource),
						acctests.DeleteSoc2bdResource(theResource, resource.Soc2bdGroup),
					),
					ExpectNonEmptyPlan: true,
				},
				{
					Config: terraformResourceSoc2bdGroup(terraformResourceName, groupName),
					Check: acctests.ComposeTestCheckFunc(
						acctests.CheckSoc2bdResourceExists(theResource),
					),
				},
			},
		})
	})
}

func TestAccSoc2bdGroupWithSecurityPolicy(t *testing.T) {
	t.Run("Test Soc2bd Resource : Acc Group Create/Update - With Security Policy", func(t *testing.T) {
		const terraformResourceName = "test004"
		theResource := acctests.TerraformGroup(terraformResourceName)
		name := test.RandomName()

		securityPolicies, err := acctests.ListSecurityPolicies()
		if err != nil {
			t.Skip("can't run test:", err)
		}

		testPolicy := securityPolicies[0]

		sdk.Test(t, sdk.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			CheckDestroy:      acctests.CheckSoc2bdGroupDestroy,
			Steps: []sdk.TestStep{
				{
					Config: terraformResourceSoc2bdGroup(terraformResourceName, name),
					Check: acctests.ComposeTestCheckFunc(
						acctests.CheckSoc2bdResourceExists(theResource),
						sdk.TestCheckResourceAttr(theResource, attr.Name, name),
					),
				},
				{
					Config: terraformResourceSoc2bdGroupWithSecurityPolicy(terraformResourceName, name, testPolicy.ID),
					Check: acctests.ComposeTestCheckFunc(
						acctests.CheckSoc2bdResourceExists(theResource),
						sdk.TestCheckResourceAttr(theResource, attr.Name, name),
						sdk.TestCheckResourceAttr(theResource, attr.SecurityPolicyID, testPolicy.ID),
					),
				},
				{
					// expecting no changes
					PlanOnly: true,
					Config:   terraformResourceSoc2bdGroup(terraformResourceName, name),
					Check: acctests.ComposeTestCheckFunc(
						acctests.CheckSoc2bdResourceExists(theResource),
						sdk.TestCheckResourceAttr(theResource, attr.Name, name),
					),
				},
			},
		})
	})
}

func terraformResourceSoc2bdGroupWithSecurityPolicy(terraformResourceName, name, securityPolicyID string) string {
	return fmt.Sprintf(`
	resource "soc2bd_group" "%s" {
	  name = "%s"
	  security_policy_id = "%s"
	}
	`, terraformResourceName, name, securityPolicyID)
}

func TestAccSoc2bdGroupUsersAuthoritativeByDefault(t *testing.T) {
	t.Run("Test Soc2bd Resource : Acc Group Users Authoritative By Default", func(t *testing.T) {
		const terraformResourceName = "test005"
		theResource := acctests.TerraformGroup(terraformResourceName)
		groupName := test.RandomName()

		users, userIDs := genNewUsers("u005", 3)

		sdk.Test(t, sdk.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			CheckDestroy:      acctests.CheckSoc2bdGroupDestroy,
			Steps: []sdk.TestStep{
				{
					Config: terraformResourceSoc2bdGroupWithUsers(terraformResourceName, groupName, users, userIDs[:1]),
					Check: acctests.ComposeTestCheckFunc(
						sdk.TestCheckResourceAttr(theResource, userIdsLen, "1"),
						acctests.CheckGroupUsersLen(theResource, 1),
					),
				},
				{
					Config: terraformResourceSoc2bdGroupWithUsers(terraformResourceName, groupName, users, userIDs[:1]),
					Check: acctests.ComposeTestCheckFunc(
						// added new user to the group though API
						acctests.AddGroupUser(theResource, groupName, userIDs[1]),
						acctests.WaitTestFunc(),
						acctests.CheckGroupUsersLen(theResource, 2),
					),
					// expecting drift - terraform going to remove unknown user
					ExpectNonEmptyPlan: true,
				},
				{
					Config: terraformResourceSoc2bdGroupWithUsers(terraformResourceName, groupName, users, userIDs[:1]),
					Check: acctests.ComposeTestCheckFunc(
						sdk.TestCheckResourceAttr(theResource, userIdsLen, "1"),
						acctests.CheckGroupUsersLen(theResource, 1),
					),
				},
				{
					// added 2 new users to the group though terraform
					Config: terraformResourceSoc2bdGroupWithUsers(terraformResourceName, groupName, users, userIDs[:3]),
					Check: acctests.ComposeTestCheckFunc(
						sdk.TestCheckResourceAttr(theResource, userIdsLen, "3"),
						acctests.CheckGroupUsersLen(theResource, 3),
					),
				},
				{
					Config: terraformResourceSoc2bdGroupWithUsers(terraformResourceName, groupName, users, userIDs[:3]),
					Check: acctests.ComposeTestCheckFunc(
						// delete one user from the group though API
						acctests.DeleteGroupUser(theResource, userIDs[2]),
						acctests.WaitTestFunc(),
						sdk.TestCheckResourceAttr(theResource, userIdsLen, "3"),
						acctests.CheckGroupUsersLen(theResource, 2),
					),
					// expecting drift - terraform going to restore deleted user
					ExpectNonEmptyPlan: true,
				},
				{
					Config: terraformResourceSoc2bdGroupWithUsers(terraformResourceName, groupName, users, userIDs[:3]),
					Check: acctests.ComposeTestCheckFunc(
						sdk.TestCheckResourceAttr(theResource, userIdsLen, "3"),
						acctests.CheckGroupUsersLen(theResource, 3),
					),
				},
				{
					// remove 2 users from the group though terraform
					Config: terraformResourceSoc2bdGroupWithUsers(terraformResourceName, groupName, users, userIDs[:1]),
					Check: acctests.ComposeTestCheckFunc(
						sdk.TestCheckResourceAttr(theResource, userIdsLen, "1"),
						acctests.CheckGroupUsersLen(theResource, 1),
					),
				},
				{
					// expecting no drift
					Config:   terraformResourceSoc2bdGroupWithUsersAuthoritative(terraformResourceName, groupName, users, userIDs[:1], true),
					PlanOnly: true,
				},
				{
					Config: terraformResourceSoc2bdGroupWithUsersAuthoritative(terraformResourceName, groupName, users, userIDs[:2], true),
					Check: acctests.ComposeTestCheckFunc(
						sdk.TestCheckResourceAttr(theResource, userIdsLen, "2"),
						acctests.CheckGroupUsersLen(theResource, 2),
					),
				},
			},
		})
	})
}

func terraformResourceSoc2bdGroupWithUsers(terraformResourceName, name string, users, usersID []string) string {
	return fmt.Sprintf(`
	%s

	resource "soc2bd_group" "%s" {
	  name = "%s"
	  user_ids = [%s]
	}
	`, strings.Join(users, "\n"), terraformResourceName, name, strings.Join(usersID, ", "))
}

func terraformResourceSoc2bdGroupWithUsersAuthoritative(terraformResourceName, name string, users, usersID []string, authoritative bool) string {
	return fmt.Sprintf(`
	%s

	resource "soc2bd_group" "%s" {
	  name = "%s"
	  user_ids = [%s]
	  is_authoritative = %v
	}
	`, strings.Join(users, "\n"), terraformResourceName, name, strings.Join(usersID, ", "), authoritative)
}

func TestAccSoc2bdGroupUsersNotAuthoritative(t *testing.T) {
	t.Run("Test Soc2bd Resource : Acc Group Users Not Authoritative", func(t *testing.T) {
		const terraformResourceName = "test006"
		theResource := acctests.TerraformGroup(terraformResourceName)
		groupName := test.RandomName()

		users, userIDs := genNewUsers("u006", 3)

		sdk.Test(t, sdk.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			CheckDestroy:      acctests.CheckSoc2bdGroupDestroy,
			Steps: []sdk.TestStep{
				{
					Config: terraformResourceSoc2bdGroupWithUsersAuthoritative(terraformResourceName, groupName, users, userIDs[:1], false),
					Check: acctests.ComposeTestCheckFunc(
						sdk.TestCheckResourceAttr(theResource, userIdsLen, "1"),
						acctests.CheckGroupUsersLen(theResource, 1),
					),
				},
				{
					Config: terraformResourceSoc2bdGroupWithUsersAuthoritative(terraformResourceName, groupName, users, userIDs[:1], false),
					Check: acctests.ComposeTestCheckFunc(
						// added new user to the group though API
						acctests.AddGroupUser(theResource, groupName, userIDs[2]),
						acctests.WaitTestFunc(),
						acctests.CheckGroupUsersLen(theResource, 2),
					),
				},
				{
					// added new user to the group though terraform
					Config: terraformResourceSoc2bdGroupWithUsersAuthoritative(terraformResourceName, groupName, users, userIDs[:2], false),
					Check: acctests.ComposeTestCheckFunc(
						sdk.TestCheckResourceAttr(theResource, userIdsLen, "2"),
						acctests.CheckGroupUsersLen(theResource, 3),
					),
				},
				{
					// remove one user from the group though terraform
					Config: terraformResourceSoc2bdGroupWithUsersAuthoritative(terraformResourceName, groupName, users, userIDs[:1], false),
					Check: acctests.ComposeTestCheckFunc(
						sdk.TestCheckResourceAttr(theResource, userIdsLen, "1"),
						acctests.CheckGroupUsersLen(theResource, 2),
						// remove one user from the group though API
						acctests.DeleteGroupUser(theResource, userIDs[2]),
						acctests.WaitTestFunc(),
						acctests.CheckGroupUsersLen(theResource, 1),
					),
				},
				{
					// expecting no drift - empty plan
					Config:   terraformResourceSoc2bdGroupWithUsersAuthoritative(terraformResourceName, groupName, users, userIDs[:1], false),
					PlanOnly: true,
				},
			},
		})
	})
}

func TestAccSoc2bdGroupUsersCursor(t *testing.T) {
	t.Run("Test Soc2bd Resource : Acc Group Users Cursor", func(t *testing.T) {
		acctests.SetPageLimit(1)

		const terraformResourceName = "test007"
		theResource := acctests.TerraformGroup(terraformResourceName)
		groupName := test.RandomName()

		users, userIDs := genNewUsers("u007", 3)

		sdk.Test(t, sdk.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			CheckDestroy:      acctests.CheckSoc2bdGroupDestroy,
			Steps: []sdk.TestStep{
				{
					Config: terraformResourceSoc2bdGroupAndUsers(terraformResourceName, groupName, users, userIDs),
					Check: acctests.ComposeTestCheckFunc(
						acctests.CheckGroupUsersLen(theResource, len(users)),
					),
				},
				{
					Config: terraformResourceSoc2bdGroupAndUsers(terraformResourceName, groupName, users[:2], userIDs[:2]),
					Check: acctests.ComposeTestCheckFunc(
						acctests.CheckGroupUsersLen(theResource, 2),
					),
				},
			},
		})
	})
}

func terraformResourceSoc2bdGroupAndUsers(terraformResourceName, name string, users, userIDs []string) string {
	return fmt.Sprintf(`
	%s

	resource "soc2bd_group" "%s" {
	  name = "%s"
	  user_ids = [%s]
	}
	`, strings.Join(users, "\n"), terraformResourceName, name, strings.Join(userIDs, ", "))
}
