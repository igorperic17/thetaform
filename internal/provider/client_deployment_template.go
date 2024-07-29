package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type DeploymentTemplate struct {
	ID             string            `json:"id"`
	Name           string            `json:"name"`
	Description    string            `json:"description"`
	Tags           []string          `json:"tags"`
	Category       string            `json:"category"`
	ProjectID      string            `json:"project_id"`
	ContainerImage string            `json:"container_image"`
	ContainerPort  int64             `json:"container_port"`
	ContainerArgs  string            `json:"container_args"`
	EnvVars        map[string]string `json:"env_vars"`
	RequireEnvVars bool              `json:"require_env_vars"`
	Rank           int64             `json:"rank"`
	IconURL        string            `json:"icon_url"`
	CreateTime     time.Time         `json:"create_time"`
}

type DeploymentTemplateRequest struct {
	Name           string            `json:"name"`
	ProjectID      string            `json:"project_id"`
	Description    string            `json:"description"`
	ContainerImage string            `json:"container_image"`
	ContainerPort  int64             `json:"container_port"`
	ContainerArgs  string            `json:"container_args"`
	EnvVars        map[string]string `json:"env_vars"`
	Tags           []string          `json:"tags"`
	IconURL        string            `json:"icon_url"`
}

type DeploymentTemplateResponse struct {
	ID             string            `json:"id"`
	Name           string            `json:"name"`
	ProjectID      string            `json:"project_id"`
	Description    string            `json:"description"`
	ContainerImage string            `json:"container_image"`
	ContainerPort  int64             `json:"container_port"`
	ContainerArgs  string            `json:"container_args"`
	EnvVars        map[string]string `json:"env_vars"`
	Tags           []string          `json:"tags"`
	IconURL        string            `json:"icon_url"`
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

func (c *Client) CreateDeploymentTemplate(template DeploymentTemplateRequest) (*DeploymentTemplateResponse, error) {
	url := "https://controller.thetaedgecloud.com/deployment_template"

	jsonData, err := json.Marshal(template)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	setCommonHeaders(req, c)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request error: %s", resp.Status)
	}

	var respData DeploymentTemplateResponse
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, err
	}

	return &respData, nil
}

func (c *Client) UpdateDeploymentTemplate(templateID string, template DeploymentTemplateRequest) (*DeploymentTemplateResponse, error) {
	url := fmt.Sprintf("https://controller.thetaedgecloud.com/deployment_template/%s", templateID)

	jsonData, err := json.Marshal(template)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	setCommonHeaders(req, c)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request error: %s", resp.Status)
	}

	var respData DeploymentTemplateResponse
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, err
	}

	return &respData, nil
}

func (c *Client) GetDeploymentTemplates(projectID string, page, number int) ([]DeploymentTemplate, error) {
	url := fmt.Sprintf("https://controller.thetaedgecloud.com/deployment_template/list_custom_templates?project_id=%s&page=%d&number=%d", projectID, page, number)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", c.authToken)
	req.Header.Set("X-Auth-Id", c.userID)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request error: %s", resp.Status)
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
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
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

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("X-Auth-Token", c.authToken)
	req.Header.Set("X-Auth-Id", c.userID)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("API request error: %s", resp.Status)
	}

	var respData DeleteDeploymentTemplateResponse
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return false, err
	}

	if respData.Status != "success" {
		return false, fmt.Errorf("API response error: %s", respData.Status)
	}

	return respData.Body, nil
}

func setCommonHeaders(req *http.Request, c *Client) {
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Accept-Language", "en-GB,en;q=0.8")
	req.Header.Set("Origin", "https://www.thetaedgecloud.com")
	req.Header.Set("Referer", "https://www.thetaedgecloud.com/")
	req.Header.Set("Sec-Ch-Ua", "\"Brave\";v=\"123\", \"Not:A-Brand\";v=\"8\", \"Chromium\";v=\"123\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Gpc", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
	req.Header.Set("X-Auth-Id", c.userID)
	req.Header.Set("X-Auth-Token", c.authToken)
	req.Header.Set("X-Platform", "web")
}
