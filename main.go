package main

import (
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

var (
	version = "dev"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return soc2bd.Provider(version)
		},
	})
}
