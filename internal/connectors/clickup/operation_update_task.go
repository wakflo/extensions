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

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type updateTaskOperationProps struct {
	TaskID      string `json:"task-id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
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
				"task-id": autoform.NewShortTextField().
					SetDisplayName("Task ID").
					SetDescription("The task Id to be updated, without the '#' symbol. e.g, if the ID is #8235ck99n, just enter 8235ck99n. ").
					SetRequired(true).
					Build(),
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

	input := sdk.InputToType[updateTaskOperationProps](ctx)

	priority, err := strconv.Atoi(input.Priority)
	if err != nil {
		return nil, err
	}

	updatedTaskData := map[string]interface{}{
		"name":        input.Name,
		"description": input.Description,
		"priority":    priority,
	}

	taskJSON, err := json.Marshal(updatedTaskData)
	if err != nil {
		return nil, err
	}

	reqURL := "https://api.clickup.com/api/v2/task/" + input.TaskID
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
