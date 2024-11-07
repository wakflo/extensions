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

package zohoinventory

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gookit/goutil/arrutil"
	fastshot "github.com/opus-domini/fast-shot"

	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
)

var (
	// #nosec
	tokenURL   = "https://accounts.zoho.com/oauth/v2/token"
	authURL    = "https://accounts.zoho.com/oauth/v2/auth"
	sharedAuth = autoform.NewOAuthField(authURL, &tokenURL, []string{
		"ZohoInventory.FullAccess.all",
	}).
		Build()
)

const baseURL = "https://www.zohoapis.com/inventory"

func getZohoClient(accessToken, url string) (map[string]interface{}, error) {
	fullURL := baseURL + url
	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Zoho-oauthtoken "+accessToken)
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return result, nil
}

func getOrganizationsInput() *sdkcore.AutoFormSchema {
	getOrganizations := func(ctx *sdkcore.DynamicFieldContext) (interface{}, error) {
		client := fastshot.NewClient("https://www.zohoapis.com/inventory").
			Auth().BearerToken(ctx.Auth.AccessToken).
			Header().
			AddAccept("application/json").
			Build()

		rsp, err := client.GET("/v1/organizations").Send()
		if err != nil {
			return nil, err
		}

		if rsp.Status().IsError() {
			return nil, errors.New(rsp.Status().Text())
		}

		bytes, err := io.ReadAll(rsp.Raw().Body) //nolint:bodyclose
		if err != nil {
			return nil, err
		}

		var organizations Organizations
		err = json.Unmarshal(bytes, &organizations)
		if err != nil {
			return nil, err
		}

		organization := organizations.Organizations
		return arrutil.Map[Organization, map[string]any](organization, func(input Organization) (target map[string]any, find bool) {
			return map[string]any{
				"id":   input.OrganizationID,
				"name": input.Name,
			}, true
		}), nil
	}

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName("Organizations").
		SetDescription("Select organization").
		SetDependsOn([]string{"connection"}).
		SetDynamicOptions(&getOrganizations).
		SetRequired(true).Build()
}
