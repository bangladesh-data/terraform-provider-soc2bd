package resource

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/attr"
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/model"
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/provider/resource"
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/test"
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/test/acctests"
	sdk "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSoc2bdUserCreateUpdate(t *testing.T) {
	t.Run("Test Soc2bd Resource : Acc User Create/Update", func(t *testing.T) {
		const terraformResourceName = "test001"
		theResource := acctests.TerraformUser(terraformResourceName)
		email := test.RandomEmail()
		firstName := test.RandomName()
		lastName := test.RandomName()
		role := model.UserRoleSupport

		sdk.Test(t, sdk.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			CheckDestroy:      acctests.CheckSoc2bdUserDestroy,
			Steps: []sdk.TestStep{
				{
					Config: terraformResourceSoc2bdUser(terraformResourceName, email),
					Check: acctests.ComposeTestCheckFunc(
						acctests.CheckSoc2bdResourceExists(theResource),
						sdk.TestCheckResourceAttr(theResource, attr.Email, email),
					),
				},
				{
					Config: terraformResourceSoc2bdUserWithFirstName(terraformResourceName, email, firstName),
					Check: acctests.ComposeTestCheckFunc(
						acctests.CheckSoc2bdResourceExists(theResource),
						sdk.TestCheckResourceAttr(theResource, attr.Email, email),
						sdk.TestCheckResourceAttr(theResource, attr.FirstName, firstName),
					),
				},
				{
					Config: terraformResourceSoc2bdUserWithLastName(terraformResourceName, email, lastName),
					Check: acctests.ComposeTestCheckFunc(
						acctests.CheckSoc2bdResourceExists(theResource),
						sdk.TestCheckResourceAttr(theResource, attr.Email, email),
						sdk.TestCheckResourceAttr(theResource, attr.FirstName, firstName),
						sdk.TestCheckResourceAttr(theResource, attr.LastName, lastName),
					),
				},
				{
					Config: terraformResourceSoc2bdUserWithRole(terraformResourceName, email, role),
					Check: acctests.ComposeTestCheckFunc(
						acctests.CheckSoc2bdResourceExists(theResource),
						sdk.TestCheckResourceAttr(theResource, attr.Email, email),
						sdk.TestCheckResourceAttr(theResource, attr.FirstName, firstName),
						sdk.TestCheckResourceAttr(theResource, attr.LastName, lastName),
						sdk.TestCheckResourceAttr(theResource, attr.Role, role),
					),
				},
			},
		})
	})
}

func terraformResourceSoc2bdUser(terraformResourceName, email string) string {
	return fmt.Sprintf(`
	resource "soc2bd_user" "%s" {
	  email = "%s"
	  send_invite = false
	}
	`, terraformResourceName, email)
}

func terraformResourceSoc2bdUserWithFirstName(terraformResourceName, email, firstName string) string {
	return fmt.Sprintf(`
	resource "soc2bd_user" "%s" {
	  email = "%s"
	  first_name = "%s"
	  send_invite = false
	}
	`, terraformResourceName, email, firstName)
}

func terraformResourceSoc2bdUserWithLastName(terraformResourceName, email, lastName string) string {
	return fmt.Sprintf(`
	resource "soc2bd_user" "%s" {
	  email = "%s"
	  last_name = "%s"
	  send_invite = false
	}
	`, terraformResourceName, email, lastName)
}

func terraformResourceSoc2bdUserWithRole(terraformResourceName, email, role string) string {
	return fmt.Sprintf(`
	resource "soc2bd_user" "%s" {
	  email = "%s"
	  role = "%s"
	  send_invite = false
	}
	`, terraformResourceName, email, role)
}

func TestAccSoc2bdUserFullCreate(t *testing.T) {
	t.Run("Test Soc2bd Resource : Acc User Full Create", func(t *testing.T) {
		const terraformResourceName = "test002"
		theResource := acctests.TerraformUser(terraformResourceName)
		email := test.RandomEmail()
		firstName := test.RandomName()
		lastName := test.RandomName()
		role := test.RandomUserRole()

		sdk.Test(t, sdk.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			CheckDestroy:      acctests.CheckSoc2bdUserDestroy,
			Steps: []sdk.TestStep{
				{
					Config: terraformResourceSoc2bdUserFull(terraformResourceName, email, firstName, lastName, role),
					Check: acctests.ComposeTestCheckFunc(
						acctests.CheckSoc2bdResourceExists(theResource),
						sdk.TestCheckResourceAttr(theResource, attr.Email, email),
						sdk.TestCheckResourceAttr(theResource, attr.FirstName, firstName),
						sdk.TestCheckResourceAttr(theResource, attr.LastName, lastName),
						sdk.TestCheckResourceAttr(theResource, attr.Role, role),
					),
				},
			},
		})
	})
}

func terraformResourceSoc2bdUserFull(terraformResourceName, email, firstName, lastName, role string) string {
	return fmt.Sprintf(`
	resource "soc2bd_user" "%s" {
	  email = "%s"
	  first_name = "%s"
	  last_name = "%s"
	  role = "%s"
	  send_invite = false
	}
	`, terraformResourceName, email, firstName, lastName, role)
}

func TestAccSoc2bdUserReCreation(t *testing.T) {
	t.Run("Test Soc2bd Resource : Acc User ReCreation", func(t *testing.T) {
		const terraformResourceName = "test003"
		theResource := acctests.TerraformUser(terraformResourceName)
		email1 := test.RandomEmail()
		email2 := test.RandomEmail()

		sdk.Test(t, sdk.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			CheckDestroy:      acctests.CheckSoc2bdUserDestroy,
			Steps: []sdk.TestStep{
				{
					Config: terraformResourceSoc2bdUser(terraformResourceName, email1),
					Check: acctests.ComposeTestCheckFunc(
						acctests.CheckSoc2bdResourceExists(theResource),
						sdk.TestCheckResourceAttr(theResource, attr.Email, email1),
					),
				},
				{
					Config: terraformResourceSoc2bdUser(terraformResourceName, email2),
					Check: acctests.ComposeTestCheckFunc(
						acctests.CheckSoc2bdResourceExists(theResource),
						sdk.TestCheckResourceAttr(theResource, attr.Email, email2),
					),
				},
			},
		})
	})
}

func TestAccSoc2bdUserUpdateState(t *testing.T) {
	t.Run("Test Soc2bd Resource : Acc User Update State", func(t *testing.T) {
		const terraformResourceName = "test004"
		theResource := acctests.TerraformUser(terraformResourceName)
		email := test.RandomEmail()

		sdk.Test(t, sdk.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			CheckDestroy:      acctests.CheckSoc2bdUserDestroy,
			Steps: []sdk.TestStep{
				{
					Config: terraformResourceSoc2bdUser(terraformResourceName, email),
					Check: acctests.ComposeTestCheckFunc(
						acctests.CheckSoc2bdResourceExists(theResource),
						sdk.TestCheckResourceAttr(theResource, attr.Email, email),
					),
				},
				{
					Config:      terraformResourceSoc2bdUserDisabled(terraformResourceName, email),
					ExpectError: regexp.MustCompile("User in PENDING state cannot be updated to the state: DISABLED"),
				},
			},
		})
	})
}

func terraformResourceSoc2bdUserDisabled(terraformResourceName, email string) string {
	return fmt.Sprintf(`
	resource "soc2bd_user" "%s" {
	  email = "%s"
	  send_invite = false
	  is_active = false
	}
	`, terraformResourceName, email)
}

func TestAccSoc2bdUserDelete(t *testing.T) {
	t.Run("Test Soc2bd Resource : Acc User Delete", func(t *testing.T) {
		const terraformResourceName = "test005"
		theResource := acctests.TerraformUser(terraformResourceName)

		sdk.Test(t, sdk.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			CheckDestroy:      acctests.CheckSoc2bdUserDestroy,
			Steps: []sdk.TestStep{
				{
					Config:  terraformResourceSoc2bdUser(terraformResourceName, test.RandomEmail()),
					Destroy: true,
					Check: acctests.ComposeTestCheckFunc(
						acctests.CheckSoc2bdResourceDoesNotExists(theResource),
					),
				},
			},
		})
	})
}

func TestAccSoc2bdUserReCreateAfterDeletion(t *testing.T) {
	t.Run("Test Soc2bd Resource : Acc User Create After Deletion", func(t *testing.T) {
		const terraformResourceName = "test006"
		theResource := acctests.TerraformUser(terraformResourceName)
		email := test.RandomEmail()

		sdk.Test(t, sdk.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			CheckDestroy:      acctests.CheckSoc2bdUserDestroy,
			Steps: []sdk.TestStep{
				{
					Config: terraformResourceSoc2bdUser(terraformResourceName, email),
					Check: acctests.ComposeTestCheckFunc(
						acctests.CheckSoc2bdResourceExists(theResource),
						acctests.DeleteSoc2bdResource(theResource, resource.Soc2bdUser),
					),
					ExpectNonEmptyPlan: true,
				},
				{
					Config: terraformResourceSoc2bdUser(terraformResourceName, email),
					Check: acctests.ComposeTestCheckFunc(
						acctests.CheckSoc2bdResourceExists(theResource),
					),
				},
			},
		})
	})
}

func TestAccSoc2bdUserCreateWithUnknownRole(t *testing.T) {
	t.Run("Test Soc2bd Resource : Acc User Create With Unknown Role", func(t *testing.T) {
		const terraformResourceName = "test007"

		sdk.Test(t, sdk.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			CheckDestroy:      acctests.CheckSoc2bdUserDestroy,
			Steps: []sdk.TestStep{
				{
					Config:      terraformResourceSoc2bdUserWithRole(terraformResourceName, test.RandomEmail(), "UnknownRole"),
					ExpectError: regexp.MustCompile(`Error: expected role to be one of \[ADMIN DEVOPS SUPPORT MEMBER\], got UnknownRole`),
				},
			},
		})
	})
}

func TestAccSoc2bdUserCreateWithoutEmail(t *testing.T) {
	t.Run("Test Soc2bd Resource : Acc User Create Without Email", func(t *testing.T) {
		const terraformResourceName = "test008"

		sdk.Test(t, sdk.TestCase{
			ProviderFactories: acctests.ProviderFactories,
			PreCheck:          func() { acctests.PreCheck(t) },
			CheckDestroy:      acctests.CheckSoc2bdUserDestroy,
			Steps: []sdk.TestStep{
				{
					Config:      terraformResourceSoc2bdUserWithoutEmail(terraformResourceName),
					ExpectError: regexp.MustCompile("Error: Missing required argument"),
				},
			},
		})
	})
}

func terraformResourceSoc2bdUserWithoutEmail(terraformResourceName string) string {
	return fmt.Sprintf(`
	resource "soc2bd_user" "%s" {
	  send_invite = false
	}
	`, terraformResourceName)
}

func genNewUsers(resourcePrefix string, count int) ([]string, []string) {
	users := make([]string, 0, count)
	userIDs := make([]string, 0, count)

	for i := 0; i < count; i++ {
		resourceName := fmt.Sprintf("%s_%d", resourcePrefix, i+1)
		users = append(users, terraformResourceSoc2bdUser(resourceName, test.RandomEmail()))
		userIDs = append(userIDs, fmt.Sprintf("soc2bd_user.%s.id", resourceName))
	}

	return users, userIDs
}
