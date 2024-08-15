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

type createFolderlesslistOperationProps struct {
	SpaceID string `json:"space-id"`
	Name    string `json:"name"`
}

type CreateFolderlessListOperation struct {
	options *sdk.OperationInfo
}

func NewCreateFolderlessListOperation() *CreateFolderlessListOperation {
	return &CreateFolderlessListOperation{
		options: &sdk.OperationInfo{
			Name:        "Create Folderless List",
			Description: "Create a new folderless list in a ClickUp workspace and space",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"space-id": autoform.NewShortTextField().
					SetDisplayName("Space ID").
					SetDescription("The ID of the Space").
					SetRequired(true).
					Build(),
				"name": autoform.NewShortTextField().
					SetDisplayName("List Name").
					SetDescription("The name of the list").
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

func (c *CreateFolderlessListOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing clickup auth token")
	}
	accessToken := ctx.Auth.AccessToken

	input := sdk.InputToType[createFolderlesslistOperationProps](ctx)
	reqURL := "https://api.clickup.com/api/v2/space/" + input.SpaceID + "/list"

	response, _ := createItem(accessToken, input.Name, reqURL)

	return response, nil
}

func (c *CreateFolderlessListOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CreateFolderlessListOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
