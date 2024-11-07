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

package toggl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	fastshot "github.com/opus-domini/fast-shot"

	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
)

var sharedAuth = autoform.NewCustomAuthField().
	SetFields(map[string]*sdkcore.AutoFormSchema{
		"api-key": autoform.NewShortTextField().
			SetDisplayName("API Key").
			SetDescription("API Key").
			SetRequired(true).
			Build(),
	}).
	Build()

const baseURL = "https://api.track.toggl.com/api"

func createProject(apiKey, workspaceID, name string, activeValue bool) (interface{}, error) {
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

func getProjects(apiKey, workspaceID string, sinceDate int64) (interface{}, error) {
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

func getWorkSpaceInput() *sdkcore.AutoFormSchema {
	getWorkspaces := func(ctx *sdkcore.DynamicFieldContext) (interface{}, error) {

		qu := fastshot.NewClient(baseURL).
			Auth().BasicAuth(ctx.Auth.Extra["api-key"], "api_token").
			Header().
			AddAccept("application/json").
			Build().GET("/v9/workspaces")

		rsp, err := qu.Send()
		if err != nil {
			return nil, err
		}

		if rsp.Status().IsError() {
			return nil, errors.New(rsp.Status().Text())
		}
		defer rsp.Raw().Body.Close()

		bytes, err := io.ReadAll(rsp.Raw().Body) //nolint:bodyclose
		if err != nil {
			return nil, err
		}

		var workspaces []WorkspaceUser
		err = json.Unmarshal(bytes, &workspaces)
		if err != nil {
			return nil, err
		}

		return workspaces, nil
	}

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName("Workspaces").
		SetDescription("Select a workspace").
		SetDynamicOptions(&getWorkspaces).
		SetRequired(true).Build()
}
