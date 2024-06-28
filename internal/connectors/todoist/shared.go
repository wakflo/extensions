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

package todoist

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/gookit/goutil/arrutil"
	fastshot "github.com/opus-domini/fast-shot"
	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

var (
	// #nosec
	tokenURL   = "https://todoist.com/oauth/access_token"
	sharedAuth = autoform.NewOAuthField("https://todoist.com/oauth/authorize", &tokenURL, []string{
		"data:read_write",
	}).SetRequired(true).Build()
)
var baseAPI = "https://api.todoist.com/rest/v2"

func getProjectsInput() *sdkcore.AutoFormSchema {
	getProjects := func(ctx *sdkcore.DynamicFieldContext) (interface{}, error) {
		client := fastshot.NewClient(baseAPI).
			Auth().BearerToken(ctx.Auth.AccessToken).
			Header().
			AddAccept("application/json").
			Build()

		rsp, err := client.GET("/projects").Send()
		if err != nil {
			return nil, err
		}

		if rsp.IsError() {
			return nil, errors.New(rsp.StatusText())
		}

		bytes, err := io.ReadAll(rsp.RawBody())
		if err != nil {
			return nil, err
		}

		var projects []Project
		err = json.Unmarshal(bytes, &projects)
		if err != nil {
			return nil, err
		}

		return projects, nil
	}

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName("Project").
		SetDescription("Task project ID. If not set, task is put to user's Inbox.").
		SetDynamicOptions(&getProjects).
		SetRequired(false).Build()
}

type getSectionsFilter struct {
	ProjectID *string `json:"project_id"`
	SectionID *string `json:"section_id"`
	Label     *string `json:"label"`
	Filter    *string `json:"filter"`
	Lang      *string `json:"lang"`
	IDs       []int   `json:"ids"`
}

func getSectionsInput() *sdkcore.AutoFormSchema {
	getProjects := func(ctx *sdkcore.DynamicFieldContext) (interface{}, error) {
		input := sdk.DynamicInputToType[getSectionsFilter](ctx)

		qu := fastshot.NewClient(baseAPI).
			Auth().BearerToken(ctx.Auth.AccessToken).
			Header().
			AddAccept("application/json").
			Build().GET("/sections")

		if input.ProjectID != nil {
			qu = qu.Query().AddParam("project_id", *input.ProjectID)
		}

		if input.SectionID != nil {
			qu = qu.Query().AddParam("section_id", *input.SectionID)
		}

		rsp, err := qu.Send()
		if err != nil {
			return nil, err
		}

		if rsp.IsError() {
			return nil, errors.New(rsp.StatusText())
		}

		bytes, err := io.ReadAll(rsp.RawBody())
		if err != nil {
			return nil, err
		}

		var projects []ProjectSection
		err = json.Unmarshal(bytes, &projects)
		if err != nil {
			return nil, err
		}

		return projects, nil
	}

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName("Section").
		SetDescription("A section for the task. It should be a Section ID under the same project").
		SetDynamicOptions(&getProjects).
		SetRequired(false).Build()
}

var viewStyleOptions = []*sdkcore.AutoFormSchema{
	{Const: "list", Title: "List"},
	{Const: "board", Title: "Board"},
}

func getTasksInput(title string, desc string, required bool) *sdkcore.AutoFormSchema {
	getProjects := func(ctx *sdkcore.DynamicFieldContext) (interface{}, error) {
		input := sdk.DynamicInputToType[getSectionsFilter](ctx)

		qu := fastshot.NewClient(baseAPI).
			Auth().BearerToken(ctx.Auth.AccessToken).
			Header().
			AddAccept("application/json").
			Build().GET("/tasks")

		if input.ProjectID != nil {
			qu = qu.Query().AddParam("project_id", *input.ProjectID)
		}

		if input.SectionID != nil {
			qu = qu.Query().AddParam("section_id", *input.SectionID)
		}

		if input.Label != nil {
			qu = qu.Query().AddParam("label", *input.Label)
		}

		if input.Filter != nil {
			qu = qu.Query().AddParam("filter", *input.Filter)
		}

		if input.Lang != nil {
			qu = qu.Query().AddParam("lang", *input.Lang)
		}

		rsp, err := qu.Send()
		if err != nil {
			return nil, err
		}

		if rsp.IsError() {
			return nil, errors.New(rsp.StatusText())
		}

		bytes, err := io.ReadAll(rsp.RawBody())
		if err != nil {
			return nil, err
		}

		var tasks []Task
		err = json.Unmarshal(bytes, &tasks)
		if err != nil {
			return nil, err
		}

		return arrutil.Map[Task, map[string]any](tasks, func(input Task) (target map[string]any, find bool) {
			return map[string]any{
				"id":   input.ID,
				"name": input.Content,
			}, true
		}), nil
	}

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName(title).
		SetDescription(desc).
		SetDynamicOptions(&getProjects).
		SetRequired(required).Build()
}
