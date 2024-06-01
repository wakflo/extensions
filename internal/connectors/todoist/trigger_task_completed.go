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
	"errors"

	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type taskCompletedTriggerProps struct {
	ProjectId *string `json:"projectId"`
}
type TaskCompletedTrigger struct {
	options *sdk.TriggerInfo
}

// NewTaskCompletedTrigger creates a new TaskCompletedTrigger object.
func NewTaskCompletedTrigger() *TaskCompletedTrigger {
	return &TaskCompletedTrigger{
		options: &sdk.TriggerInfo{
			Name:        "Task Completed",
			Description: "Triggers when a new task is completed",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"projectId": getProjectsInput(),
			},
			SampleOutput: map[string]interface{}{
				"content":      "Buy Milk",
				"meta_data":    nil,
				"user_id":      "2671355",
				"task_id":      "2995104339",
				"note_count":   0,
				"project_id":   "2203306141",
				"section_id":   "7025",
				"completed_at": "2015-02-17T15:40:41.000000Z",
				"id":           "1899066186",
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (t TaskCompletedTrigger) Run(ctx *sdk.RunContext) (sdk.Json, error) {
	if ctx.Auth.Token == nil {
		return nil, errors.New("missing todoist auth token")
	}

	_ = sdk.InputToType[taskCompletedTriggerProps](ctx)
	return nil, nil
}

func (t TaskCompletedTrigger) Test(ctx *sdk.RunContext) (sdk.Json, error) {
	return t.Run(ctx)
}

func (t TaskCompletedTrigger) OnEnabled(ctx *sdk.RunContext) error {
	return nil
}

func (t TaskCompletedTrigger) OnDisabled(ctx *sdk.RunContext) error {
	return nil
}

func (t TaskCompletedTrigger) GetInfo() *sdk.TriggerInfo {
	return t.options
}
