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
	"github.com/juicycleff/smartform/v1"
	fastshot "github.com/opus-domini/fast-shot"

	sdkcontext "github.com/wakflo/go-sdk/v2/context"

	"github.com/wakflo/go-sdk/v2"

	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

const BaseAPI = "https://app.asana.com/api/1.0"

var (
	asanaForm = smartform.NewAuthForm("asana-auth", "Asana OAuth", smartform.AuthStrategyOAuth2)
	_         = asanaForm.
			OAuthField("oauth", "Asana OAuth").
			AuthorizationURL("https://app.asana.com/-/oauth_authorize").
			TokenURL("https://app.asana.com/-/oauth_token").
			Scopes([]string{
			"default",
		}).
		Required(true).
		Build()
)

var AsanaSharedAuth = asanaForm.Build()

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

func RegisterWorkspacesProps(form *smartform.FormBuilder) *smartform.FieldBuilder {
	getWorkspaces := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		client := fastshot.NewClient(BaseAPI).
			Auth().BearerToken(authCtx.AccessToken).
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

	return form.SelectField("workspace_id", "Workspace").
		Placeholder("Select a workspace").
		Required(false).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getWorkspaces)).
				WithFieldReference("connection", "connection").
				WithSearchSupport().
				End().
				// RefreshOn("connection").
				GetDynamicSource(),
		).
		HelpText("Task workspace ID.")
}

func RegisterProjectsProps(form *smartform.FormBuilder) *smartform.FieldBuilder {
	getProjects := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		qu := fastshot.NewClient(BaseAPI).
			Auth().BearerToken(authCtx.AccessToken).
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

	return form.SelectField("project_id", "Project").
		Placeholder("Select a project").
		Required(false).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getProjects)).
				WithSearchSupport().
				End().
				GetDynamicSource(),
		).
		HelpText("A project to create the task under")
}

func RegisterTasksProps(form *smartform.FormBuilder) *smartform.FieldBuilder {
	getTasks := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		input := sdk.DynamicInputToType[struct {
			ProjectID string `json:"project_id,omitempty"`
		}](ctx)

		baseURL := BaseAPI + "/tasks"
		params := url.Values{}

		params.Add("project", input.ProjectID)
		params.Add("opt_fields", "gid,name")

		reqURL := baseURL + "?" + params.Encode()

		req, err := http.NewRequest(http.MethodGet, reqURL, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Add("Accept", "application/json")
		req.Header.Add("Authorization", "Bearer "+authCtx.AccessToken)

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

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

	return form.SelectField("task_id", "Tasks").
		Placeholder("Select a task").
		Required(true).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getTasks)).
				WithFieldReference("project_id", "project_id").
				WithSearchSupport().
				End().
				RefreshOn("project_id").
				GetDynamicSource(),
		).
		HelpText("Select a task to retrieve")
}
