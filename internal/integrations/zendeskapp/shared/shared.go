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

	"github.com/juicycleff/smartform/v1"
)

var form = smartform.NewAuthForm("zendesk-auth", "Zendesk Auth", smartform.AuthStrategyCustom)

var _ = form.TextField("email", "Agent Email").
	Required(true).
	Placeholder("Enter your Zendesk email").
	HelpText("The email address you use to login to Zendesk")

var _ = form.TextField("api-token", "API Token").
	Required(true).
	Placeholder("Enter your API token").
	HelpText("The API token you can generate in Zendesk")

var _ = form.TextField("subdomain", "Organization").
	Required(true).
	Placeholder("e.g. wakflohelp").
	HelpText("The subdomain of your Zendesk instance")

var SharedAuth = form.Build()

func ZendeskRequest(method, fullURL, email, apiToken string, request []byte) (map[string]interface{}, error) {
	req, err := http.NewRequest(method, fullURL, bytes.NewBuffer(request))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(email+"/token", apiToken)
	req.Header.Set("Accept", "application/json")
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
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return result, nil
}
