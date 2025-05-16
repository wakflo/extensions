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

package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/asana/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type listTasksActionProps struct {
	Project   *string `json:"project_id"`
	Workspace *string `json:"workspace_id"`
	Assignee  *string `json:"assignee"`
	Completed *bool   `json:"completed"`
	Limit     *int    `json:"limit"`
}

type ListTasksAction struct{}

// Metadata returns metadata about the action
func (l *ListTasksAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "list_tasks",
		DisplayName:   "List All Tasks",
		Description:   "Retrieve a list of tasks from Asana",
		Type:          core.ActionTypeAction,
		Documentation: listTasksDocs,
		SampleOutput: map[string]any{
			"data": []map[string]any{
				{
					"gid":       "12345",
					"name":      "Example Task 1",
					"completed": false,
				},
				{
					"gid":       "67890",
					"name":      "Example Task 2",
					"completed": true,
				},
			},
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (l *ListTasksAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("list_tasks", "List All Tasks")

	shared.RegisterProjectsProps(form)

	form.NumberField("limit", "Limit").
		Placeholder("Enter a limit").
		Required(false).
		DefaultValue(50).
		HelpText("Maximum number of tasks to return")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (l *ListTasksAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (l *ListTasksAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[listTasksActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	if input.Project == nil {
		return nil, errors.New("you must specify a project")
	}

	baseURL := shared.BaseAPI + "/tasks"
	params := url.Values{}

	if input.Project != nil {
		params.Add("project", *input.Project)
	}

	if input.Limit != nil {
		params.Add("limit", strconv.Itoa(*input.Limit))
	}

	params.Add("opt_fields", "gid,name,completed,due_on,notes,assignee,assignee_status,created_at,modified_at,projects,workspace")

	reqURL := baseURL
	if len(params) > 0 {
		reqURL = baseURL + "?" + params.Encode()
	}

	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+authCtx.Token.AccessToken)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		fmt.Println("Error response from Asana:", string(body))
		return nil, fmt.Errorf("request failed with status code: %d, response: %s", res.StatusCode, string(body))
	}

	var response interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func NewListTasksAction() sdk.Action {
	return &ListTasksAction{}
}
