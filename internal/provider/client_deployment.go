package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Deployment struct {
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	ProjectID         string            `json:"project_id"`
	DeploymentImageID string            `json:"deployment_image_id"`
	ContainerImage    string            `json:"container_image"`
	MinReplicas       int               `json:"min_replicas"`
	MaxReplicas       int               `json:"max_replicas"`
	VMID              string            `json:"vm_id"`
	Annotations       map[string]string `json:"annotations"`
	EnvVars           map[string]string `json:"env_vars"`
	Suffix            string            `json:"suffix"`
	URL               string            `json:"url"`
}

func (c *Client) CreateDeployment(deployment *Deployment) (*Deployment, error) {
	url := fmt.Sprintf("%s/deployment", c.baseURL)
	body, err := json.Marshal(deployment)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", c.authToken)
	req.Header.Set("X-Auth-Id", c.userID)

	fmt.Printf("Request URL: %s\n", url)
	fmt.Printf("Request Body: %s\n", body)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("Response Status: %s\n", resp.Status)
	fmt.Printf("Response Body: %s\n", respBody)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request error: %s", resp.Status)
	}

	var respData struct {
		Status string `json:"status"`
		Body   string `json:"body"`
	}
	if err := json.Unmarshal(respBody, &respData); err != nil {
		return nil, err
	}

	suffix := strings.TrimPrefix(respData.Body, "Custom deployment initiated. Access it at: https://")
	suffix = strings.TrimSuffix(suffix, "\n")

	deployment.Suffix = suffix
	deployment.URL = "https://" + suffix

	return deployment, nil
}

func (c *Client) GetDeployment(suffix string) (*Deployment, error) {
	url := fmt.Sprintf("%s/deployment/1/%s?project_id=%s", c.baseURL, suffix, "prj_8qf89pmjgdqurbaqfpdu3u854s6p")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
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
		Status string     `json:"status"`
		Body   Deployment `json:"body"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, err
	}

	return &respData.Body, nil
}

func (c *Client) UpdateDeployment(suffix string, deployment *Deployment) (*Deployment, error) {
	url := fmt.Sprintf("%s/deployment/1/%s?project_id=%s", c.baseURL, suffix, "prj_8qf89pmjgdqurbaqfpdu3u854s6p")
	body, err := json.Marshal(deployment)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
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
		Status string     `json:"status"`
		Body   Deployment `json:"body"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, err
	}

	return &respData.Body, nil
}

func (c *Client) DeleteDeployment(suffix string) error {
	url := fmt.Sprintf("%s/deployment/1/%s?project_id=%s", c.baseURL, suffix, "prj_8qf89pmjgdqurbaqfpdu3u854s6p")

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Auth-Token", c.authToken)
	req.Header.Set("X-Auth-Id", c.userID)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request error: %s", resp.Status)
	}

	return nil
}
