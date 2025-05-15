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

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"

	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

const baseURL = "https://api.typeform.com/"

var (
	typeformForm = smartform.NewAuthForm("typeform-auth", "Typeform OAuth", smartform.AuthStrategyOAuth2)
	_            = typeformForm.
			OAuthField("oauth", "Typeform OAuth").
			AuthorizationURL(baseURL + "oauth/authorize").
			TokenURL(baseURL + "oauth/token").
			Scopes([]string{"accounts:read forms:write forms:read responses:read responses:write"}).
			Build()
)

var SharedTypeformAuth = typeformForm.Build()

func RegisterTypeformFormsProps(form *smartform.FormBuilder, label string, hint string, required bool) {
	getTypeformForms := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		// Define the Typeform API URL for listing forms
		url := baseURL + "forms"

		// Create a new HTTP GET request
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		// Set the required headers
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *&ctx.Auth().Token.AccessToken))
		req.Header.Set("Content-Type", "application/json")

		// Create a new HTTP client and send the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		// Check if the response status indicates an error
		if resp.StatusCode >= 400 {
			return nil, fmt.Errorf("error: %s", resp.Status)
		}

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		// Parse the response body into the struct
		var typeformResponse FormsResponse
		if err := json.Unmarshal(body, &typeformResponse); err != nil {
			return nil, err
		}

		// Extract forms from the response
		forms := typeformResponse.Items

		var options []map[string]interface{}
		for _, form := range forms {
			options = append(options, map[string]interface{}{
				"id":   form.ID,
				"name": form.Title,
			})
		}

		return ctx.Respond(options, len(options))
	}

	form.SelectField(label, label).
		Placeholder("Enter a value.").
		Required(required).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getTypeformForms)).
				// WithFieldReference("state", "state").
				WithSearchSupport().
				WithPagination(10).
				End().
				// RefreshOn("state").
				GetDynamicSource(),
		).
		HelpText(hint)
}

func GetFormResponses(accessToken, formID string) (map[string]interface{}, error) {
	// Define the API URL for retrieving form responses
	url := fmt.Sprintf(baseURL+"forms/%s/responses", formID)

	// Create a new HTTP GET request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// Set required headers
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Content-Type", "application/json")

	// Initialize the HTTP client and execute the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Parse the response body
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	// Check for a successful response
	if res.StatusCode != http.StatusOK {
		return response, fmt.Errorf("error: %v", response)
	}

	return response, nil
}
