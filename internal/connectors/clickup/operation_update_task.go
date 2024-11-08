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

package clickup

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type updateTaskOperationProps struct {
	TaskID      string `json:"task-id"`
	ListID      string `json:"list-id"`
	WorkspaceID string `json:"workspace-id"`
	SpaceID     string `json:"space-id"`
	FolderID    string `json:"folder-id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	AssigneeID  string `json:"assignee-id"`
}

type UpdateTaskOperation struct {
	options *sdk.OperationInfo
}

func NewUpdateTaskOperation() *UpdateTaskOperation {
	return &UpdateTaskOperation{
		options: &sdk.OperationInfo{
			Name:        "Update Task",
			Description: "Update task in a ClickUp workspace and list",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"workspace-id": getWorkSpaceInput("Workspaces", "select a workspace", true),
				"space-id":     getSpacesInput("Spaces", "select a space", true),
				"folder-id":    getFoldersInput("Folders", "select a folder", true),
				"list-id":      getListsInput("Lists", "select a list to create task in", true),
				"task-id":      getTasksInput("Tasks", "select a task to update", true),
				"assignee-id":  getAssigneeInput("Assignees", "select a assignee", true),
				"name": autoform.NewShortTextField().
					SetDisplayName("Task Name").
					SetDescription("The name of task to update").
					Build(),
				"description": autoform.NewLongTextField().
					SetDisplayName("Task Description").
					SetDescription("The description of task to update").
					Build(),
				"priority": autoform.NewSelectField().
					SetDisplayName("Priority").
					SetDescription("The priority level of the task").
					SetOptions(clickupPriorityType).
					Build(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *UpdateTaskOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing clickup auth token")
	}
	accessToken := ctx.Auth.AccessToken

	input, errs := sdk.InputToTypeSafely[updateTaskOperationProps](ctx)
	if errs != nil {
		return nil, errs
	}
	updatedTaskData := map[string]interface{}{}

	if input.Name != "" {
		updatedTaskData["name"] = input.Name
	}
	if input.Description != "" {
		updatedTaskData["description"] = input.Name
	}

	if input.Priority != "" {
		priority, err := strconv.Atoi(input.Priority)
		if err != nil {
			return nil, err
		}
		updatedTaskData["priority"] = priority
	}

	if input.AssigneeID != "" {
		assigneeStrings := strings.Split(input.AssigneeID, ",")
		assignees := make([]int, len(assigneeStrings))

		for i, assigneeStr := range assigneeStrings {
			assignee, err := strconv.Atoi(strings.TrimSpace(assigneeStr))
			if err != nil {
				return nil, err
			}
			assignees[i] = assignee
		}
		assigneesObject := map[string][]int{
			"add": assignees,
			"rem": {},
		}

		updatedTaskData["assignees"] = assigneesObject
	}

	taskJSON, err := json.Marshal(updatedTaskData)
	if err != nil {
		return nil, err
	}

	reqURL := baseURL + "/v2/task/" + input.TaskID
	req, err := http.NewRequest(http.MethodPut, reqURL, bytes.NewBuffer(taskJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", accessToken)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var response sdk.JSON
	fmt.Println(response)
	return map[string]interface{}{
		"Report": "Task updated successfully",
	}, nil
}

func (c *UpdateTaskOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *UpdateTaskOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
