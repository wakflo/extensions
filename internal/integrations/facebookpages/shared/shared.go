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
	"errors"
	"fmt"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
	"io"
	"net/http"

	"github.com/wakflo/go-sdk/autoform"
)

var (
	// #nosec
	tokenURL                = "https://graph.facebook.com/oauth/access_token"
	authURL                 = "https://graph.facebook.com/oauth/authorize"
	FacebookPagesSharedAuth = autoform.NewOAuthField(authURL, &tokenURL, []string{
		"pages_show_list pages_manage_posts business_management pages_read_engagement publish_video pages_manage_engagement pages_read_user_engagement",
	}).
		Build()
)

const baseURL = "https://graph.facebook.com/v22.0"

func MakeFacebookRequest(method, accessToken, url string, body map[string]interface{}) (map[string]interface{}, error) {
	fullURL := baseURL + url
	jsonData, err := json.Marshal(body)
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
		return nil, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(responseBody))
	}

	var result map[string]interface{}
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return result, nil
}

func GetFacebookPageInput(title string, desc string, required bool) *sdkcore.AutoFormSchema {
	getPages := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		reqURL := baseURL + "/me/accounts"
		req, err := http.NewRequest(http.MethodGet, reqURL, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating request: %v", err)
		}

		query := req.URL.Query()
		query.Add("access_token", ctx.Auth.AccessToken)
		query.Add("fields", "id,name")
		req.URL.RawQuery = query.Encode()

		req.Header.Add("Authorization", "Bearer "+ctx.Auth.AccessToken)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error sending request: %v", err)
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %v", err)
		}

		if res.StatusCode != http.StatusOK {
			var apiError map[string]interface{}
			if errs := json.Unmarshal(body, &apiError); errs == nil {
				return nil, fmt.Errorf("API error: %v", apiError["error"])
			}
			return nil, fmt.Errorf("API request failed with status code %d: %s", res.StatusCode, string(body))
		}

		var respData PagesResponse
		err = json.Unmarshal(body, &respData)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling response: %v", err)
		}

		var pages []map[string]string
		for _, page := range respData.Data {
			pages = append(pages, map[string]string{
				"id":   page.ID,
				"name": page.Name,
			})
		}

		return ctx.Respond(pages, len(pages))
	}

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName(title).
		SetDescription(desc).
		SetDynamicOptions(&getPages).
		SetRequired(required).
		Build()
}

func GetPagePostsInput(title string, desc string, required bool) *sdkcore.AutoFormSchema {
	getPages := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		input := sdk.DynamicInputToType[struct {
			PageID string `json:"page_id"`
		}](ctx)

		if input.PageID == "" {
			return nil, errors.New("please select a page")
		}

		// Fetch the page access token
		pageAccessToken, err := GetPageAccessToken(ctx.Auth.AccessToken, input.PageID)
		if err != nil {
			return nil, fmt.Errorf("error fetching page access token: %v", err)
		}

		// Use the page access token to fetch posts
		reqURL := fmt.Sprintf("%s/%s/feed", baseURL, input.PageID)
		req, err := http.NewRequest(http.MethodGet, reqURL, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating request: %v", err)
		}

		query := req.URL.Query()
		query.Add("access_token", pageAccessToken) // Use the page access token here
		query.Add("fields", "id,message,created_time")
		req.URL.RawQuery = query.Encode()

		req.Header.Add("Authorization", "Bearer "+pageAccessToken) // Use the page access token here

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error sending request: %v", err)
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %v", err)
		}

		if res.StatusCode != http.StatusOK {
			var apiError map[string]interface{}
			if errs := json.Unmarshal(body, &apiError); errs == nil {
				return nil, fmt.Errorf("API error: %v", apiError["error"])
			}
			return nil, fmt.Errorf("API request failed with status code %d: %s", res.StatusCode, string(body))
		}

		var respData PostsResponse
		err = json.Unmarshal(body, &respData)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling response: %v", err)
		}

		var posts []map[string]string
		for _, post := range respData.Data {
			posts = append(posts, map[string]string{
				"id":   post.ID,
				"name": post.Message,
			})
		}

		return ctx.Respond(posts, len(posts))
	}

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName(title).
		SetDescription(desc).
		SetDynamicOptions(&getPages).
		SetRequired(required).
		Build()
}

func GetPageAccessToken(userAccessToken, pageID string) (string, error) {
	endpoint := "/me/accounts"
	result, err := MakeFacebookRequest(http.MethodGet, userAccessToken, endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("error fetching user pages: %v", err)
	}

	pages, ok := result["data"].([]interface{})
	if !ok {
		return "", fmt.Errorf("invalid response format: missing 'data' field")
	}

	for _, page := range pages {
		pageMap, ok := page.(map[string]interface{})
		if !ok {
			continue
		}
		if pageMap["id"] == pageID {
			token, ok := pageMap["access_token"].(string)
			if !ok {
				return "", fmt.Errorf("invalid access token format")
			}
			return token, nil
		}
	}

	return "", fmt.Errorf("page not found")
}

func ActionFunc(pageAccessToken, endpoint string, payload map[string]interface{}) ([]map[string]interface{}, error) {
	result, err := MakeFacebookRequest(http.MethodGet, pageAccessToken, endpoint, payload)
	if err != nil {
		return nil, fmt.Errorf("error fetching page posts: %v", err)
	}

	posts, ok := result["data"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format: missing 'data' field")
	}

	var postList []map[string]interface{}
	for _, post := range posts {
		postMap, ok := post.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid post format")
		}
		postList = append(postList, postMap)
	}

	return postList, nil
}

func PostActionFunc(pageAccessToken, method, endpoint string, payload map[string]interface{}) (map[string]interface{}, error) {
	result, err := MakeFacebookRequest(method, pageAccessToken, endpoint, payload)
	if err != nil {
		return nil, fmt.Errorf("error performing POST request: %v", err)
	}
	return result, nil
}
