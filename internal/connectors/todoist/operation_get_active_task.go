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
	"fmt"
	"io"

	fastshot "github.com/opus-domini/fast-shot"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type getActiveTaskOperationProps struct {
	TaskID string `json:"taskId"`
}

type GetActiveTaskOperation struct {
	options *sdk.OperationInfo
}

func NewGetActiveTaskOperation() *GetActiveTaskOperation {
	return &GetActiveTaskOperation{
		options: &sdk.OperationInfo{
			Name:        "Get Active Task",
			Description: "Returns a single active (non-completed) task by ID",
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"taskId": getTasksInput("Task ID", "ID of the active task you want to retrieve", true),
			},
			SampleOutput: map[string]interface{}{},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
			RequireAuth: true,
		},
	}
}

func (c *GetActiveTaskOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	input := sdk.InputToType[getActiveTaskOperationProps](ctx)

	qu := fastshot.NewClient(baseApi).
		Auth().BearerToken(ctx.Auth.AccessToken).
		Header().
		AddAccept("application/json").
		Build().GET(fmt.Sprintf("/tasks/%s", input.TaskID))

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

	var task Task
	err = json.Unmarshal(bytes, &task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (c *GetActiveTaskOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *GetActiveTaskOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
