// resource_gravitee_subscription.go
package gravitee

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGraviteeSubscription() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGraviteeSubscriptionCreate,
		ReadContext:   resourceGraviteeSubscriptionRead,
		UpdateContext: resourceGraviteeSubscriptionUpdate,
		DeleteContext: resourceGraviteeSubscriptionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"api_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the API this subscription belongs to",
			},
			"application_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the application subscribing to the API",
			},
			"plan_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the plan to subscribe to",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the subscription",
			},
			"auto_validate": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to automatically validate the subscription",
			},
			"consumer_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"entrypoint_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"channel": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"entrypoint_configuration": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"callback_url": {
										Type:     schema.TypeString,
										Required: true,
									},
									"headers": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Required: true,
												},
												"value": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Metadata for the subscription",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceGraviteeSubscriptionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*gravitee.Client)
	apiID := d.Get("api_id").(string)

	subscription := expandSubscription(d)

	createdSubscription, err := client.CreateSubscription(apiID, subscription)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdSubscription.ID)

	// Auto-validate the subscription if enabled
	if d.Get("auto_validate").(bool) && createdSubscription.Status != "ACCEPTED" {
		err = client.AcceptSubscription(apiID, createdSubscription.ID, "Auto-approved by Terraform")
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceGraviteeSubscriptionRead(ctx, d, m)
}

func resourceGraviteeSubscriptionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*gravitee.Client)
	apiID := d.Get("api_id").(string)

	subscription, err := client.GetSubscription(apiID, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if subscription == nil {
		d.SetId("")
		return nil
	}

	// Flatten the Subscription object and set to ResourceData
	if err := flattenSubscription(d, subscription); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceGraviteeSubscriptionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*gravitee.Client)
	apiID := d.Get("api_id").(string)

	// Check if consumer configuration or metadata is being changed
	if d.HasChange("consumer_configuration") || d.HasChange("metadata") {
		subscription := expandSubscription(d)
		subscription.ID = d.Id()

		err := client.UpdateSubscription(apiID, subscription)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Check if plan is being changed
	if d.HasChange("plan_id") {
		// Transfer to a new plan
		err := client.TransferSubscription(apiID, d.Id(), d.Get("plan_id").(string))
		if err != nil {
			return diag.FromErr(err)
		}

		// Auto-validate the transferred subscription if enabled
		if d.Get("auto_validate").(bool) {
			err = client.AcceptSubscription(apiID, d.Id(), "Auto-approved transfer by Terraform")
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return resourceGraviteeSubscriptionRead(ctx, d, m)
}

func resourceGraviteeSubscriptionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*gravitee.Client)
	apiID := d.Get("api_id").(string)

	err := client.CloseSubscription(apiID, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

// Helper functions for expanding and flattening Subscription objects
func expandSubscription(d *schema.ResourceData) *gravitee.Subscription {
	subscription := &gravitee.Subscription{
		PlanID:        d.Get("plan_id").(string),
		ApplicationID: d.Get("application_id").(string),
	}

	if v, ok := d.GetOk("consumer_configuration"); ok && len(v.([]interface{})) > 0 {
		subscription.ConsumerConfiguration = expandConsumerConfiguration(v.([]interface{})[0].(map[string]interface{}))
	}

	if v, ok := d.GetOk("metadata"); ok {
		metadataMap := v.(map[string]interface{})
		metadata := make(map[string]string, len(metadataMap))
		for k, v := range metadataMap {
			metadata[k] = v.(string)
		}
		subscription.Metadata = metadata
	}

	return subscription
}

func expandConsumerConfiguration(config map[string]interface{}) *gravitee.ConsumerConfiguration {
	consumerConfig := &gravitee.ConsumerConfiguration{
		EntrypointID: config["entrypoint_id"].(string),
	}

	if v, ok := config["channel"]; ok {
		consumerConfig.Channel = v.(string)
	}

	if v, ok := config["entrypoint_configuration"]; ok && len(v.([]interface{})) > 0 {
		entrypointConfig := v.([]interface{})[0].(map[string]interface{})

		consumerConfig.EntrypointConfiguration = &gravitee.EntrypointConfiguration{
			CallbackURL: entrypointConfig["callback_url"].(string),
		}

		if headersRaw, ok := entrypointConfig["headers"]; ok {
			headers := make([]gravitee.Header, 0)
			for _, h := range headersRaw.([]interface{}) {
				header := h.(map[string]interface{})
				headers = append(headers, gravitee.Header{
					Name:  header["name"].(string),
					Value: header["value"].(string),
				})
			}
			consumerConfig.EntrypointConfiguration.Headers = headers
		}
	}

	return consumerConfig
}

func flattenSubscription(d *schema.ResourceData, subscription *gravitee.Subscription) error {
	d.Set("plan_id", subscription.PlanID)
	d.Set("application_id", subscription.ApplicationID)
	d.Set("status", subscription.Status)

	if subscription.ConsumerConfiguration != nil {
		consumerConfig := flattenConsumerConfiguration(subscription.ConsumerConfiguration)
		if err := d.Set("consumer_configuration", []interface{}{consumerConfig}); err != nil {
			return err
		}
	}

	if subscription.Metadata != nil {
		if err := d.Set("metadata", subscription.Metadata); err != nil {
			return err
		}
	}

	return nil
}

func flattenConsumerConfiguration(config *gravitee.ConsumerConfiguration) map[string]interface{} {
	result := map[string]interface{}{
		"entrypoint_id": config.EntrypointID,
	}

	if config.Channel != "" {
		result["channel"] = config.Channel
	}

	if config.EntrypointConfiguration != nil {
		entrypointConfig := map[string]interface{}{
			"callback_url": config.EntrypointConfiguration.CallbackURL,
		}

		if len(config.EntrypointConfiguration.Headers) > 0 {
			headers := make([]map[string]interface{}, len(config.EntrypointConfiguration.Headers))
			for i, header := range config.EntrypointConfiguration.Headers {
				headers[i] = map[string]interface{}{
					"name":  header.Name,
					"value": header.Value,
				}
			}
			entrypointConfig["headers"] = headers
		}

		result["entrypoint_configuration"] = []interface{}{entrypointConfig}
	}

	return result
}