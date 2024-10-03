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

package toggl

import (
	"errors"
	"log"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type createProjectProps struct {
	WorkspaceID string `json:"workspace_id"`
	Name        string `json:"name"`
	Active      bool   `json:"active"`
}

type ListTicketsOperation struct {
	options *sdk.OperationInfo
}

func NewListTicketsOperation() *ListTicketsOperation {
	return &ListTicketsOperation{
		options: &sdk.OperationInfo{
			Name:        "Create Project",
			Description: "Create a project",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"workspace_id": getWorkSpaceInput(),
				"name": autoform.NewShortTextField().
					SetDisplayName("Project Name").
					SetDescription("Project Name").
					SetRequired(true).
					Build(),
				"active": autoform.NewBooleanField().
					SetDisplayName("Active").
					SetDescription("make project active").
					SetDefaultValue(true).
					Build(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *ListTicketsOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing toggl api key")
	}
	apiKey := ctx.Auth.Extra["api-key"]
	input := sdk.InputToType[createProjectProps](ctx)

	response, err := createProject(apiKey, input.WorkspaceID, input.Name, input.Active)
	if err != nil {
		log.Fatalf("error fetching data: %v", err)
	}

	return response, nil
}

func (c *ListTicketsOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *ListTicketsOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
