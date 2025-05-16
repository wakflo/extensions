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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/asana/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type createTaskActionProps struct {
	Name      string  `json:"name"`
	Project   *string `json:"project_id"`
	Workspace *string `json:"workspace_id"`
}

type CreateTaskAction struct{}

// Metadata returns metadata about the action
func (c *CreateTaskAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_task",
		DisplayName:   "Create Task",
		Description:   "Create a new task",
		Type:          core.ActionTypeAction,
		Documentation: createTaskDocs,
		SampleOutput: map[string]any{
			"Result": "Task created Successfully",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (c *CreateTaskAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_task", "Create Task")

	form.TextField("name", "Task Name").
		Placeholder("Enter task name").
		Required(true).
		HelpText("The name task's Name.")

	shared.RegisterWorkspacesProps(form)
	shared.RegisterProjectsProps(form)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (c *CreateTaskAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (c *CreateTaskAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[createTaskActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	var projects []string
	if input.Project != nil {
		projects = append(projects, *input.Project)
	}

	taskData := map[string]interface{}{
		"data": map[string]interface{}{
			"name":      input.Name,
			"workspace": input.Workspace,
			"projects":  projects,
		},
	}

	taskJSON, err := json.Marshal(taskData)
	if err != nil {
		return nil, err
	}

	reqURL := shared.BaseAPI + "/tasks"

	req, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewBuffer(taskJSON))
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

	fmt.Println(string(body))

	return map[string]interface{}{
		"Result": "Task created Successfully",
	}, nil
}

func NewCreateTaskAction() sdk.Action {
	return &CreateTaskAction{}
}
