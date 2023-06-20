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

func datasourceRemoteNetworkRead(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	networkID := resourceData.Get(attr.ID).(string)
	networkName := resourceData.Get(attr.Name).(string)

	network, err := c.ReadRemoteNetwork(ctx, networkID, networkName)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := resourceData.Set(attr.Name, network.Name); err != nil {
		return diag.FromErr(err)
	}

	resourceData.SetId(network.ID)

	return nil
}

func RemoteNetwork() *schema.Resource {
	return &schema.Resource{
		Description: "A Remote Network represents a single private network in Soc2bd that can have one or more Connectors and Resources assigned to it. You must create a Remote Network before creating Resources and Connectors that belong to it. For more information, see Soc2bd's [documentation](https://docs.soc2bd.com/docs/remote-networks).",
		ReadContext: datasourceRemoteNetworkRead,
		Schema: map[string]*schema.Schema{
			attr.ID: {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The ID of the Remote Network",
				ExactlyOneOf: []string{attr.Name},
			},
			attr.Name: {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The name of the Remote Network",
				ExactlyOneOf: []string{attr.ID},
			},
			attr.Location: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: fmt.Sprintf("The location of the Remote Network. Must be one of the following: %s.", strings.Join(model.Locations, ", ")),
			},
		},
	}
}
