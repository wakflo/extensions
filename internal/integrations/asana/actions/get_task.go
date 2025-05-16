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
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/asana/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type getTaskActionProps struct {
	TaskID    string `json:"task_id"`
	ProjectID string `json:"project_id"`
}

type GetTaskAction struct{}

// Metadata returns metadata about the action
func (g *GetTaskAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_task",
		DisplayName:   "Get Task",
		Description:   "Retrieve details for a specific task by ID",
		Type:          core.ActionTypeAction,
		Documentation: getTaskDocs,
		SampleOutput: map[string]any{
			"data": map[string]any{
				"gid":         "12345",
				"name":        "Example Task",
				"completed":   false,
				"description": "This is an example task",
			},
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (g *GetTaskAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_task", "Get Task")

	shared.RegisterProjectsProps(form)
	shared.RegisterTasksProps(form)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (g *GetTaskAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (g *GetTaskAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[getTaskActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	baseURL := shared.BaseAPI + "/tasks/" + input.TaskID
	params := url.Values{}

	params.Add("opt_fields", "gid,name,completed,due_on,notes,assignee,description,assignee_status,created_at,modified_at,projects,subtasks,parent,tags,workspace")

	reqURL := baseURL
	if len(params) > 0 {
		reqURL = baseURL + "?" + params.Encode()
	}

	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}

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
		return nil, fmt.Errorf("request failed with status code: %d", res.StatusCode)
	}

	var response interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func NewGetTaskAction() sdk.Action {
	return &GetTaskAction{}
}
