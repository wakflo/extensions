package shared

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	telegramAPIBaseURL = "https://api.telegram.org/bot"
)

// GetTelegramClient makes a request to the Telegram API
func GetTelegramClient(token string, endpoint string, params map[string]interface{}) (map[string]interface{}, error) {
	baseURL := fmt.Sprintf("%s%s/%s", telegramAPIBaseURL, token, endpoint)

	var reqBody io.Reader = nil
	method := "GET"

	if len(params) > 0 {
		method = "POST"
		jsonData, err := json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request parameters: %v", err)
		}
		reqBody = strings.NewReader(string(jsonData))
	}

	req, err := http.NewRequest(method, baseURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	if method == "POST" {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("telegram API returned status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var responseObj map[string]interface{}
	if err := json.Unmarshal(body, &responseObj); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	// Check if the Telegram API returned an error
	if success, ok := responseObj["ok"].(bool); !ok || !success {
		errorDesc := "unknown error"
		if desc, ok := responseObj["description"].(string); ok {
			errorDesc = desc
		}
		return nil, errors.New(errorDesc)
	}

	return responseObj, nil
}

// UploadTelegramFile uploads a file to Telegram
func UploadTelegramFile(token string, method string, fileParam string, fileURL string, params map[string]string) (map[string]interface{}, error) {
	baseURL := fmt.Sprintf("%s%s/%s", telegramAPIBaseURL, token, method)

	// If fileURL is a web URL, we'll use sendPhoto with the URL parameter
	if strings.HasPrefix(fileURL, "http://") || strings.HasPrefix(fileURL, "https://") {
		// Add file URL to params
		if params == nil {
			params = make(map[string]string)
		}
		params[fileParam] = fileURL

		// Convert params to URL values
		values := url.Values{}
		for k, v := range params {
			values.Add(k, v)
		}

		// Make request
		resp, err := http.PostForm(baseURL, values)
		if err != nil {
			return nil, fmt.Errorf("error executing request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("telegram API returned status code %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %v", err)
		}

		var responseObj map[string]interface{}
		if err := json.Unmarshal(body, &responseObj); err != nil {
			return nil, fmt.Errorf("error unmarshaling response: %v", err)
		}

		// Check if the Telegram API returned an error
		if success, ok := responseObj["ok"].(bool); !ok || !success {
			errorDesc := "unknown error"
			if desc, ok := responseObj["description"].(string); ok {
				errorDesc = desc
			}
			return nil, errors.New(errorDesc)
		}

		return responseObj, nil
	}

	// For local files or more complex scenarios, we'd need to implement multipart/form-data upload
	// But for simplicity, we'll limit to URL-based uploads in this implementation
	return nil, errors.New("only URL-based file uploads are supported in this implementation")
}

func SendPhotoViaJSON(token string, params map[string]interface{}) (map[string]interface{}, error) {
	// The Telegram API accepts photo URLs via JSON too
	return GetTelegramClient(token, "sendPhoto", params)
}
