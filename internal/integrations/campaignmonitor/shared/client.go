package shared

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	sdkcore "github.com/wakflo/go-sdk/core"
)

const BaseURL = "https://api.createsend.com/api/v3.3"

func GetCampaignMonitorClient(apiKey, clientID, endpoint string, method string, body interface{}) (sdkcore.JSON, error) {
	if !strings.HasPrefix(endpoint, "/") {
		endpoint = "/" + endpoint
	}

	fullURL := BaseURL + endpoint

	var req *http.Request
	var err error

	if body != nil && (method == "POST" || method == "PUT") {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request body: %v", err)
		}
		req, err = http.NewRequest(method, fullURL, bytes.NewBuffer(jsonBody))
		if err != nil {
			return nil, fmt.Errorf("error creating request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, fullURL, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating request: %v", err)
		}
	}

	req.SetBasicAuth(apiKey, "x")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	errCode := 400

	if resp.StatusCode >= errCode {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var result sdkcore.JSON
	if len(respBody) > 0 {
		err = json.Unmarshal(respBody, &result)
		if err != nil {
			return map[string]interface{}{
				"rawResponse": string(respBody),
			}, nil
		}
	} else {
		result = map[string]interface{}{}
	}

	return result, nil
}

func BuildQueryString(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}

	values := url.Values{}
	for key, value := range params {
		values.Add(key, value)
	}

	return "?" + values.Encode()
}
