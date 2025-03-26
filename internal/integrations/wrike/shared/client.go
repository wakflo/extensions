package shared

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	sdkcore "github.com/wakflo/go-sdk/core"
)

const (
	// WrikeAPIBaseURL is the base URL for the Wrike API
	WrikeAPIBaseURL = "https://www.wrike.com/api/v4"
)

// GetWrikeClient makes a request to the Wrike API with the given access token and endpoint
func GetWrikeClient(accessToken, endpoint string) (sdkcore.JSON, error) {
	return makeWrikeAPIRequest(http.MethodGet, accessToken, endpoint, nil)
}

// PostWrikeClient makes a POST request to the Wrike API with the given access token, endpoint, and data
func PostWrikeClient(accessToken, endpoint string, data map[string]interface{}) (sdkcore.JSON, error) {
	return makeWrikeAPIRequest(http.MethodPost, accessToken, endpoint, data)
}

// PutWrikeClient makes a PUT request to the Wrike API with the given access token, endpoint, and data
func PutWrikeClient(accessToken, endpoint string, data map[string]interface{}) (sdkcore.JSON, error) {
	return makeWrikeAPIRequest(http.MethodPut, accessToken, endpoint, data)
}

// makeWrikeAPIRequest makes a request to the Wrike API with the given method, access token, endpoint, and data
func makeWrikeAPIRequest(method, accessToken, endpoint string, data map[string]interface{}) (sdkcore.JSON, error) {
	client := &http.Client{}

	url := WrikeAPIBaseURL + endpoint
	var req *http.Request
	var err error

	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return nil, err
		}
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("wrike API error: %s - %s", resp.Status, string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, err
	}

	return result, nil
}
