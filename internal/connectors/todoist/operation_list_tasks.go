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

	fastshot "github.com/opus-domini/fast-shot"
	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type listTasksOperationProps struct {
	ProjectId string `json:"projectId"`
}

type ListTasksOperation struct {
	options *sdk.OperationInfo
}

func NewListTasksOperation() *ListTasksOperation {
	return &ListTasksOperation{
		options: &sdk.OperationInfo{
			Name:        "List Tasks",
			Description: "Returns a list containing all active tasks",
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"project_id": getProjectsInput(),
				"section_id": getSectionsInput(),
				"label": autoform.NewShortTextField().
					SetDisplayName("Label").
					SetDescription("Filter tasks by label name.").
					SetRequired(false).Build(),
				"filter": autoform.NewShortTextField().
					SetDisplayName("Filter").
					SetDescription("Filter by any supported filter. Multiple filters (using the comma , operator) are not supported.").
					SetRequired(false).Build(),
				"lang": autoform.NewShortTextField().
					SetDisplayName("Lang").
					SetDescription("IETF language tag defining what language filter is written in, if differs from default English.").
					SetRequired(false).Build(),
				"ids": autoform.NewArrayField().
					SetDisplayName("IDs").
					SetDescription("A list of the task IDs to retrieve, this should be a comma separated list.").
					SetItems(
						autoform.NewShortTextField().
							SetDisplayName("ID").
							SetDescription("id").
							SetRequired(true).
							Build(),
					).
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

func (c *ListTasksOperation) Run(ctx *sdk.RunContext) (sdk.Json, error) {
	_ = sdk.InputToType[listTasksOperationProps](ctx)

	client := fastshot.NewClient(baseApi).
		Auth().BearerToken(ctx.Auth.AccessToken).
		Header().
		AddAccept("application/json").
		Build()

	rsp, err := client.GET("/tasks").Send()
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

	var tasks []Task
	err = json.Unmarshal(bytes, &tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (c *ListTasksOperation) Test(ctx *sdk.RunContext) (sdk.Json, error) {
	return c.Run(ctx)
}

func (c *ListTasksOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
