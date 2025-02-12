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
	"errors"
	"net/http"

	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
	sdk2 "github.com/wakflo/go-sdk/sdk"
)

type deleteTaskOperationProps struct {
	TaskID string `json:"task-id"`
}

type DeleteTaskOperation struct {
	options *sdk.OperationInfo
}

func NewDeleteTaskOperation() *DeleteTaskOperation {
	return &DeleteTaskOperation{
		options: &sdk.OperationInfo{
			Name:        "Delete Task",
			Description: "Delete task in a ClickUp workspace and list",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"workspace-id": getWorkSpaceInput("Workspaces", "select a workspace", true),
				"space-id":     getSpacesInput("Spaces", "select a space", true),
				"folder-id":    getFoldersInput("Folders", "select a folder", true),
				"list-id":      getListsInput("Lists", "select a list", true),
				"task-id":      getTasksInput("Tasks", "select a task", true),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *DeleteTaskOperation) Run(ctx *sdk.RunContext) (sdk2.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing clickup auth token")
	}
	accessToken := ctx.Auth.AccessToken

	input := sdk.InputToType[deleteTaskOperationProps](ctx)

	reqURL := baseURL + "/v2/task/" + input.TaskID
	req, err := http.NewRequest(http.MethodDelete, reqURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", accessToken)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return map[string]interface{}{
		"task": "Task Deleted",
	}, nil
}

func (c *DeleteTaskOperation) Test(ctx *sdk.RunContext) (sdk2.JSON, error) {
	return c.Run(ctx)
}

func (c *DeleteTaskOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
