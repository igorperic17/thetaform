package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const loginURL = "https://api.thetaedgecloud.com/user/login?expand=redirect_project_id.org_id"
const baseURL = "https://controller.thetaedgecloud.com"

type Client struct {
	httpClient        *http.Client
	email             string
	password          string
	authToken         string
	redirectProjectID string
	userID            string
	orgID             string
}

func NewClient(email, password string) *Client {
	client := &Client{
		httpClient: &http.Client{},
		email:      email,
		password:   password,
	}
	err := client.authenticate()
	if err != nil {
		fmt.Printf("Error authenticating: %s\n", err)
	}
	return client
}

func (c *Client) authenticate() error {
	payload := map[string]string{
		"email":    c.email,
		"password": c.password,
	}
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", loginURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("authentication failed: %s", body)
	}

	var authResponse struct {
		Status string `json:"status"`
		Body   struct {
			Users []struct {
				ID                string `json:"id"`
				AuthToken         string `json:"auth_token"`
				RedirectProjectID string `json:"redirect_project_id"`
			} `json:"users"`
			Projects []struct {
				ID    string `json:"id"`
				OrgID string `json:"org_id"`
			} `json:"projects"`
		} `json:"body"`
	}
	if err := json.Unmarshal(body, &authResponse); err != nil {
		return err
	}

	if authResponse.Status != "success" {
		return fmt.Errorf("authentication failed: %s", body)
	}

	c.authToken = authResponse.Body.Users[0].AuthToken
	c.redirectProjectID = authResponse.Body.Users[0].RedirectProjectID
	c.userID = authResponse.Body.Users[0].ID
	c.orgID = authResponse.Body.Projects[0].OrgID

	return nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.authToken))
	req.Header.Set("X-Auth-Id", c.userID)
	req.Header.Set("X-Auth-Token", c.authToken)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API request error: %s", body)
	}

	return body, nil
}

type Endpoint struct {
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
}

func (c *Client) CreateEndpoint(endpoint *Endpoint) (*Endpoint, error) {
	url := fmt.Sprintf("%s/deployment", baseURL)
	jsonBody, err := json.Marshal(endpoint)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var createdEndpoint Endpoint
	if err := json.Unmarshal(body, &createdEndpoint); err != nil {
		return nil, err
	}

	return &createdEndpoint, nil
}

func (c *Client) GetEndpoint(id string) (*Endpoint, error) {
	url := fmt.Sprintf("%s/deployment/%s", baseURL, id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var endpoint Endpoint
	if err := json.Unmarshal(body, &endpoint); err != nil {
		return nil, err
	}

	return &endpoint, nil
}

func (c *Client) UpdateEndpoint(id string, endpoint *Endpoint) (*Endpoint, error) {
	url := fmt.Sprintf("%s/deployment/%s", baseURL, id)
	jsonBody, err := json.Marshal(endpoint)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var updatedEndpoint Endpoint
	if err := json.Unmarshal(body, &updatedEndpoint); err != nil {
		return nil, err
	}

	return &updatedEndpoint, nil
}

func (c *Client) DeleteEndpoint(id string) error {
	url := fmt.Sprintf("%s/deployment/%s", baseURL, id)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	return err
}
