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

type listProjectUsersOperationProps struct {
	ProjectId string `json:"id"`
}

type ListProjectCollaboratorsOperation struct {
	options *sdk.OperationInfo
}

func NewListProjectCollaboratorsOperation() *ListProjectCollaboratorsOperation {
	return &ListProjectCollaboratorsOperation{
		options: &sdk.OperationInfo{
			Name:        "List Project Collaborators",
			Description: "Returns a list containing all collaborators of a shared project.",
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"id": autoform.NewShortTextField().
					SetDisplayName("Project ID").
					SetDescription("ID of the project.").
					SetRequired(true).Build(),
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

func (c *ListProjectCollaboratorsOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	input := sdk.InputToType[listProjectUsersOperationProps](ctx)

	client := fastshot.NewClient(baseApi).
		Auth().BearerToken(ctx.Auth.AccessToken).
		Header().
		AddAccept("application/json").
		Build()

	rsp, err := client.GET(fmt.Sprintf("/projects/%s/collaborators", input.ProjectId)).Send()
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

	var projects []Collaborator
	err = json.Unmarshal(bytes, &projects)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (c *ListProjectCollaboratorsOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *ListProjectCollaboratorsOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
