package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const baseURL = "https://controller.thetaedgecloud.com"

type Client struct {
	httpClient *http.Client
	apiKey     string
	apiSecret  string
	authID     string
	authToken  string
}

func NewClient(apiKey, apiSecret string) *Client {
	client := &Client{
		httpClient: &http.Client{},
		apiKey:     apiKey,
		apiSecret:  apiSecret,
	}
	err := client.authenticate()
	if err != nil {
		fmt.Printf("Error authenticating: %s\n", err)
	}
	return client
}

func (c *Client) authenticate() error {
	url := fmt.Sprintf("%s/auth/login", baseURL)
	payload := map[string]string{
		"api_key":    c.apiKey,
		"api_secret": c.apiSecret,
	}
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
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
		AuthID    string `json:"auth_id"`
		AuthToken string `json:"auth_token"`
	}
	if err := json.Unmarshal(body, &authResponse); err != nil {
		return err
	}

	c.authID = authResponse.AuthID
	c.authToken = authResponse.AuthToken
	return nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.authToken))
	req.Header.Set("X-Auth-Id", c.authID)
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
