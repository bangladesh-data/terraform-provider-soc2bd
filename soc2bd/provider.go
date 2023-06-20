package soc2bd

import (
	"context"
	"fmt"
	"time"

	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/attr"
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/client"
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/provider/datasource"
	"github.com/bangladesh-data/terraform-provider-soc2bd/soc2bd/internal/provider/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	DefaultHTTPTimeout  = "10"
	DefaultHTTPMaxRetry = "10"
	DefaultURL          = "soc2bd.com"

	// EnvAPIToken env var for Token.
	EnvAPIToken     = "SOC2BD_API_TOKEN" //#nosec
	EnvNetwork      = "SOC2BD_NETWORK"
	EnvURL          = "SOC2BD_URL"
	EnvHTTPTimeout  = "SOC2BD_HTTP_TIMEOUT"
	EnvHTTPMaxRetry = "SOC2BD_HTTP_MAX_RETRY"
)

func Provider(version string) *schema.Provider {
	provider := &schema.Provider{
		Schema: providerOptions(),
		ResourcesMap: map[string]*schema.Resource{
			resource.Soc2bdRemoteNetwork:     resource.RemoteNetwork(),
			resource.Soc2bdConnector:         resource.Connector(),
			resource.Soc2bdConnectorTokens:   resource.ConnectorTokens(),
			resource.Soc2bdGroup:             resource.Group(),
			resource.Soc2bdResource:          resource.Resource(),
			resource.Soc2bdServiceAccount:    resource.ServiceAccount(),
			resource.Soc2bdServiceAccountKey: resource.ServiceKey(),
			resource.Soc2bdUser:              resource.User(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			datasource.Soc2bdGroup:            datasource.Group(),
			datasource.Soc2bdGroups:           datasource.Groups(),
			datasource.Soc2bdRemoteNetwork:    datasource.RemoteNetwork(),
			datasource.Soc2bdRemoteNetworks:   datasource.RemoteNetworks(),
			datasource.Soc2bdUser:             datasource.User(),
			datasource.Soc2bdUsers:            datasource.Users(),
			datasource.Soc2bdConnector:        datasource.Connector(),
			datasource.Soc2bdConnectors:       datasource.Connectors(),
			datasource.Soc2bdResource:         datasource.Resource(),
			datasource.Soc2bdResources:        datasource.Resources(),
			datasource.Soc2bdServiceAccounts:  datasource.ServiceAccounts(),
			datasource.Soc2bdSecurityPolicy:   datasource.SecurityPolicy(),
			datasource.Soc2bdSecurityPolicies: datasource.SecurityPolicies(),
		},
	}
	provider.ConfigureContextFunc = configure(version, provider)

	return provider
}

func providerOptions() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		attr.APIToken: {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			DefaultFunc: schema.EnvDefaultFunc(EnvAPIToken, nil),
			Description: fmt.Sprintf("The access key for API operations. You can retrieve this\n"+
				"from the Soc2bd Admin Console ([documentation](https://docs.soc2bd.com/docs/api-overview)).\n"+
				"Alternatively, this can be specified using the %s environment variable.", EnvAPIToken),
		},
		attr.Network: {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   false,
			DefaultFunc: schema.EnvDefaultFunc(EnvNetwork, nil),
			Description: fmt.Sprintf("Your Soc2bd network ID for API operations.\n"+
				"You can find it in the Admin Console URL, for example:\n"+
				"`autoco.soc2bd.com`, where `autoco` is your network ID\n"+
				"Alternatively, this can be specified using the %s environment variable.", EnvNetwork),
		},
		attr.URL: {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   false,
			DefaultFunc: schema.EnvDefaultFunc(EnvURL, DefaultURL),
			Description: fmt.Sprintf("The default is '%s'\n"+
				"This is optional and shouldn't be changed under normal circumstances.", DefaultURL),
		},
		attr.HTTPTimeout: {
			Type:        schema.TypeInt,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc(EnvHTTPTimeout, DefaultHTTPTimeout),
			Description: fmt.Sprintf("Specifies a time limit in seconds for the http requests made. The default value is %s seconds.\n"+
				"Alternatively, this can be specified using the %s environment variable", DefaultHTTPTimeout, EnvHTTPTimeout),
		},
		attr.HTTPMaxRetry: {
			Type:        schema.TypeInt,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc(EnvHTTPMaxRetry, DefaultHTTPMaxRetry),
			Description: fmt.Sprintf("Specifies a retry limit for the http requests made. The default value is %s.\n"+
				"Alternatively, this can be specified using the %s environment variable", DefaultHTTPMaxRetry, EnvHTTPMaxRetry),
		},
	}
}

func configure(version string, _ *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		apiToken := d.Get(attr.APIToken).(string)
		network := d.Get(attr.Network).(string)
		url := d.Get(attr.URL).(string)
		httpTimeout := d.Get(attr.HTTPTimeout).(int)
		httpMaxRetry := d.Get(attr.HTTPMaxRetry).(int)

		if network != "" {
			return client.NewClient(url,
					apiToken,
					network,
					time.Duration(httpTimeout)*time.Second,
					httpMaxRetry,
					version),
				nil
		}

		return nil, diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create Soc2bd client",
				Detail:   "Unable to create anonymous Soc2bd client, network has to be provided",
			},
		}
	}
}
