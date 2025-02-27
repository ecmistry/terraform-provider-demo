package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: Provider,
	})
}

// Provider returns a terraform-plugin-sdk/v2/helper/schema.Provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"management_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GRAVITEE_MANAGEMENT_URL", nil),
				Description: "URL of the Gravitee Management API",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GRAVITEE_USERNAME", nil),
				Description: "Username for Gravitee API Management",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("GRAVITEE_PASSWORD", nil),
				Description: "Password for Gravitee API Management",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"gravitee_api":          resourceGraviteeAPI(),
			"gravitee_plan":         resourceGraviteePlan(),
			"gravitee_subscription": resourceGraviteeSubscription(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"gravitee_api": dataSourceGraviteeAPI(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

// providerConfigure configures the provider with auth credentials and API client
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	managementURL := d.Get("management_url").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	// Initialize the client
	client := &gravitee.Client{
		ManagementURL: managementURL,
		Username:      username,
		Password:      password,
		HTTPClient:    &http.Client{},
	}

	// Test the connection
	err := client.TestConnection()
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return client, nil
}