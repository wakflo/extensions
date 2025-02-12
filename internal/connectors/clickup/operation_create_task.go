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
	sdk2 "github.com/wakflo/go-sdk/sdk"
)

type createTaskOperationProps struct {
	ListID      string `json:"list-id"`
	WorkspaceID string `json:"workspace-id"`
	SpaceID     string `json:"space-id"`
	FolderID    string `json:"folder-id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	AssigneeID  string `json:"assignee-id"`
}

type CreateTaskOperation struct {
	options *sdk.OperationInfo
}

func NewCreateTaskOperation() *CreateTaskOperation {
	return &CreateTaskOperation{
		options: &sdk.OperationInfo{
			Name:        "Create Task",
			Description: "Create a new task in a specified list",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"workspace-id": getWorkSpaceInput("Workspaces", "select a workspace", true),
				"space-id":     getSpacesInput("Spaces", "select a space", true),
				"folder-id":    getFoldersInput("Folders", "select a folder", true),
				"list-id":      getListsInput("Lists", "select a list to create task in", true),
				"assignee-id":  getAssigneeInput("Assignees", "select a assignee", true),
				"name": autoform.NewShortTextField().
					SetDisplayName("Task Name").
					SetDescription("The name of the task").
					SetRequired(true).
					Build(),
				"description": autoform.NewLongTextField().
					SetDisplayName("Task Description").
					SetDescription("The description of the task").
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

func (c *CreateTaskOperation) Run(ctx *sdk.RunContext) (sdk2.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing clickup auth token")
	}
	accessToken := ctx.Auth.AccessToken

	input, errs := sdk.InputToTypeSafely[createTaskOperationProps](ctx)
	if errs != nil {
		return nil, errs
	}

	taskData := map[string]interface{}{
		"name":        input.Name,
		"description": input.Description,
	}

	if input.Priority != "" {
		priority, err := strconv.Atoi(input.Priority)
		if err != nil {
			return nil, err
		}
		taskData["priority"] = priority
	}

	if input.AssigneeID != "" {
		assigneeStrings := strings.Split(input.AssigneeID, ",")
		assignees := make([]int, len(assigneeStrings))

		for i, assigneeStr := range assigneeStrings {
			assignee, err := strconv.Atoi(strings.TrimSpace(assigneeStr))
			if err != nil {
				fmt.Println("Assignee conversion error:", err)
				return nil, err
			}
			assignees[i] = assignee
		}
		taskData["assignees"] = assignees
	}

	taskJSON, err := json.Marshal(taskData)
	if err != nil {
		return nil, err
	}

	reqURL := baseURL + "/v2/list/" + input.ListID + "/task"
	req, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewBuffer(taskJSON))
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

	var response sdk2.JSON
	if newErr := json.NewDecoder(res.Body).Decode(&response); newErr != nil {
		return nil, err
	}

	return response, nil
}

func (c *CreateTaskOperation) Test(ctx *sdk.RunContext) (sdk2.JSON, error) {
	return c.Run(ctx)
}

func (c *CreateTaskOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
