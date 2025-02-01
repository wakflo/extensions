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
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/wakflo/go-sdk/autoform"
)

var (
	// #nosec
	tokenURL            = "https://accounts.zoho.com/oauth/v2/token"
	authURL             = "https://accounts.zoho.com/oauth/v2/auth"
	ZohoSalesSharedAuth = autoform.NewOAuthField(authURL, &tokenURL, []string{
		"SalesIQ.chatdetails.UPDATE SalesIQ.chatdetails.READ SalesIQ.visitordetails.UPDATE SalesIQ.visitordetails.READ",
	}).
		Build()
)

const baseURL = "https://salesiq.zoho.com/api/v1"

func GetZohoClient(accessToken, url string) (map[string]interface{}, error) {
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
