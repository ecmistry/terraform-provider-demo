// client_subscription.go
package gravitee

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Subscription represents a Gravitee API subscription
type Subscription struct {
	ID                    string                  `json:"id,omitempty"`
	PlanID                string                  `json:"planId"`
	ApplicationID         string                  `json:"applicationId"`
	Status                string                  `json:"status,omitempty"`
	ConsumerConfiguration *ConsumerConfiguration  `json:"consumerConfiguration,omitempty"`
	Metadata              map[string]string       `json:"metadata,omitempty"`
}

// ConsumerConfiguration represents subscription consumer configuration
type ConsumerConfiguration struct {
	EntrypointID            string                   `json:"entrypointId"`
	Channel                 string                   `json:"channel,omitempty"`
	EntrypointConfiguration *EntrypointConfiguration `json:"entrypointConfiguration,omitempty"`
}

// EntrypointConfiguration represents entrypoint configuration
type EntrypointConfiguration struct {
	CallbackURL string   `json:"callbackUrl"`
	Headers     []Header `json:"headers,omitempty"`
}

// Header represents a HTTP header
type Header struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Create a Subscription
func (c *Client) CreateSubscription(apiID string, subscription *Subscription) (*Subscription, error) {
	body, err := json.Marshal(subscription)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/management/v2/environments/DEFAULT/apis/%s/subscriptions", c.ManagementURL, apiID), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.Username, c.Password)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	var createdSubscription Subscription
	if err := json.NewDecoder(resp.Body).Decode(&createdSubscription); err != nil {
		return nil, err
	}

	return &createdSubscription, nil
}

// Get a Subscription by ID
func (c *Client) GetSubscription(apiID string, subscriptionID string) (*Subscription, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/management/v2/environments/DEFAULT/apis/%s/subscriptions/%s", c.ManagementURL, apiID, subscriptionID), nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.Username, c.Password)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	var subscription Subscription
	if err := json.NewDecoder(resp.Body).Decode(&subscription); err != nil {
		return nil, err
	}

	return &subscription, nil
}

// Update a Subscription
func (c *Client) UpdateSubscription(apiID string, subscription *Subscription) error {
	body, err := json.Marshal(map[string]interface{}{
		"configuration": subscription.ConsumerConfiguration,
		"metadata":      subscription.Metadata,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/management/v2/environments/DEFAULT/apis/%s/subscriptions/%s", c.ManagementURL, apiID, subscription.ID), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.Username, c.Password)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	return nil
}

// Close a Subscription
func (c *Client) CloseSubscription(apiID string, subscriptionID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/management/v2/environments/DEFAULT/apis/%s/subscriptions/%s", c.ManagementURL, apiID, subscriptionID), nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.Username, c.Password)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	return nil
}

// Transfer a Subscription to a new plan
func (c *Client) TransferSubscription(apiID string, subscriptionID string, newPlanID string) error {
	body, err := json.Marshal(map[string]string{
		"plan": newPlanID,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/management/v2/environments/DEFAULT/apis/%s/subscriptions/%s/_transfer", c.ManagementURL, apiID, subscriptionID), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.Username, c.Password)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	return nil
}

// Accept a Subscription
func (c *Client) AcceptSubscription(apiID string, subscriptionID string, reason string) error {
	body, err := json.Marshal(map[string]string{
		"reason": reason,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/management/v2/environments/DEFAULT/apis/%s/subscriptions/%s/_accept", c.ManagementURL, apiID, subscriptionID), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.Username, c.Password)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	return nil
}