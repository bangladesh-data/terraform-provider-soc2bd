package datasource

import (
	"context"
	"fmt"
	"strings"

	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/attr"
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/client"
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/model"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceRemoteNetworksRead(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)

	remoteNetworks, err := client.ReadRemoteNetworks(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := resourceData.Set(attr.RemoteNetworks, convertRemoteNetworksToTerraform(remoteNetworks)); err != nil {
		return diag.FromErr(err)
	}

	resourceData.SetId("all-remote-networks")

	return nil
}

func RemoteNetworks() *schema.Resource {
	return &schema.Resource{
		Description: "A Remote Network represents a single private network in Soc2bd that can have one or more Connectors and Resources assigned to it. You must create a Remote Network before creating Resources and Connectors that belong to it. For more information, see Soc2bd's [documentation](https://docs.soc2bd.com/docs/remote-networks).",
		ReadContext: datasourceRemoteNetworksRead,
		Schema: map[string]*schema.Schema{
			attr.RemoteNetworks: {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of Remote Networks",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						attr.ID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the Remote Network",
						},
						attr.Name: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the Remote Network",
						},
						attr.Location: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: fmt.Sprintf("The location of the Remote Network. Must be one of the following: %s.", strings.Join(model.Locations, ", ")),
						},
					},
				},
			},
		},
	}
}
