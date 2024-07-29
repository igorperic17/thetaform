package provider

import (
	"compress/gzip"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Organization represents the organization structure
type Organization struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	LogoURL      string `json:"logo_url"`
	CreateTime   string `json:"create_time"`
	UserJoinTime string `json:"user_join_time"`
	UserRole     string `json:"user_role"`
	Disabled     bool   `json:"disabled"`
	Suspended    bool   `json:"suspended"`
	Email        string `json:"email"`
}

// GetOrganizations fetches the list of organizations for the authenticated user
func (c *Client) GetOrganizations() ([]Organization, error) {
	url := fmt.Sprintf("%s/user/%s/orgs", c.baseURL, c.userID)
	println(url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-auth-id", c.userID)
	req.Header.Set("x-auth-token", c.authToken)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "gzip, deflate")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var bodyBytes []byte
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		gz, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		defer gz.Close()
		bodyBytes, err = ioutil.ReadAll(gz)
		if err != nil {
			return nil, err
		}
	case "deflate":
		df, err := zlib.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		defer df.Close()
		bodyBytes, err = ioutil.ReadAll(df)
		if err != nil {
			return nil, err
		}
	default:
		bodyBytes, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
	}

	if resp.StatusCode != http.StatusOK {
		bodyString := string(bodyBytes)
		return nil, fmt.Errorf("failed to get organizations: %s, response: %s", resp.Status, bodyString)
	}

	var respData struct {
		Status string `json:"status"`
		Body   struct {
			Organizations []Organization `json:"organizations"`
		} `json:"body"`
	}

	if err := json.Unmarshal(bodyBytes, &respData); err != nil {
		return nil, fmt.Errorf("error decoding response: %s, response: %s", err, bodyBytes)
	}

	return respData.Body.Organizations, nil
}
