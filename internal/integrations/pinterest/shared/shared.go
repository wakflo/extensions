// Copyright 2022-present Wakflo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package shared

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

var (
	// #nosec
	tokenURL = "https://api.pinterest.com/v5/oauth/token"
	authURL  = "https://www.pinterest.com/oauth"
)

var (
	form = smartform.NewAuthForm("pinterest-auth", "Pinterest Oauth", smartform.AuthStrategyOAuth2)
	_    = form.OAuthField("oauth", "Pinterest Oauth").
		AuthorizationURL(authURL).
		TokenURL(tokenURL).
		Scopes([]string{
			"pins:read", "pins:write", "boards:read", "boards:write", "ads:read", "ads:write", "boards:read_secret", "pins:read_secret", "user_accounts:read",
		}).Build()
)

var SharedPinterestAuth = form.Build()

var baseURL = "https://api.pinterest.com/v5"

func GetAdAccountProps(form *smartform.FormBuilder) {
	listAdAccounts := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		// Pinterest API endpoint for fetching ad accounts
		endpoint := baseURL + "/ad_accounts"

		// Create a new HTTP GET request
		req, err := http.NewRequest(http.MethodGet, endpoint, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %v", err)
		}

		// Set required headers according to Pinterest API documentation
		req.Header.Set("Authorization", "Bearer "+authCtx.Token.AccessToken)
		req.Header.Set("Accept", "application/json")

		// Add query parameters
		q := req.URL.Query()
		q.Add("page_size", "100") // Maximum allowed by Pinterest API
		req.URL.RawQuery = q.Encode()

		// Create HTTP client and send request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("ERROR: Failed to send request: %v\n", err)
			return nil, fmt.Errorf("failed to send request: %v", err)
		}
		defer resp.Body.Close()

		// Check response status
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			fmt.Printf("ERROR Response Body: %s\n", string(body))
			fmt.Printf("===========================================\n")
			return nil, fmt.Errorf("Pinterest API error (status %d): %s", resp.StatusCode, string(body))
		}

		// Read response body
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("ERROR: Failed to read response body: %v\n", err)
			return nil, fmt.Errorf("failed to read response body: %v", err)
		}

		// Parse the response
		var result map[string]interface{}
		if err := json.Unmarshal(responseBody, &result); err != nil {
			fmt.Printf("ERROR: Failed to unmarshal response: %v\n", err)
			return nil, fmt.Errorf("failed to unmarshal response: %v", err)
		}

		// Pinterest API returns ad accounts in "items" array
		itemsRaw, ok := result["items"].([]interface{})
		if !ok {
			fmt.Printf("ERROR: Items field not found or not an array\n")
			return nil, fmt.Errorf("items field not found or not an array")
		}

		var options []map[string]interface{}
		for i, itemRaw := range itemsRaw {
			adAccount, ok := itemRaw.(map[string]interface{})
			if !ok {
				fmt.Printf("WARNING: Skipping invalid ad account at index %d\n", i)
				continue
			}

			// Extract ad account details
			id, _ := adAccount["id"].(string)
			name, _ := adAccount["name"].(string)

			// Include additional info if available
			currency, _ := adAccount["currency"].(string)
			country, _ := adAccount["country"].(string)

			displayName := name
			if currency != "" || country != "" {
				displayName = fmt.Sprintf("%s (%s/%s)", name, country, currency)
			}
			fmt.Printf("  Display Name: %s\n", displayName)

			options = append(options, map[string]interface{}{
				"id":   id,
				"name": displayName,
			})
		}

		return ctx.Respond(options, len(options))
	}

	form.SelectField("ad_account_id", "Ad Account").
		Placeholder("Select an ad account").
		Required(false).
		// VisibleWhenEquals("is_busness_account", "yes").
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&listAdAccounts)).
				WithSearchSupport().
				End().
				GetDynamicSource(),
		).
		HelpText("Select a Pinterest ad account to use for creating pins")
}

func RegisterBoardPinsProps(form *smartform.FormBuilder) *smartform.FieldBuilder {
	getBoardPins := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		// Get the board_id from the dynamic input
		input := sdk.DynamicInputToType[struct {
			BoardID string `json:"board_id"`
		}](ctx)

		// Validate board ID
		if input.BoardID == "" {
			return nil, fmt.Errorf("board ID is required")
		}

		// Pinterest API endpoint for fetching pins from a specific board
		endpoint := fmt.Sprintf("%s/boards/%s/pins", baseURL, input.BoardID)

		// Create a new HTTP GET request
		req, err := http.NewRequest(http.MethodGet, endpoint, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %v", err)
		}

		// Set required headers
		req.Header.Set("Authorization", "Bearer "+authCtx.Token.AccessToken)
		req.Header.Set("Accept", "application/json")

		// Add query parameters
		q := req.URL.Query()
		q.Add("page_size", "100")
		q.Add("creative_types", "REGULAR,VIDEO") // Include both regular and video pins
		req.URL.RawQuery = q.Encode()

		// Create HTTP client and send request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to send request: %v", err)
		}
		defer resp.Body.Close()

		// Check response status
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("Pinterest API error (status %d): %s", resp.StatusCode, string(body))
		}

		// Read response body
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %v", err)
		}

		// Parse the response
		var result map[string]interface{}
		if err := json.Unmarshal(responseBody, &result); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %v", err)
		}

		// Pinterest API returns pins in "items" array
		itemsRaw, ok := result["items"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("items field not found or not an array")
		}

		var options []map[string]interface{}
		for _, itemRaw := range itemsRaw {
			pin, ok := itemRaw.(map[string]interface{})
			if !ok {
				continue
			}

			// Extract pin details
			id, _ := pin["id"].(string)
			title, _ := pin["title"].(string)
			description, _ := pin["description"].(string)
			createdAt, _ := pin["created_at"].(string)

			// Extract media info if available
			var mediaType string
			if media, ok := pin["media"].(map[string]interface{}); ok {
				mediaType, _ = media["media_type"].(string)
			}

			// Create display name with media type indicator
			displayName := title
			if displayName == "" && description != "" {
				if len(description) > 40 {
					displayName = description[:37] + "..."
				} else {
					displayName = description
				}
			}
			if displayName == "" {
				displayName = fmt.Sprintf("Pin %s", id)
			}

			// Add media type indicator
			if mediaType == "video" {
				displayName = "🎥 " + displayName
			} else if mediaType == "image" {
				displayName = "🖼️ " + displayName
			}

			// Add creation date if available
			if createdAt != "" {
				// Parse and format date (Pinterest uses RFC3339 format)
				if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
					displayName = fmt.Sprintf("%s (%s)", displayName, t.Format("Jan 2, 2006"))
				}
			}

			options = append(options, map[string]interface{}{
				"id":   id,
				"name": displayName,
			})
		}

		return ctx.Respond(options, len(options))
	}

	return form.SelectField("pin_id", "Board Pin").
		Placeholder("Select a pin from the board").
		Required(false).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getBoardPins)).
				WithFieldReference("board_id", "board_id").
				WithSearchSupport().
				End().
				RefreshOn("board_id").
				GetDynamicSource(),
		).
		HelpText("Select a pin from the selected board")
}

func RegisterBoardsProps(form *smartform.FormBuilder) *smartform.FieldBuilder {
	getBoards := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		// Pinterest API endpoint for fetching user's boards

		endpoint := baseURL + "/boards"

		// Create a new HTTP GET request
		req, err := http.NewRequest(http.MethodGet, endpoint, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %v", err)
		}

		// Set required headers according to Pinterest API documentation
		req.Header.Set("Authorization", "Bearer "+authCtx.Token.AccessToken)
		req.Header.Set("Accept", "application/json")

		// Add query parameters
		q := req.URL.Query()
		q.Add("page_size", "100") // Maximum allowed by Pinterest API
		q.Add("privacy", "ALL")   // Include both public and secret boards
		req.URL.RawQuery = q.Encode()

		// Create HTTP client and send request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to send request: %v", err)
		}
		defer resp.Body.Close()

		// Check response status
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("Pinterest API error (status %d): %s", resp.StatusCode, string(body))
		}

		// Read response body
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %v", err)
		}

		// Parse the response
		var result struct {
			Items []struct {
				ID            string `json:"id"`
				Name          string `json:"name"`
				Description   string `json:"description,omitempty"`
				Privacy       string `json:"privacy,omitempty"`
				PinCount      int    `json:"pin_count,omitempty"`
				FollowerCount int    `json:"follower_count,omitempty"`
				CreatedAt     string `json:"created_at,omitempty"`
			} `json:"items"`
		}

		if err := json.Unmarshal(responseBody, &result); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %v", err)
		}

		// Transform boards into options format
		var options []map[string]interface{}
		for _, board := range result.Items {
			// Create display name with additional info
			displayName := board.Name

			// Add privacy indicator for secret boards
			if board.Privacy == "SECRET" {
				displayName = "🔒 " + displayName
			}

			// Add pin count if available
			if board.PinCount > 0 {
				displayName = fmt.Sprintf("%s (%d pins)", displayName, board.PinCount)
			}

			options = append(options, map[string]interface{}{
				"id":   board.ID,
				"name": displayName,
			})
		}

		// Sort boards alphabetically
		sort.Slice(options, func(i, j int) bool {
			nameI, _ := options[i]["name"].(string)
			nameJ, _ := options[j]["name"].(string)
			return strings.ToLower(nameI) < strings.ToLower(nameJ)
		})

		return ctx.Respond(options, len(options))
	}

	return form.SelectField("board_id", "Board").
		Placeholder("Select a board").
		Required(true).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getBoards)).
				WithSearchSupport().
				End().
				GetDynamicSource(),
		).
		HelpText("Select a Pinterest board to pin to")
}

// API Calls
func GetPin(accessToken string, pinID string) (map[string]interface{}, error) {
	// Define the API URL for getting a specific pin
	url := baseURL + "/pins/" + pinID

	// Create a new HTTP GET request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// Set required headers
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	// Initialize the HTTP client and execute the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Parse the response body
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	// Check for a successful response
	if res.StatusCode != http.StatusOK {
		return response, fmt.Errorf("error getting pin, status code: %d, response: %v", res.StatusCode, response)
	}

	return response, nil
}

func DeletePin(accessToken string, pinID string) (map[string]interface{}, error) {
	// Define the API URL for deleting a specific pin
	url := baseURL + "/pins/" + pinID

	// Create a new HTTP DELETE request
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	// Set required headers
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	// Initialize the HTTP client and execute the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Parse the response body
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	// Check for a successful response (204 No Content is expected for successful deletion)
	if res.StatusCode != http.StatusNoContent && res.StatusCode != http.StatusOK {
		return response, fmt.Errorf("error deleting pin, status code: %d, response: %v", res.StatusCode, response)
	}

	return response, nil
}

func SearchPins(accessToken string, query string, bookmark string, adAccountID string) (map[string]interface{}, error) {
	// Define the API URL for searching pins
	url := baseURL + "/search/pins"

	// Create a new HTTP GET request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// Add query parameters
	q := req.URL.Query()
	q.Add("query", query)

	// Add bookmark for pagination if provided
	if bookmark != "" {
		q.Add("bookmark", bookmark)
	}

	// Add ad_account_id if provided (for business accounts)
	if adAccountID != "" {
		q.Add("ad_account_id", adAccountID)
	}

	req.URL.RawQuery = q.Encode()

	// Set required headers
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	// Initialize the HTTP client and execute the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Log the raw response for debugging
	fmt.Printf("Response Status: %d\n", res.StatusCode)
	fmt.Printf("Raw Response: %s\n", string(body))

	// Parse the response body
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	// Check for a successful response
	if res.StatusCode != http.StatusOK {
		return response, fmt.Errorf("error searching pins, status code: %d, response: %v", res.StatusCode, response)
	}

	return response, nil
}

func UpdatePin(accessToken string, pinID string, updateData map[string]interface{}) (map[string]interface{}, error) {
	// Define the API URL for updating a specific pin
	url := baseURL + "/pins/" + pinID

	// Convert update data to JSON
	jsonData, err := json.Marshal(updateData)
	if err != nil {
		return nil, err
	}

	// Create a new HTTP PATCH request
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	// Set required headers
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	// Initialize the HTTP client and execute the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Parse the response body
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	// Check for a successful response
	if res.StatusCode != http.StatusOK {
		return response, fmt.Errorf("error updating pin, status code: %d, response: %v", res.StatusCode, response)
	}

	return response, nil
}

func GetPinAnalytics(accessToken string, pinID string, metricTypes []string, startDate, endDate string, appTypes, splitField, adAccountID string) (map[string]interface{}, error) {
	// Define the API URL for getting pin analytics
	baseURL := baseURL + "/pins/" + pinID + "/analytics"

	// Create query parameters
	params := url.Values{}

	// Required parameters
	if len(metricTypes) > 0 {
		params.Set("metric_types", strings.Join(metricTypes, ","))
	}

	// Optional parameters
	if startDate != "" {
		params.Set("start_date", startDate)
	}
	if endDate != "" {
		params.Set("end_date", endDate)
	}
	if appTypes != "" {
		params.Set("app_types", appTypes)
	}
	if splitField != "" {
		params.Set("split_field", splitField)
	}
	if adAccountID != "" {
		params.Set("ad_account_id", adAccountID)
	}

	// Build URL with query parameters
	fullURL := baseURL
	if len(params) > 0 {
		fullURL = baseURL + "?" + params.Encode()
	}

	// Create a new HTTP GET request
	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, err
	}

	// Set required headers
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	// Initialize the HTTP client and execute the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Parse the response body
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	// Check for a successful response
	if res.StatusCode != http.StatusOK {
		return response, fmt.Errorf("error getting pin analytics, status code: %d, response: %v", res.StatusCode, response)
	}

	return response, nil
}

func CreatePin(accessToken string, pinData map[string]interface{}) (map[string]interface{}, error) {
	// Define the API URL for creating a pin
	url := baseURL + "/pins"

	// Convert pin data to JSON
	jsonData, err := json.Marshal(pinData)
	if err != nil {
		return nil, err
	}

	// Create a new HTTP POST request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	// Set required headers
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	// Initialize the HTTP client and execute the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Parse the response body
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	// Check for a successful response (201 Created for new resources)
	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK {
		return response, fmt.Errorf("error creating pin, status code: %d, response: %v", res.StatusCode, response)
	}

	return response, nil
}
