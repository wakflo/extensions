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
	"net/url"
	"strconv"
	"time"

	"github.com/juicycleff/smartform/v1"
)

var (
	form = smartform.NewAuthForm("cin7-auth", "Cin7 API Authentication", smartform.AuthStrategyCustom)

	_ = form.TextField("account_id", "Account ID (Required)").
		Required(true).
		HelpText("Your Account ID")

	_ = form.TextField("key", "API Application Key (Required)").
		Required(true).
		HelpText("API Application Key")

	Cin7SharedAuth = form.Build()
)

const baseURL = "https://inventory.dearsystems.com"

func FetchData(endpoint, accountID, applicationKey string, queryParams map[string]interface{}) (map[string]interface{}, error) {
	params := url.Values{}
	for key, value := range queryParams {
		switch v := value.(type) {
		case string:
			params.Add(key, v)
		case int:
			params.Add(key, strconv.Itoa(v))
		case float64:
			params.Add(key, strconv.FormatFloat(v, 'f', -1, 64))
		case bool:
			params.Add(key, strconv.FormatBool(v))
		case time.Time:
			params.Add(key, v.Format(time.RFC3339))
		default:
			params.Add(key, fmt.Sprintf("%v", v))
		}
	}

	fullURL := fmt.Sprintf("%s%s?%s", baseURL, endpoint, params.Encode())

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Auth-Accountid", accountID)
	req.Header.Set("Api-Auth-Applicationkey", applicationKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %w", err)
	}

	return result, nil
}
