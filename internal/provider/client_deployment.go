package provider

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DeploymentCreateRequest struct {
	ID                types.String            `tfsdk:"id"`
	Name              types.String            `tfsdk:"name"`
	ProjectID         types.String            `tfsdk:"project_id"`
	DeploymentImageID types.String            `tfsdk:"deployment_image_id"`
	ContainerImage    types.String            `tfsdk:"container_image"`
	MinReplicas       types.Int64             `tfsdk:"min_replicas"`
	MaxReplicas       types.Int64             `tfsdk:"max_replicas"`
	VMID              types.String            `tfsdk:"vm_id"`
	Annotations       map[string]types.String `tfsdk:"annotations"`
	AuthUsername      types.String            `tfsdk:"auth_username"`
	AuthPassword      types.String            `tfsdk:"auth_password"`
	URL               types.String            `tfsdk:"deployment_url"`
}

type DeploymentCreateRequestNative struct {
	Name              string            `json:"name"`
	ProjectID         string            `json:"project_id"`
	DeploymentImageID string            `json:"deployment_image_id"`
	ContainerImage    string            `json:"container_image"`
	MinReplicas       int64             `json:"min_replicas"`
	MaxReplicas       int64             `json:"max_replicas"`
	VMID              string            `json:"vm_id"`
	Annotations       map[string]string `json:"annotations"` // Ensure correct format
	AuthUsername      string            `json:"auth_username"`
	AuthPassword      string            `json:"auth_password"`
	URL               string            `json:"deployment_url"`
}

// Deployment represents the structure of a deployment response.
type Deployment struct {
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	ProjectID         string            `json:"project_id"`
	DeploymentImageID string            `json:"deployment_image_id"`
	ContainerImage    string            `json:"container_image"`
	MinReplicas       int64             `json:"min_replicas"`
	MaxReplicas       int64             `json:"max_replicas"`
	VMID              string            `json:"vm_id"`
	Annotations       map[string]string `json:"annotations"`
	AuthUsername      string            `json:"auth_username"`
	AuthPassword      string            `json:"auth_password"`
	ContainerPort     int64             `json:"container_port"`
	URL               string            `json:"deployment_url"`
}

func (c *Client) CreateDeployment(req DeploymentCreateRequestNative) (*Deployment, error) {
	url := fmt.Sprintf("%s/deployment", c.baseControllerURL)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	respBody, err := sendRequest(c, "POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}

	var result struct {
		Status string `json:"status"`
		Body   string `json:"body"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if result.Status != "success" {
		return nil, fmt.Errorf("API returned an error: %s", result.Body)
	}

	// Extract URL from the body string
	deploymentURL := extractDeploymentURL(result.Body)

	// Extract ID from the URL
	deploymentID, err := extractDeploymentID(deploymentURL)
	if err != nil {
		return nil, fmt.Errorf("failed to extract deployment ID: %v", err)
	}

	// Return the deployment with the ID and URL set
	return &Deployment{
		ID:                deploymentID,
		Name:              req.Name,
		ProjectID:         req.ProjectID,
		DeploymentImageID: req.DeploymentImageID,
		ContainerImage:    req.ContainerImage,
		MinReplicas:       req.MinReplicas,
		MaxReplicas:       req.MaxReplicas,
		VMID:              req.VMID,
		Annotations:       req.Annotations,
		AuthUsername:      req.AuthUsername,
		AuthPassword:      req.AuthPassword,
		URL:               deploymentURL,
	}, nil
}

func extractDeploymentURL(body string) string {
	// Extract the URL from the response body string
	// Assuming the URL is at the end of the string
	// This is a simple example, adapt as needed
	matches := regexp.MustCompile(`https://[^\s]+`).FindString(body)
	return matches
}

func extractDeploymentID(url string) (string, error) {
	// URL format is {name}-{id}.{rest_of_theta_domain}
	parts := strings.Split(url, ".")
	if len(parts) < 2 {
		return "", fmt.Errorf("unexpected URL format")
	}

	// Extract the part before the first dot
	beforeDot := parts[0]

	// Extract the ID from the part before the dot
	parts = strings.Split(beforeDot, "-")
	if len(parts) < 2 {
		return "", fmt.Errorf("unexpected URL format for ID extraction")
	}

	// The ID is the part after the name
	return parts[len(parts)-1], nil
}

func (c *Client) GetDeploymentByID(id string, projectID string) (*Deployment, error) {
	// URL to list all deployments
	url := fmt.Sprintf("%s/deployments/list?project_id=%s", c.baseControllerURL, projectID)

	// Perform the HTTP request
	respBody, err := sendRequest(c, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}

	// Process the response body
	var rawResponse map[string]interface{}
	if err := json.Unmarshal(respBody, &rawResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	// Ensure the response contains a "body" field
	body, ok := rawResponse["body"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format: missing or invalid body field")
	}

	// Search for the deployment with the matching Suffix
	for _, item := range body {
		deploymentData, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected response format: invalid item format")
		}

		// Check if the "Suffix" matches the requested ID
		if suffix, ok := deploymentData["Suffix"].(string); ok && suffix == id {
			// Create and populate Deployment struct
			deployment := &Deployment{
				ID:                suffix,
				Name:              getStringValue(deploymentData, "Name"),
				ProjectID:         getStringValue(deploymentData, "ProjectID"),
				DeploymentImageID: getStringValue(deploymentData, "DeploymentImageID"),
				ContainerImage:    getStringValue(deploymentData, "ContainerImage"),
				MinReplicas:       getInt64Value(deploymentData, "MinReplicas"),
				MaxReplicas:       getInt64Value(deploymentData, "MaxReplicas"),
				VMID:              getStringValue(deploymentData, "VMID"),
				Annotations:       convertToStringMap(getMapValue(deploymentData, "Annotations")),
				AuthUsername:      getStringValue(deploymentData, "AuthUsername"),
				AuthPassword:      getStringValue(deploymentData, "AuthPassword"),
				URL:               getStringValue(deploymentData, "Endpoint"),
			}

			return deployment, nil
		}
	}

	// Deployment not found
	return nil, fmt.Errorf("Deployment not found")
}

// Utility function to safely get a string value from a map
func getStringValue(data map[string]interface{}, key string) string {
	if value, ok := data[key].(string); ok {
		return value
	}
	return ""
}

// Utility function to safely get an int64 value from a map
func getInt64Value(data map[string]interface{}, key string) int64 {
	if value, ok := data[key].(float64); ok {
		return int64(value)
	}
	return 0
}

// Utility function to safely get a map value from a map
func getMapValue(data map[string]interface{}, key string) map[string]interface{} {
	if value, ok := data[key].(map[string]interface{}); ok {
		return value
	}
	return nil
}

func convertToStringMap(input map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for key, value := range input {
		strValue, ok := value.(string)
		if !ok {
			// Handle cases where the value is not a string
			strValue = fmt.Sprintf("%v", value) // Convert non-string values to string
		}
		result[key] = strValue
	}
	return result
}

func (c *Client) UpdateDeployment(id string, projectID string, req DeploymentCreateRequestNative) (*Deployment, error) {
	url := fmt.Sprintf("%s/deployment/1/%s?project_id=%s", c.baseControllerURL, id, projectID)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	// Send the request using the utility function
	respBody, err := sendRequest(c, "PUT", url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}

	var result struct {
		Status string     `json:"status"`
		Body   Deployment `json:"body"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if result.Status != "success" {
		return nil, fmt.Errorf("API returned an error: %s", result.Status)
	}

	return &result.Body, nil
}
func (c *Client) DeleteDeployment(id string, projectID string) (bool, error) {
	url := fmt.Sprintf("%s/deployment/1/%s?project_id=%s", c.baseControllerURL, id, projectID)

	// Send the request using the utility function
	respBody, err := sendRequest(c, "DELETE", url, nil)
	if err != nil {
		return false, fmt.Errorf("failed to send request: %v", err)
	}

	// Check the response body to determine success
	if string(respBody) == "" {
		return true, nil
	}

	return false, fmt.Errorf("unexpected response body: %s", string(respBody))
}
