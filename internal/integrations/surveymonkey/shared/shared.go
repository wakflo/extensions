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

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

var (
	surveyMonkeyTokenURL = "https://api.surveymonkey.com/oauth/token"
	surveyMonkeyAuthURL  = "https://api.surveymonkey.com/oauth/authorize"

	// Setting up OAuthField for SurveyMonkey
	// SharedAuth = auth.NewOauth2Auth(
	// 	surveyMonkeyAuthURL,
	// 	&surveyMonkeyTokenURL,
	// 	[]string{
	// 		"contacts_write",
	// 		"responses_read",
	// 		"responses_read_detail",
	// 		"webhooks_read",
	// 		"webhooks_write",
	// 	},
	// ).SetExcludedQueryParams([]string{"access_type", "token_access_type", "prompt"}).Build()
)

var (
	form = smartform.NewAuthForm("survey-monkey-auth", "SurveyMonkey Oauth", smartform.AuthStrategyOAuth2)
	_    = form.OAuthField("oauth", "SurveyMonkey Oauth").
		AuthorizationURL(surveyMonkeyAuthURL).
		TokenURL(surveyMonkeyTokenURL).
		Scopes([]string{
			"contacts_write",
			"responses_read",
			"responses_read_detail",
			"webhooks_read",
			"webhooks_write",
		}).
		Build()
)

var SurveyMonkeySharedAuth = form.Build()

var baseURL = "https://api.surveymonkey.com/v3"

func CreateContactList(accessToken, listName string) (map[string]interface{}, error) {
	// Define the API URL for creating contact lists
	url := baseURL + "/contact_lists"

	// Set up the payload with the contact list name
	payload := map[string]interface{}{
		"name": listName,
	}

	// Marshal the payload into JSON format
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// Create a new HTTP POST request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonPayload))
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
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		return response, fmt.Errorf("error: %v", response)
	}

	return response, nil
}

func GetSurveysProp(id string, title string, desc string, required bool, form *smartform.FormBuilder) *smartform.FieldBuilder {
	getSurveys := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		// Define the SurveyMonkey API URL for listing surveys
		url := baseURL + "/surveys"

		// Create a new HTTP GET request
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		tokenSource := ctx.Auth().Token
		if tokenSource == nil {
			return nil, errors.New("missing authentication token")
		}
		token := tokenSource.AccessToken

		// Set the required headers
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
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
		var surveysResponse SurveyMonkeySurveysResponse
		if err := json.Unmarshal(body, &surveysResponse); err != nil {
			return nil, err
		}

		// Extract surveys from the response
		surveys := surveysResponse.Data

		// log surveys
		fmt.Println("surveys----------------", surveys)

		var options []map[string]interface{}
		for _, survey := range surveys {
			options = append(options, map[string]interface{}{
				"value": survey.ID,
				"label": survey.Title,
			})
		}

		return ctx.Respond(options, len(options))
	}

	return form.SelectField(id, title).
		Placeholder("Select survey").
		Required(required).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getSurveys)).
				WithSearchSupport().
				WithPagination(10).
				End().
				GetDynamicSource(),
		).
		HelpText(desc)
}
