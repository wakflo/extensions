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
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type listProjectsOperationProps struct {
	ProjectId string `json:"project_id"`
}

type ListProjectsOperation struct {
	options *sdk.OperationInfo
}

func NewListProjectsOperation() *ListProjectsOperation {
	return &ListProjectsOperation{
		options: &sdk.OperationInfo{
			Name:        "List Project",
			Description: "Returns a list containing all user projects.",
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"project_id": getProjectsInput(),
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

func (c *ListProjectsOperation) Run(ctx *sdk.RunContext) (sdk.Json, error) {
	_ = sdk.InputToType[listProjectsOperationProps](ctx)

	client := fastshot.NewClient(baseApi).
		Auth().BearerToken(ctx.Auth.AccessToken).
		Header().
		AddAccept("application/json").
		Build()

	rsp, err := client.GET("/projects").Send()
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

	var projects []Project
	err = json.Unmarshal(bytes, &projects)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (c *ListProjectsOperation) Test(ctx *sdk.RunContext) (sdk.Json, error) {
	return c.Run(ctx)
}

func (c *ListProjectsOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
