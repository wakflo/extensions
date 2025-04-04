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
	"io"
	"net/http"

	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
)

var (
	// #nosec
	tokenURL   = "https://accounts.zoho.com/oauth/v2/token"
	authURL    = "https://accounts.zoho.com/oauth/v2/auth"
	SharedAuth = autoform.NewOAuthField(authURL, &tokenURL, []string{
		"ZohoCRM.modules.ALL", "ZohoCRM.settings.ALL",
	}).Build()
)

const BaseURL = "https://www.zohoapis.com/crm/v7/"

func GetZohoCRMClient(accessToken, method, endpoint string, body interface{}) (map[string]interface{}, error) {
	fullURL := BaseURL + endpoint

	var req *http.Request
	var err error

	if body != nil {
		bodyJSON, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request body: %v", err)
		}

		req, err = http.NewRequest(method, fullURL, bytes.NewBuffer(bodyJSON))
		if err != nil {
			return nil, fmt.Errorf("error creating request: %v", err)
		}
	} else {
		req, err = http.NewRequest(method, fullURL, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating request: %v", err)
		}
	}

	req.Header.Set("Authorization", "Zoho-oauthtoken "+accessToken)
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

func GetModulesFunction() *sdkcore.DynamicOptionsFn {
	getModules := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		result, err := GetZohoCRMClient(ctx.Auth.AccessToken, http.MethodGet, "settings/modules", nil)
		if err != nil {
			return nil, fmt.Errorf("error fetching modules: %v", err)
		}

		modules, ok := result["modules"].([]interface{})
		if !ok {
			return nil, errors.New("invalid response format: modules field is missing or invalid")
		}

		var items []map[string]any
		for _, module := range modules {
			moduleData, ok := module.(map[string]interface{})
			if !ok {
				continue
			}

			apiName, ok := moduleData["api_name"].(string)
			if !ok {
				continue
			}

			displayName := apiName
			if pluralLabel, ok := moduleData["plural_label"].(string); ok {
				displayName = pluralLabel
			}

			items = append(items, map[string]any{
				"id":   apiName,
				"name": displayName,
			})
		}

		return ctx.Respond(items, len(items))
	}

	return &getModules
}
