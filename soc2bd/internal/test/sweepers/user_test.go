package sweepers

import (
	"context"

	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/client"
	soc2bd "github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/provider/resource"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	name := soc2bd.Soc2bdUser
	resource.AddTestSweepers(name, &resource.Sweeper{
		Name: name,
		F: newTestSweeper(name,
			func(client *client.Client, ctx context.Context) ([]Resource, error) {
				resources, err := client.ReadUsers(ctx)
				if err != nil {
					return nil, err
				}

				items := make([]Resource, 0, len(resources))
				for _, r := range resources {
					items = append(items, r)
				}
				return items, nil
			},
			func(client *client.Client, ctx context.Context, id string) error {
				return client.DeleteUser(ctx, id)
			},
		),
	})
}
