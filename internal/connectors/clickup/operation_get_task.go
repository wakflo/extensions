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

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type getTaskOperationProps struct {
	TaskID string `json:"task-id"`
}

type GetTaskOperation struct {
	options *sdk.OperationInfo
}

func NewGetTaskOperation() *GetTaskOperation {
	return &GetTaskOperation{
		options: &sdk.OperationInfo{
			Name:        "Get Task",
			Description: "Get a particular task",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"task-id": autoform.NewShortTextField().
					SetDisplayName("Task ID").
					SetDescription("Task ID").
					SetRequired(true).
					Build(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *GetTaskOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing clickup auth token")
	}
	accessToken := ctx.Auth.AccessToken

	input := sdk.InputToType[getTaskOperationProps](ctx)

	url := "https://api.clickup.com/api/v2/task/" + input.TaskID

	tasks, err := getData(accessToken, url)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (c *GetTaskOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *GetTaskOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
