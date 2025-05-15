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

	"github.com/juicycleff/smartform/v1"

	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

const BaseAPI = "https://api.todoist.com/rest/v2"

var (
	todoistForm = smartform.NewAuthForm("todoist-auth", "Todoist OAuth", smartform.AuthStrategyOAuth2)
	_           = todoistForm.
			OAuthField("oauth", "Todoist OAuth").
			AuthorizationURL("https://todoist.com/oauth/authorize").
			TokenURL("https://todoist.com/oauth/access_token").
			Scopes([]string{"data:read_write"}).
			Required(true).
			Build()
)

var SharedTodoistAuth = todoistForm.Build()

func RegisterProjectsProps(form *smartform.FormBuilder) {
	getProjects := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}
		url := BaseAPI + "/projects"

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Accept", "application/json")

		req.Header.Set("Authorization", "Bearer "+authCtx.Token.AccessToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("API returned non-200 status code: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var projects []Project
		err = json.Unmarshal(body, &projects)
		if err != nil {
			return nil, err
		}

		return ctx.Respond(projects, len(projects))
	}

	form.SelectField("project_id", "Projects").
		Placeholder("Select a project").
		Required(true).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getProjects)).
				WithSearchSupport().
				WithPagination(10).
				End().
				GetDynamicSource(),
		).
		HelpText("Select a project")
}

type getSectionsFilter struct {
	ProjectID *string `json:"project_id"`
	SectionID *string `json:"section_id"`
	Label     *string `json:"label"`
	Filter    *string `json:"filter"`
	Lang      *string `json:"lang"`
	IDs       []int   `json:"ids"`
}

func RegisterSectionsProps(form *smartform.FormBuilder) {
	getSections := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		input := sdk.DynamicInputToType[getSectionsFilter](ctx)
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}
		baseUrl := BaseAPI + "/sections"

		queryParams := url.Values{}

		if input.ProjectID != nil {
			queryParams.Add("project_id", *input.ProjectID)
		}

		if input.SectionID != nil {
			queryParams.Add("section_id", *input.SectionID)
		}

		fullUrl := baseUrl
		if len(queryParams) > 0 {
			fullUrl += "?" + queryParams.Encode()
		}

		req, err := http.NewRequest(http.MethodGet, fullUrl, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Accept", "application/json")

		req.Header.Set("Authorization", "Bearer "+authCtx.Token.AccessToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("API returned non-200 status code: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var projects []ProjectSection
		err = json.Unmarshal(body, &projects)
		if err != nil {
			return nil, err
		}

		return ctx.Respond(projects, len(projects))
	}

	form.SelectField("section_id", "Section").
		Placeholder("Select a section").
		Required(true).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getSections)).
				WithSearchSupport().
				WithPagination(10).
				End().
				GetDynamicSource(),
		).
		HelpText("A section for the task. It should be a Section ID under the same project")
}

func RegisterTasksProps(form *smartform.FormBuilder, formId string, title string, desc string, required bool) {
	getTasks := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		input := sdk.DynamicInputToType[getSectionsFilter](ctx)

		// Build the URL with query parameters
		baseUrl := BaseAPI + "/tasks"

		// Create query parameters
		queryParams := url.Values{}

		if input.ProjectID != nil {
			queryParams.Add("project_id", *input.ProjectID)
		}

		if input.SectionID != nil {
			queryParams.Add("section_id", *input.SectionID)
		}

		if input.Label != nil {
			queryParams.Add("label", *input.Label)
		}

		if input.Filter != nil {
			queryParams.Add("filter", *input.Filter)
		}

		if input.Lang != nil {
			queryParams.Add("lang", *input.Lang)
		}

		fullUrl := baseUrl
		if len(queryParams) > 0 {
			fullUrl += "?" + queryParams.Encode()
		}

		req, err := http.NewRequest(http.MethodGet, fullUrl, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Accept", "application/json")

		req.Header.Set("Authorization", "Bearer "+authCtx.Token.AccessToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("API returned non-200 status code: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var tasks []Task
		err = json.Unmarshal(body, &tasks)
		if err != nil {
			return nil, err
		}

		items := make([]map[string]interface{}, 0, len(tasks))
		for _, task := range tasks {
			items = append(items, map[string]interface{}{
				"id":   task.ID,
				"name": task.Content,
			})
		}

		return ctx.Respond(items, len(items))
	}

	form.SelectField(formId, title).
		Placeholder("Select a task").
		Required(required).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getTasks)).
				WithSearchSupport().
				WithPagination(10).
				End().
				GetDynamicSource(),
		).
		HelpText(desc)
}
