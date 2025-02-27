// resource_gravitee_plan.go
package gravitee

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGraviteePlan() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGraviteePlanCreate,
		ReadContext:   resourceGraviteePlanRead,
		UpdateContext: resourceGraviteePlanUpdate,
		DeleteContext: resourceGraviteePlanDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"api_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the API this plan belongs to",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the plan",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the plan",
			},
			"definition_version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Definition version (e.g., V4)",
			},
			"security_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Security type for the plan (e.g., KEY_LESS, subscription)",
			},
			"mode": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Mode of the plan (e.g., STANDARD, PUSH)",
			},
			"characteristics": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Characteristics of the plan",
			},
			"validation": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Validation mode for the plan",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the plan",
			},
			"auto_publish": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
// resource_gravitee_plan.go (continued)
				Description: "Whether to automatically publish the plan",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceGraviteePlanCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*gravitee.Client)
	apiID := d.Get("api_id").(string)

	plan := expandPlan(d)

	createdPlan, err := client.CreatePlan(apiID, plan)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdPlan.ID)

	// Publish the plan if auto_publish is enabled
	if d.Get("auto_publish").(bool) {
		err = client.PublishPlan(apiID, createdPlan.ID)
		if err != nil {
			return diag.FromErr(err)
		}
		d.Set("status", "PUBLISHED")
	}

	return resourceGraviteePlanRead(ctx, d, m)
}

func resourceGraviteePlanRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*gravitee.Client)
	apiID := d.Get("api_id").(string)

	plan, err := client.GetPlan(apiID, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if plan == nil {
		d.SetId("")
		return nil
	}

	// Flatten the Plan object and set to ResourceData
	if err := flattenPlan(d, plan); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceGraviteePlanUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*gravitee.Client)
	apiID := d.Get("api_id").(string)

	plan := expandPlan(d)
	plan.ID = d.Id()

	err := client.UpdatePlan(apiID, plan)
	if err != nil {
		return diag.FromErr(err)
	}

	// Publish the plan if auto_publish is enabled and there are changes that require republishing
	if d.Get("auto_publish").(bool) && d.HasChange("name") || d.HasChange("description") || d.HasChange("security_type") || d.HasChange("mode") || d.HasChange("characteristics") {
		err = client.PublishPlan(apiID, plan.ID)
		if err != nil {
			return diag.FromErr(err)
		}
		d.Set("status", "PUBLISHED")
	}

	return resourceGraviteePlanRead(ctx, d, m)
}

func resourceGraviteePlanDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*gravitee.Client)
	apiID := d.Get("api_id").(string)

	err := client.DeletePlan(apiID, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

// Helper functions for expanding and flattening Plan objects
func expandPlan(d *schema.ResourceData) *gravitee.Plan {
	securityType := d.Get("security_type").(string)

	plan := &gravitee.Plan{
		Name:              d.Get("name").(string),
		Description:       d.Get("description").(string),
		DefinitionVersion: d.Get("definition_version").(string),
		Mode:              d.Get("mode").(string),
		Security: &gravitee.Security{
			Type: securityType,
		},
	}

	if v, ok := d.GetOk("characteristics"); ok {
		characteristicsRaw := v.([]interface{})
		characteristics := make([]string, len(characteristicsRaw))
		for i, v := range characteristicsRaw {
			characteristics[i] = v.(string)
		}
		plan.Characteristics = characteristics
	}

	if v, ok := d.GetOk("validation"); ok {
		plan.Validation = v.(string)
	}

	return plan
}

func flattenPlan(d *schema.ResourceData, plan *gravitee.Plan) error {
	d.Set("name", plan.Name)
	d.Set("description", plan.Description)
	d.Set("definition_version", plan.DefinitionVersion)
	d.Set("mode", plan.Mode)

	if plan.Security != nil {
		d.Set("security_type", plan.Security.Type)
	}

	if err := d.Set("characteristics", plan.Characteristics); err != nil {
		return err
	}

	if plan.Validation != "" {
		d.Set("validation", plan.Validation)
	}

	d.Set("status", plan.Status)

	return nil
}