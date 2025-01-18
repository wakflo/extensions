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
	"io"
	"time"

	fastshot "github.com/opus-domini/fast-shot"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type createTaskOperationProps struct {
	Content     string     `json:"content"`
	Description *string    `json:"description"`
	ProjectID   *string    `json:"project_id"`
	SectionID   *string    `json:"section_id"`
	ParentID    *string    `json:"parent_id"`
	Labels      []string   `json:"labels"`
	Order       *string    `json:"order"`
	Priority    *int       `json:"priority"`
	DueDate     *time.Time `json:"dueDate"`
}

type CreateTaskOperation struct {
	options *sdk.OperationInfo
}

func NewCreateTaskOperation() *CreateTaskOperation {
	return &CreateTaskOperation{
		options: &sdk.OperationInfo{
			Name:        "Create Task",
			Description: "Creates a new task",
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"content": autoform.NewMarkdownField().
					SetDisplayName("Content").
					SetDescription("The task's content. It may contain some markdown-formatted text and hyperlinks").
					SetRequired(true).Build(),

				"description": autoform.NewMarkdownField().
					SetDisplayName("Description").
					SetDescription("A description for the task. This value may contain some markdown-formatted text and hyperlinks.").
					SetRequired(false).Build(),

				"project_id": getProjectsInput(),
				"section_id": getSectionsInput(),
				"parent_id":  getTasksInput("Parent Task ID", "Parent task ID.", false),
				"order": autoform.NewNumberField().
					SetDisplayName("Order").
					SetDescription("Non-zero integer value used by clients to sort tasks under the same parent.").
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

func (c *CreateTaskOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	input := sdk.InputToType[createTaskOperationProps](ctx)

	qu := fastshot.NewClient(baseAPI).
		Auth().BearerToken(ctx.Auth.AccessToken).
		Header().
		AddAccept("application/json").
		Build().POST("/tasks").Body().AsJSON(input)

	rsp, err := qu.Send()
	if err != nil {
		return nil, err
	}

	if rsp.Status().IsError() {
		return nil, errors.New(rsp.Status().Text())
	}

	bytes, err := io.ReadAll(rsp.Raw().Body)
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

func (c *CreateTaskOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CreateTaskOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
