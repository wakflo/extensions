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

	sdkcore "github.com/wakflo/go-sdk/core"

	"github.com/wakflo/go-sdk/autoform"
)

var (
	// #nosec
	tokenURL            = "https://graph.facebook.com/oauth/access_token"
	authURL             = "https://graph.facebook.com/oauth/authorize"
	InstagramSharedAuth = autoform.NewOAuthField(authURL, &tokenURL, []string{
		"instagram_basic instagram_content_publish business_management pages_show_list",
	}).
		Build()
)

const baseURL = "https://graph.facebook.com/v22.0"

func MakeInstagramRequest(method, accessToken, url string, body map[string]interface{}) (map[string]interface{}, error) {
	fullURL := baseURL + url
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	req, err := http.NewRequest(method, fullURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		var apiError map[string]interface{}
		if err := json.Unmarshal(responseBody, &apiError); err == nil {
			return nil, fmt.Errorf("API error: %v", apiError["error"])
		}
		return nil, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(responseBody))
	}

	var result map[string]interface{}
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return result, nil
}

func GetInstagramAccountInput(title string, desc string, required bool) *sdkcore.AutoFormSchema {
	getAccounts := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		// First get Facebook pages
		reqURL := baseURL + "/me/accounts"
		req, err := http.NewRequest(http.MethodGet, reqURL, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating request: %v", err)
		}

		query := req.URL.Query()
		query.Add("access_token", ctx.Auth.AccessToken)
		query.Add("fields", "id,instagram_business_account{id,name,username}")
		req.URL.RawQuery = query.Encode()

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error sending request: %v", err)
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %v", err)
		}

		var response struct {
			Data []struct {
				ID                       string `json:"id"`
				InstagramBusinessAccount struct {
					ID       string `json:"id"`
					Name     string `json:"name"`
					Username string `json:"username"`
				} `json:"instagram_business_account"`
			} `json:"data"`
		}

		if err := json.Unmarshal(body, &response); err != nil {
			return nil, fmt.Errorf("error unmarshaling response: %v", err)
		}

		var accounts []map[string]string
		for _, page := range response.Data {
			if page.InstagramBusinessAccount.ID != "" {
				accounts = append(accounts, map[string]string{
					"id": page.InstagramBusinessAccount.ID,
					"name": fmt.Sprintf("%s (@%s)",
						page.InstagramBusinessAccount.Name,
						page.InstagramBusinessAccount.Username,
					),
				})
			}
		}

		return ctx.Respond(accounts, len(accounts))
	}

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName(title).
		SetDescription(desc).
		SetDynamicOptions(&getAccounts).
		SetRequired(required).
		Build()
}

// CreateReelContainer Creates a container for Instagram reel
func CreateReelContainer(accessToken, instagramAccountId string, videoUrl string, caption string) (string, error) {
	endpoint := fmt.Sprintf("/%s/media", instagramAccountId)

	body := map[string]interface{}{
		"media_type":    "REELS",
		"video_url":     videoUrl,
		"caption":       caption,
		"share_to_feed": "true",
	}

	result, err := MakeInstagramRequest(http.MethodPost, accessToken, endpoint, body)
	if err != nil {
		return "", fmt.Errorf("error creating reel container: %v", err)
	}

	containerId, ok := result["id"].(string)
	if !ok {
		return "", fmt.Errorf("invalid container ID format in response")
	}

	return containerId, nil
}

func PublishReel(accessToken, instagramAccountId string, containerId string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("/%s/media_publish", instagramAccountId)

	body := map[string]interface{}{
		"creation_id": containerId,
	}

	result, err := MakeInstagramRequest(http.MethodPost, accessToken, endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("error publishing reel: %v", err)
	}

	return result, nil
}
