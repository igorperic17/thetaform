package provider

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/andybalholm/brotli"
)

// Utility function to handle compressed responses
func readCompressedResponse(resp *http.Response) ([]byte, error) {
	var body []byte
	var err error
	contentEncoding := resp.Header.Get("Content-Encoding")
	switch contentEncoding {
	case "gzip":
		gzr, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error creating gzip reader: %v", err)
		}
		defer gzr.Close()
		body, err = ioutil.ReadAll(gzr)
		if err != nil {
			return nil, fmt.Errorf("error reading gzipped response body: %v", err)
		}
	case "deflate":
		zr, err := zlib.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error creating deflate reader: %v", err)
		}
		defer zr.Close()
		body, err = ioutil.ReadAll(zr)
		if err != nil {
			return nil, fmt.Errorf("error reading deflated response body: %v", err)
		}
	case "br":
		br := brotli.NewReader(resp.Body)
		body, err = ioutil.ReadAll(br)
		if err != nil {
			return nil, fmt.Errorf("error reading Brotli response body: %v", err)
		}
	default:
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %v", err)
		}
	}
	return body, err
}

func sendRequest(c *Client, method, url string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	setCommonHeaders(req, c)

	fmt.Printf("DEBUG: Sending %s request to %s\n", method, url)
	fmt.Printf("DEBUG: Request body: %s\n", string(body))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	fmt.Printf("DEBUG: Response status: %s\n", resp.Status)
	for key, values := range resp.Header {
		for _, value := range values {
			fmt.Printf("DEBUG: Header: %s: %s\n", key, value)
		}
	}

	respBody, err := readCompressedResponse(resp)
	if err != nil {
		return nil, err
	}

	fmt.Printf("DEBUG: Response body: %s\n", string(respBody))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request error: %s. Response body: %s", resp.Status, string(respBody))
	}

	return respBody, nil
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
