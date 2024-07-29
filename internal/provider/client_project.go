package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Project struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	OrgID         string   `json:"org_id"`
	TvaID         string   `json:"tva_id"`
	GatewayID     *string  `json:"gateway_id"`
	CreateTime    string   `json:"create_time"`
	UserJoinTime  string   `json:"user_join_time"`
	UserIDs       []string `json:"user_ids"`
	UserRole      string   `json:"user_role"`
	TvaSecret     string   `json:"tva_secret"`
	GatewayKey    *string  `json:"gateway_key"`
	GatewaySecret *string  `json:"gateway_secret"`
	Disabled      bool     `json:"disabled"`
}

type User struct {
	ID              string  `json:"id"`
	FirstName       string  `json:"first_name"`
	LastName        *string `json:"last_name"`
	Language        string  `json:"language"`
	CreateTime      string  `json:"create_time"`
	UpdateTime      string  `json:"update_time"`
	Email           string  `json:"email"`
	EmailVerified   bool    `json:"email_verified"`
	AuthToken       string  `json:"auth_token"`
	Email2FAEnabled bool    `json:"email_2fa_enabled"`
	OTP2FAEnabled   bool    `json:"otp_2fa_enabled"`
	OTPVerified     bool    `json:"otp_verified"`
	OptOutEmails    bool    `json:"opt_out_emails"`
	OptOutTexts     bool    `json:"opt_out_texts"`
}

type ResponseBody struct {
	Projects []Project `json:"projects"`
	Users    []User    `json:"users"`
}

type APIResponse struct {
	Status string       `json:"status"`
	Body   ResponseBody `json:"body"`
}

func (c *Client) CreateProject(project *Project) (*Project, error) {
	url := fmt.Sprintf("%s/project", c.baseURL)
	body, err := json.Marshal(project)
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

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request error: %s", resp.Status)
	}

	var respData struct {
		Status string  `json:"status"`
		Body   Project `json:"body"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, err
	}

	return &respData.Body, nil
}

func (c *Client) GetProjects(orgID string) (*[]Project, error) {
	url := fmt.Sprintf("https://api.thetaedgecloud.com/user/%s/organization/%s/projects?expand=user_ids", c.userID, orgID)

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

	var respData APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, err
	}

	return &respData.Body.Projects, nil
}

func (c *Client) UpdateProject(id string, project *Project) (*Project, error) {
	url := fmt.Sprintf("https://api.thetaedgecloud.com/project/%s", id)
	body, err := json.Marshal(project)
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
		Status string  `json:"status"`
		Body   Project `json:"body"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, err
	}

	return &respData.Body, nil
}

func (c *Client) DeleteProject(id string) error {
	url := fmt.Sprintf("https://api.thetaedgecloud.com/project/%s", id)

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
