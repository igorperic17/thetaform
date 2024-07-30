package provider

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type DeploymentTemplateRequest struct {
	ID              basetypes.StringValue `json:"id,omitempty" tfsdk:"id"`
	Name            string                `json:"name" tfsdk:"name"`
	ProjectID       string                `json:"project_id" tfsdk:"project_id"`
	Description     basetypes.StringValue `json:"description,omitempty" tfsdk:"description"`
	ContainerImages []string              `json:"container_images" tfsdk:"container_images"`
	ContainerPort   basetypes.Int64Value  `json:"container_port,omitempty" tfsdk:"container_port"`
	ContainerArgs   []string              `json:"container_args,omitempty" tfsdk:"container_args"`
	EnvVars         map[string]string     `json:"env_vars,omitempty" tfsdk:"env_vars"`
	Tags            []string              `json:"tags,omitempty" tfsdk:"tags"`
	IconURL         basetypes.StringValue `json:"icon_url,omitempty" tfsdk:"icon_url"`
	RequireEnvVars  basetypes.BoolValue   `json:"require_env_vars,omitempty" tfsdk:"require_env_vars"`
	Rank            basetypes.Int64Value  `json:"rank,omitempty" tfsdk:"rank"`
	CreateTime      basetypes.StringValue `json:"create_time,omitempty" tfsdk:"create_time"`
	Category        basetypes.StringValue `json:"category,omitempty" tfsdk:"category"`
}

type DeploymentTemplateRequestNative struct {
	Name           string            `json:"name"`
	ProjectID      string            `json:"project_id"`
	Description    string            `json:"description,omitempty"`
	ContainerImage []string          `json:"container_image"`
	ContainerPort  int64             `json:"container_port,omitempty"`
	ContainerArgs  []string          `json:"container_args,omitempty"`
	EnvVars        map[string]string `json:"env_vars,omitempty"`
	Tags           []string          `json:"tags,omitempty"`
	IconURL        string            `json:"icon_url,omitempty"`
	RequireEnvVars *bool             `json:"require_env_vars,omitempty"`
	Rank           *int64            `json:"rank,omitempty"`
}

type DeploymentTemplate struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	Description     string            `json:"description"`
	Tags            []string          `json:"tags"`
	Category        string            `json:"category"`
	ProjectID       string            `json:"project_id"`
	ContainerImages []string          `json:"container_images"`
	ContainerPort   int64             `json:"container_port"`
	ContainerArgs   []string          `json:"container_args"`
	EnvVars         map[string]string `json:"env_vars"`
	RequireEnvVars  *bool             `json:"require_env_vars"`
	Rank            *int64            `json:"rank"`
	IconURL         string            `json:"icon_url"`
	CreateTime      time.Time         `json:"create_time"`
}

type CreateDeploymentTemplateResponse struct {
	Status string             `json:"status"`
	Body   DeploymentTemplate `json:"body"`
}

type UpdateDeploymentTemplateResponse struct {
	Status string             `json:"status"`
	Body   DeploymentTemplate `json:"body"`
}

type DeleteDeploymentTemplateResponse struct {
	Status string `json:"status"`
	Body   bool   `json:"body"`
}

func (c *Client) CreateDeploymentTemplate(template DeploymentTemplateRequestNative) (*DeploymentTemplate, error) {
	url := "https://controller.thetaedgecloud.com/deployment_template"

	jsonData, err := json.Marshal(template)
	if err != nil {
		return nil, err
	}

	body, err := sendRequest(c, "POST", url, jsonData)
	if err != nil {
		return nil, err
	}

	var respData struct {
		Status string             `json:"status"`
		Body   DeploymentTemplate `json:"body"`
	}
	if err := json.Unmarshal(body, &respData); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &respData.Body, nil
}

func (c *Client) UpdateDeploymentTemplate(templateID string, template DeploymentTemplateRequestNative) (*DeploymentTemplate, error) {
	url := fmt.Sprintf("https://controller.thetaedgecloud.com/deployment_template/%s", templateID)

	jsonData, err := json.Marshal(template)
	if err != nil {
		return nil, err
	}

	body, err := sendRequest(c, "PUT", url, jsonData)
	if err != nil {
		return nil, err
	}

	var respData struct {
		Status string             `json:"status"`
		Body   DeploymentTemplate `json:"body"`
	}
	if err := json.Unmarshal(body, &respData); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	if respData.Status != "success" {
		return nil, fmt.Errorf("API response error: %s", respData.Status)
	}

	return &respData.Body, nil
}

func (c *Client) GetDeploymentTemplates(projectID string, page, number int) ([]DeploymentTemplate, error) {
	url := fmt.Sprintf("https://controller.thetaedgecloud.com/deployment_template/list_custom_templates?project_id=%s&page=%d&number=%d", projectID, page, number)

	body, err := sendRequest(c, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	var respData struct {
		Status string `json:"status"`
		Body   struct {
			TotalCount string               `json:"total_count"`
			Templates  []DeploymentTemplate `json:"templates"`
			Page       int                  `json:"page"`
			Number     int                  `json:"number"`
		} `json:"body"`
	}
	if err := json.Unmarshal(body, &respData); err != nil {
		return nil, err
	}

	if respData.Status != "success" {
		return nil, fmt.Errorf("API response error: %s", respData.Status)
	}

	return respData.Body.Templates, nil
}

func (c *Client) GetDeploymentTemplateByID(projectID, templateID string) (*DeploymentTemplate, error) {
	templates, err := c.GetDeploymentTemplates(projectID, 0, 100)
	if err != nil {
		return nil, err
	}

	for _, template := range templates {
		if template.ID == templateID {
			return &template, nil
		}
	}

	return nil, fmt.Errorf("template with ID %s not found", templateID)
}

func (c *Client) DeleteDeploymentTemplate(templateID, projectID string) (bool, error) {
	url := fmt.Sprintf("https://controller.thetaedgecloud.com/deployment_template/%s?project_id=%s", templateID, projectID)

	body, err := sendRequest(c, "DELETE", url, nil)
	if err != nil {
		return false, err
	}

	var respData DeleteDeploymentTemplateResponse
	if err := json.Unmarshal(body, &respData); err != nil {
		return false, err
	}

	if respData.Status != "success" {
		return false, fmt.Errorf("API response error: %s", respData.Status)
	}

	return respData.Body, nil
}
