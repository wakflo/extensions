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
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/asana/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type updateTaskActionProps struct {
	TaskID    string  `json:"task_id"`
	Name      *string `json:"name"`
	Notes     *string `json:"notes"`
	Completed *bool   `json:"completed"`
	DueOn     *string `json:"due_on"`
}

type UpdateTaskAction struct{}

// Metadata returns metadata about the action
func (u *UpdateTaskAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "update_task",
		DisplayName:   "Update Task",
		Description:   "Update an existing task in Asana",
		Type:          core.ActionTypeAction,
		Documentation: updateTaskDocs,
		SampleOutput: map[string]any{
			"data": map[string]any{
				"gid":         "12345",
				"name":        "Updated Task Name",
				"completed":   true,
				"description": "This is an updated task",
			},
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (u *UpdateTaskAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("update_task", "Update Task")

	shared.RegisterProjectsProps(form)
	shared.RegisterTasksProps(form)

	form.TextField("name", "Name").
		Placeholder("Enter a new name").
		Required(false).
		HelpText("New name for the task")

	form.TextareaField("notes", "Notes").
		Placeholder("Enter notes").
		Required(false).
		HelpText("Notes or description for the task")

	form.CheckboxField("completed", "Completed").
		Required(false).
		HelpText("Mark the task as completed")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (u *UpdateTaskAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (u *UpdateTaskAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[updateTaskActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	taskData := map[string]interface{}{
		"data": map[string]interface{}{},
	}

	data := taskData["data"].(map[string]interface{})

	if input.Name != nil {
		data["name"] = *input.Name
	}

	if input.Notes != nil {
		data["notes"] = *input.Notes
	}

	if input.Completed != nil {
		data["completed"] = *input.Completed
	}

	if len(data) == 0 {
		return nil, errors.New("at least one field must be provided to update the task")
	}

	taskJSON, err := json.Marshal(taskData)
	if err != nil {
		return nil, err
	}

	reqURL := shared.BaseAPI + "/tasks/" + input.TaskID

	req, err := http.NewRequest(http.MethodPut, reqURL, bytes.NewBuffer(taskJSON))
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
		return nil, fmt.Errorf("request failed with status code: %d", res.StatusCode)
	}

	var response interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func NewUpdateTaskAction() sdk.Action {
	return &UpdateTaskAction{}
}
