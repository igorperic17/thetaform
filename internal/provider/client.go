package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	baseURL           string
	baseControllerURL string
	authToken         string
	userID            string
	orgID             string
	httpClient        *http.Client
}

func NewClient(email, password string) *Client {
	client := &Client{
		baseURL:           "https://api.thetaedgecloud.com",
		baseControllerURL: "https://controller.thetaedgecloud.com",
		httpClient:        &http.Client{},
	}

	authToken, userID, orgID, err := client.authenticate(email, password)
	if err != nil {
		return nil
	}

	client.authToken = authToken
	client.userID = userID
	client.orgID = orgID
	return client
}

func (c *Client) authenticate(email, password string) (string, string, string, error) {
	url := "https://api.thetaedgecloud.com/user/login?expand=redirect_project_id.org_id"
	payload := map[string]string{"email": email, "password": password}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", "", "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", "", "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Platform", "web")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", "", fmt.Errorf("authentication failed: %s", resp.Status)
	}

	var respData struct {
		Status string `json:"status"`
		Body   struct {
			Users []struct {
				AuthToken string `json:"auth_token"`
				ID        string `json:"id"`
			} `json:"users"`
			Organizations []struct {
				ID string `json:"id"`
			} `json:"organizations"`
		} `json:"body"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return "", "", "", err
	}

	if len(respData.Body.Users) == 0 {
		return "", "", "", fmt.Errorf("authentication failed: no users found")
	}

	return respData.Body.Users[0].AuthToken, respData.Body.Users[0].ID, respData.Body.Organizations[0].ID, nil
}
