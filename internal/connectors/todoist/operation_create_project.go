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

type CreateProjectOperation struct {
	options *sdk.OperationInfo
}

func NewCreateProjectOperation() *CreateProjectOperation {
	return &CreateProjectOperation{
		options: &sdk.OperationInfo{
			Name:        "Create Project",
			Description: "Create a todoist project",
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"name": autoform.NewShortTextField().
					SetDisplayName("Name").
					SetDescription("Name of the project.").
					SetRequired(true).Build(),
				"project_id": getProjectsInput(),
				"color": autoform.NewShortTextField().
					SetDisplayName("Color").
					SetDescription("The color of the project icon. Refer to the name column in the Colors guide for more info.").
					SetRequired(false).Build(),
				"is_favorite": autoform.NewBooleanField().
					SetDisplayName("Is Favourite").
					SetDescription("Whether the project is a favorite (a true or false value).").
					SetRequired(false).Build(),
				"view_style": autoform.NewSelectField().
					SetDisplayName("View Style").
					SetDescription("A string value (either list or board, default is list). This determines the way the project is displayed within the Todoist clients.").
					SetOptions(viewStyleOptions).
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

func (c *CreateProjectOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	input := sdk.InputToType[CreateProject](ctx)

	client := fastshot.NewClient(baseAPI).
		Auth().BearerToken(ctx.Auth.AccessToken).
		Header().
		AddAccept("application/json").
		Build()

	rsp, err := client.POST("/projects").Body().AsJSON(input).Send()
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

	var project Project
	err = json.Unmarshal(bytes, &project)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (c *CreateProjectOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CreateProjectOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
