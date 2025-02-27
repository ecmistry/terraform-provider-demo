// client_plan.go
package gravitee

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Plan represents a Gravitee API plan
type Plan struct {
	ID                string      `json:"id,omitempty"`
	Name              string      `json:"name"`
	Description       string      `json:"description"`
	DefinitionVersion string      `json:"definitionVersion"`
	Mode              string      `json:"mode"`
	Security          *Security   `json:"security"`
	Characteristics   []string    `json:"characteristics,omitempty"`
	Validation        string      `json:"validation,omitempty"`
	Status            string      `json:"status,omitempty"`
}

// Security represents a plan's security configuration
type Security struct {
	Type string `json:"type"`
}

// Create a Plan for an API
func (c *Client) CreatePlan(apiID string, plan *Plan) (*Plan, error) {
	body, err := json.Marshal(plan)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/management/v2/environments/DEFAULT/apis/%s/plans", c.ManagementURL, apiID), bytes.NewBuffer(body))
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

	var createdPlan Plan
	if err := json.NewDecoder(resp.Body).Decode(&createdPlan); err != nil {
		return nil, err
	}

	return &createdPlan, nil
}

// Get a Plan by ID
func (c *Client) GetPlan(apiID string, planID string) (*Plan, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/management/v2/environments/DEFAULT/apis/%s/plans/%s", c.ManagementURL, apiID, planID), nil)
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

	var plan Plan
	if err := json.NewDecoder(resp.Body).Decode(&plan); err != nil {
		return nil, err
	}

	return &plan, nil
}

// Update a Plan
func (c *Client) UpdatePlan(apiID string, plan *Plan) error {
	body, err := json.Marshal(plan)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/management/v2/environments/DEFAULT/apis/%s/plans/%s", c.ManagementURL, apiID, plan.ID), bytes.NewBuffer(body))
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

// Delete a Plan
func (c *Client) DeletePlan(apiID string, planID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/management/v2/environments/DEFAULT/apis/%s/plans/%s", c.ManagementURL, apiID, planID), nil)
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

// Publish a Plan
func (c *Client) PublishPlan(apiID string, planID string) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/management/v2/environments/DEFAULT/apis/%s/plans/%s/_publish", c.ManagementURL, apiID, planID), nil)
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