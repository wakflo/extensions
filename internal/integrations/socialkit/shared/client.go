// shared/client.go
package shared

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	BaseURL = "https://api.socialkit.dev"
)

// Response represents a generic API response from SocialKit
type Response struct {
	Data interface{}
}

// GetSocialKitClient makes a GET request to the SocialKit API
func GetSocialKitClient(accessKey string, endpoint string, queryParams map[string]string) (map[string]interface{}, error) {
	// Create HTTP client
	client := &http.Client{
		Timeout: time.Second * 30,
	}

	// Build URL with query parameters
	url := fmt.Sprintf("%s%s", BaseURL, endpoint)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %v", err)
	}

	// Add query parameters
	q := req.URL.Query()
	q.Add("access_key", accessKey)

	for key, value := range queryParams {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	// Add headers
	req.Header.Add("Accept", "application/json")
	req.Header.Add("x-access-key", accessKey)

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error executing request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body: %v", err)
	}

	// Parse response
	var result map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("Error parsing response: %v", err)
	}

	return result, nil
}
