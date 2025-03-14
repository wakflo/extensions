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

	"github.com/wakflo/go-sdk/autoform"
)

var (
	// #nosec
	tokenURL       = "https://api.infusionsoft.com/token"
	authURL        = "https://accounts.infusionsoft.com/app/oauth/authorize"
	KeapSharedAuth = autoform.NewOAuthField(authURL, &tokenURL, []string{
		"full",
	}).
		Build()
)

const (
	baseURL     = "https://api.infusionsoft.com/crm/rest/v1"
	contentType = "application/json"
)

// MakeKeapRequest makes a request to the Keap API
func MakeKeapRequest(accessToken, method, endpoint string, payload interface{}) (map[string]interface{}, error) {
	url := baseURL + endpoint
	client := &http.Client{}

	var reqBody io.Reader
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", contentType)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if len(body) > 0 {
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, err
		}
	}

	return result, nil
}

type ContactData struct {
	ID           string                 `json:"id,omitempty"`
	GivenName    string                 `json:"given_name,omitempty"`
	FamilyName   string                 `json:"family_name,omitempty"`
	Email        string                 `json:"email,omitempty"`
	PhoneNumbers []PhoneNumber          `json:"phone_numbers,omitempty"`
	Addresses    []Address              `json:"addresses,omitempty"`
	CustomFields map[string]interface{} `json:"custom_fields,omitempty"`
}

type PhoneNumber struct {
	Type   string `json:"type,omitempty"`
	Number string `json:"number,omitempty"`
}

type Address struct {
	Type       string `json:"type,omitempty"`
	Line1      string `json:"line1,omitempty"`
	Line2      string `json:"line2,omitempty"`
	Locality   string `json:"locality,omitempty"`
	Region     string `json:"region,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
	Country    string `json:"country,omitempty"`
}
