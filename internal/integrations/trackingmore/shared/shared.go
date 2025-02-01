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

	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
)

var TrackingMoreSharedAuth = autoform.NewCustomAuthField().
	SetFields(map[string]*sdkcore.AutoFormSchema{
		"key": autoform.NewShortTextField().SetDisplayName("Tracking Api Key (Required)").
			SetDescription("API Key").
			SetRequired(true).
			Build(),
	}).
	Build()

const baseURL = "https://api.trackingmore.com"

func CreateATracking(endpoint, applicationKey string, payload map[string]interface{}) (map[string]interface{}, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshalling payload:", err)
		return nil, err
	}

	fullURL := baseURL + endpoint

	req, err := http.NewRequest(http.MethodPost, fullURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Tracking-Api-Key", applicationKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error making the request:", err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %w", err)
	}

	return result, nil
}

func CreateBatchTracking(apiKey string, data []map[string]interface{}) (map[string]interface{}, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %w", err)
	}

	fullURL := baseURL + "/v4/trackings/batch"
	req, err := http.NewRequest(http.MethodPost, fullURL, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Tracking-Api-Key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	var responseMap map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return responseMap, nil
}

func ListTracking(endpoint, applicationKey string, date string) (map[string]interface{}, error) {
	encodedDate := url.QueryEscape(date)
	fullURL := fmt.Sprintf("%s%s?updated_date_min=%s", baseURL, endpoint, encodedDate)

	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Tracking-Api-Key", applicationKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error making the request:", err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %w", err)
	}

	return result, nil
}
