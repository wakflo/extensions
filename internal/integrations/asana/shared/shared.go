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
	"net/url"
	"strconv"

	"github.com/gookit/goutil/arrutil"
	fastshot "github.com/opus-domini/fast-shot"

	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

var (
	// #nosec
	tokenURL        = "https://app.asana.com/-/oauth_token"
	AsanaSharedAuth = autoform.NewOAuthField("https://app.asana.com/-/oauth_authorize", &tokenURL, []string{
		"default",
	}).SetRequired(true).Build()
)
var BaseAPI = "https://app.asana.com/api/1.0"

func GetAsanaClient(accessToken string, endpoint string, method string, payload interface{}) (map[string]interface{}, error) {
	client := &http.Client{}

	var req *http.Request
	var err error

	url := fmt.Sprintf("%s%s", BaseAPI, endpoint)

	if payload != nil && (method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch) {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, err
		}
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return nil, err
		}
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetWorkspacesInput() *sdkcore.AutoFormSchema {
	getWorkspaces := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		client := fastshot.NewClient(BaseAPI).
			Auth().BearerToken(ctx.Auth.AccessToken).
			Header().
			AddAccept("application/json").
			Build()

		rsp, err := client.GET("/workspaces").Send()
		if err != nil {
			return nil, err
		}

		if rsp.Status().IsError() {
			return nil, errors.New(rsp.Status().Text())
		}

		bytes, err := io.ReadAll(rsp.Raw().Body) //nolint:bodyclose
		if err != nil {
			return nil, err
		}

		var workspaces WorkspaceResponse
		err = json.Unmarshal(bytes, &workspaces)
		if err != nil {
			return nil, err
		}

		workspace := workspaces.Data
		items := arrutil.Map[Workspace, map[string]any](workspace, func(input Workspace) (target map[string]any, find bool) {
			return map[string]any{
				"id":   input.GID,
				"name": input.Name,
			}, true
		})
		return ctx.Respond(items, len(items))
	}

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName("Workspace").
		SetDescription("Task workspace ID.").
		SetDependsOn([]string{"connection"}).
		SetDynamicOptions(&getWorkspaces).
		SetRequired(false).Build()
}

func GetProjectsInput() *sdkcore.AutoFormSchema {
	getProjects := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		qu := fastshot.NewClient(BaseAPI).
			Auth().BearerToken(ctx.Auth.AccessToken).
			Header().
			AddAccept("application/json").
			Build().GET("/projects")

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

		var response ProjectResponse
		err = json.Unmarshal(bytes, &response)
		if err != nil {
			return nil, err
		}

		projects := response.Data

		items := arrutil.Map[Project, map[string]any](projects, func(input Project) (target map[string]any, find bool) {
			return map[string]any{
				"id":   input.GID,
				"name": input.Name,
			}, true
		})

		return ctx.Respond(items, len(items))
	}

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName("Project").
		SetDescription("A project to create the task under").
		SetDynamicOptions(&getProjects).
		SetRequired(false).Build()
}

// func GetTasksInput() *sdkcore.AutoFormSchema {
// 	getTasks := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
// 		input := sdk.DynamicInputToType[struct {
// 			ProjectID string `json:"project-id,omitempty"`
// 		}](ctx)
// 		query := fastshot.NewClient(BaseAPI).
// 			Auth().BearerToken(ctx.Auth.AccessToken).
// 			Header().
// 			AddAccept("application/json").
// 			Build().GET("/tasks")

// 		query.Query().AddParam("project", input.ProjectID)

// 		rsp, err := query.Send()
// 		if err != nil {
// 			return nil, err
// 		}

// 		if rsp.Status().IsError() {
// 			return nil, errors.New(rsp.Status().Text())
// 		}
// 		defer rsp.Raw().Body.Close()

// 		bytes, err := io.ReadAll(rsp.Raw().Body) //nolint:bodyclose
// 		if err != nil {
// 			return nil, err
// 		}

// 		type Task struct {
// 			GID  string `json:"gid"`
// 			Name string `json:"name"`
// 		}

// 		type TaskResponse struct {
// 			Data []Task `json:"data"`
// 		}

// 		var response TaskResponse
// 		err = json.Unmarshal(bytes, &response)
// 		if err != nil {
// 			return nil, err
// 		}

// 		tasks := response.Data

// 		items := arrutil.Map[Task, map[string]any](tasks, func(input Task) (target map[string]any, find bool) {
// 			return map[string]any{
// 				"id":   input.GID,
// 				"name": input.Name,
// 			}, true
// 		})

// 		return ctx.Respond(items, len(items))
// 	}
// 	return autoform.NewDynamicField(sdkcore.String).
// 		SetDisplayName("Task").
// 		SetDescription("Select a task to retrieve").
// 		SetDynamicOptions(&getTasks).
// 		SetRequired(true).Build()
// }

func GetTasksInput() *sdkcore.AutoFormSchema {
	getTasks := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		input := sdk.DynamicInputToType[struct {
			ProjectID string `json:"project_id,omitempty"`
		}](ctx)
		// Build URL with query parameters
		baseURL := BaseAPI + "/tasks"
		params := url.Values{}

		params.Add("project", input.ProjectID)

		// Add fields to include in the response
		params.Add("opt_fields", "gid,name")

		// Build the complete URL with query parameters
		reqURL := baseURL + "?" + params.Encode()

		// Create the request
		req, err := http.NewRequest(http.MethodGet, reqURL, nil)
		if err != nil {
			return nil, err
		}

		// Add headers
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Authorization", "Bearer "+ctx.Auth.AccessToken)

		// Send the request
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		// Read the response
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		// Check for errors
		if res.StatusCode < 200 || res.StatusCode >= 300 {
			return nil, errors.New("request failed with status: " + strconv.Itoa(res.StatusCode))
		}

		type Task struct {
			GID  string `json:"gid"`
			Name string `json:"name"`
		}

		type TaskResponse struct {
			Data []Task `json:"data"`
		}

		var response TaskResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}

		tasks := response.Data

		items := arrutil.Map[Task, map[string]any](tasks, func(input Task) (target map[string]any, find bool) {
			return map[string]any{
				"id":   input.GID,
				"name": input.Name,
			}, true
		})

		return ctx.Respond(items, len(items))
	}

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName("Task").
		SetDescription("Select a task to retrieve").
		SetDependsOn([]string{"project"}). // Add dependency on project field
		SetDynamicOptions(&getTasks).
		SetRequired(true).Build()
}
