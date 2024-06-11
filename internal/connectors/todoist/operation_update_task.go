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
	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type updateTaskOperationProps struct {
	TaskID *string `json:"taskId"`
	UpdateTask
}

type UpdateTaskOperation struct {
	options *sdk.OperationInfo
}

func NewUpdateTaskOperation() *UpdateTaskOperation {
	return &UpdateTaskOperation{
		options: &sdk.OperationInfo{
			Name:        "Create Task",
			Description: "Create task",
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"taskId": getTasksInput("Task ID", "ID of the task to update", true),

				"content": autoform.NewMarkdownField().
					SetDisplayName("Content").
					SetDescription("The task's content. It may contain some markdown-formatted text and hyperlinks").
					SetRequired(false).Build(),

				"description": autoform.NewMarkdownField().
					SetDisplayName("Description").
					SetDescription("A description for the task. This value may contain some markdown-formatted text and hyperlinks.").
					SetRequired(false).Build(),

				"labels": autoform.NewArrayField().
					SetDisplayName("Labels").
					SetDescription("The task's labels (a list of names that may represent either personal or shared labels)").
					SetItems(
						autoform.NewShortTextField().
							SetDisplayName("Label").
							SetDescription("Label").
							SetRequired(true).
							Build(),
					).
					SetRequired(false).Build(),

				"priority": autoform.NewNumberField().
					SetDisplayName("Priority").
					SetDescription("Task priority from 1 (normal) to 4 (urgent).").
					SetRequired(false).Build(),

				"dueDate": autoform.NewDateTimeField().
					SetDisplayName("Due date").
					SetDescription("Specific date in YYYY-MM-DD format relative to user's timezone").
					SetRequired(false).Build(),

				"duration": autoform.NewNumberField().
					SetDisplayName("Duration").
					SetDescription("A positive (greater than zero) integer for the amount of duration_unit the task will take, or null to unset. If specified, you must define a duration_unit.").
					SetRequired(false).Build(),
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

func (c *UpdateTaskOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	input := sdk.InputToType[updateTaskOperationProps](ctx)

	qu := fastshot.NewClient(baseAPI).
		Auth().BearerToken(ctx.Auth.AccessToken).
		Header().
		AddAccept("application/json").
		Build().POST(fmt.Sprintf("/tasks/%v", input.TaskID)).Body().AsJSON(input.UpdateTask)

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

func (c *UpdateTaskOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *UpdateTaskOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
