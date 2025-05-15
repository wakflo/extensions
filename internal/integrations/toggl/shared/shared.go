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
	"strconv"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"

	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

var (
	form = smartform.NewAuthForm("toggl-auth", "Toggl API Authentication", smartform.AuthStrategyCustom)

	_ = form.TextField("api-key", "Api Key(Required)").
		Required(true).
		HelpText("Your toggl api key ")

	TogglSharedAuth = form.Build()
)

const baseURL = "https://api.track.toggl.com/api"

func CreateProjects(apiKey, workspaceID, name string, activeValue bool) (interface{}, error) {
	projectData := map[string]interface{}{
		"name":   name,
		"active": activeValue,
	}
	jsonData, err := json.Marshal(projectData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %v", err)
	}

	url := fmt.Sprintf("%s/v9/workspaces/%s/projects", baseURL, workspaceID)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	req.SetBasicAuth(apiKey, "api_token")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	var response interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func GetProject(apiKey, workspaceID string, sinceDate int64) (interface{}, error) {
	lastUpdate := strconv.FormatInt(sinceDate, 10)
	url := fmt.Sprintf("%s/v9/workspaces/%s/projects?query=since=%s", baseURL, workspaceID, lastUpdate)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	req.SetBasicAuth(apiKey, "api_token")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	var response interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func RegisterWorkspacesProp(form *smartform.FormBuilder) *smartform.FieldBuilder {
	getWorkspaces := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		// Create a new HTTP request
		req, err := http.NewRequest(http.MethodGet, baseURL+"/v9/workspaces", nil)
		if err != nil {
			return nil, err
		}

		// Set headers
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		req.Header.Set("Accept", "application/json")
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}
		// Set auth
		apiKeyFromContext := authCtx.Extra["api-key"]
		req.SetBasicAuth(apiKeyFromContext, "api_token")

		// Create HTTP client and send request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		// Read response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("API returned non-200 status code: %d", resp.StatusCode)
		}

		var workspaces []WorkspaceUser
		err = json.Unmarshal(body, &workspaces)
		if err != nil {
			return nil, err
		}

		// Return the response in the format expected by the dynamic field context
		return ctx.Respond(workspaces, len(workspaces))
	}

	return form.SelectField("workspaces", "Workspaces").
		Placeholder("Enter a value for Parent Folder.").
		Required(false).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getWorkspaces)).
				WithSearchSupport().
				WithPagination(10).
				End().
				GetDynamicSource(),
		).
		HelpText("Select a workspace")
}
