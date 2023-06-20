package sweepers

import (
	"context"

	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/client"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const resourceRemoteNetwork = "soc2bd_remote_network"

func init() {
	resource.AddTestSweepers(resourceRemoteNetwork, &resource.Sweeper{
		Name: resourceRemoteNetwork,
		F: newTestSweeper(resourceRemoteNetwork,
			func(client *client.Client, ctx context.Context) ([]Resource, error) {
				resources, err := client.ReadRemoteNetworks(ctx)
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
				return client.DeleteRemoteNetwork(ctx, id)
			},
		),
	})
}
