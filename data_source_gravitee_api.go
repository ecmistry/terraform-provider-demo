// data_source_gravitee_api.go
package gravitee

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGraviteeAPI() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGraviteeAPIRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the API",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the API",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the API",
			},
			"api_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Version of the API",
			},
			"definition_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Definition version of the API",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the API (PROXY, MESSAGE)",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State of the API (STARTED, STOPPED)",
			},
			"created_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Creation timestamp of the API",
			},
			"updated_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Last update timestamp of the API",
			},
			"deployed_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Last deployment timestamp of the API",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func dataSourceGraviteeAPIRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*gravitee.Client)

	apiID := d.Get("id").(string)

	api, err := client.GetAPI(apiID)
	if err != nil {
		return diag.FromErr(err)
	}