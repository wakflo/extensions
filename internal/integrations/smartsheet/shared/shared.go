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
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/juicycleff/smartform/v1"
)

var (
	smartsheetForm = smartform.NewAuthForm("smartsheet-auth", "Smartsheet OAuth", smartform.AuthStrategyOAuth2)
	_              = smartsheetForm.
			OAuthField("oauth", "Smartsheet OAuth").
			AuthorizationURL("https://app.smartsheet.com/b/authorize").
			TokenURL("https://api.smartsheet.com/2.0/token").
			Scopes([]string{
			"ADMIN_SHEETS ADMIN_USERS ADMIN_WEBHOOKS SHARE_SHEETS WRITE_SHEETS ADMIN_WORKSPACES CREATE_SHEETS READ_CONTACTS DELETE_SHEETS READ_SHEETS READ_USERS",
		}).
		Build()
)

var SmartsheetSharedAuth = smartsheetForm.Build()

const baseURL = "https://api.smartsheet.com"

func GetSmartSheetClient(accessToken, url string) (map[string]interface{}, error) {
	fullURL := baseURL + url
	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

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
		return nil, errors.New("error unmarshalling response")
	}

	return result, nil
}
