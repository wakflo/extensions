// shared/client.go
package shared

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	BaseURL = "https://www.sendowl.com/api/v1"
)

// Response represents a generic API response that can be either an array or an object
type Response struct {
	Array   []interface{}
	Map     map[string]interface{}
	IsArray bool
}

// GetSendOwlClient makes a GET request to the SendOwl API
func GetSendOwlClient(apiKey string, apiSecret string, endpoint string) (*Response, error) {
	// Create HTTP client
	client := &http.Client{
		Timeout: time.Second * 30,
	}

	// Create request
	url := fmt.Sprintf("%s%s", BaseURL, endpoint)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %v", err)
	}

	// Add basic auth
	req.SetBasicAuth(apiKey, apiSecret)
	req.Header.Add("Accept", "application/json")

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

	// Determine if response is array or object
	response := &Response{}
	bodyStr := string(bodyBytes)
	trimmedBody := strings.TrimSpace(bodyStr)

	if len(trimmedBody) > 0 && trimmedBody[0] == '[' {
		// Response is an array
		response.IsArray = true
		if err := json.Unmarshal(bodyBytes, &response.Array); err != nil {
			return nil, fmt.Errorf("Error parsing array response: %v", err)
		}
	} else {
		// Response is an object
		response.IsArray = false
		if err := json.Unmarshal(bodyBytes, &response.Map); err != nil {
			return nil, fmt.Errorf("Error parsing object response: %v", err)
		}
	}

	fmt.Println("Response 1 from SendOwl-------------------------->:", response)

	return response, nil

}

// PostSendOwlClient makes a POST request to the SendOwl API
func PostSendOwlClient(apiKey string, apiSecret string, endpoint string, body map[string]interface{}) (*Response, error) {
	// Convert body to JSON
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("Error creating request body: %v", err)
	}

	// Create HTTP client
	client := &http.Client{}

	// Create request
	url := fmt.Sprintf("%s%s", BaseURL, endpoint)
	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %v", err)
	}

	// Add basic auth
	req.SetBasicAuth(apiKey, apiSecret)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

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

	// Determine if response is array or object
	response := &Response{}
	bodyStr := string(bodyBytes)
	trimmedBody := strings.TrimSpace(bodyStr)

	if len(trimmedBody) > 0 && trimmedBody[0] == '[' {
		// Response is an array
		response.IsArray = true
		if err := json.Unmarshal(bodyBytes, &response.Array); err != nil {
			return nil, fmt.Errorf("Error parsing array response: %v", err)
		}
	} else {
		// Response is an object
		response.IsArray = false
		if err := json.Unmarshal(bodyBytes, &response.Map); err != nil {
			return nil, fmt.Errorf("Error parsing object response: %v", err)
		}
	}

	return response, nil
}
