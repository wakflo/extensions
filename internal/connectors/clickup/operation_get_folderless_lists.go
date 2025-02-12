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

	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
	sdk2 "github.com/wakflo/go-sdk/sdk"
)

type getFolderlesslistsOperationProps struct {
	SpaceID string `json:"space-id"`
}

type GetFolderlesslistOperation struct {
	options *sdk.OperationInfo
}

func NewGetFolderlesslistOperation() *GetFolderlesslistOperation {
	return &GetFolderlesslistOperation{
		options: &sdk.OperationInfo{
			Name:        "Get Folderless Lists",
			Description: "Gets Folderless Lists",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"workspace-id": getWorkSpaceInput("Workspace", "select workspace", true),
				"space-id":     getSpacesInput("Spaces", "Select space", true),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *GetFolderlesslistOperation) Run(ctx *sdk.RunContext) (sdk2.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing clickup auth token")
	}
	accessToken := ctx.Auth.AccessToken

	input := sdk.InputToType[getFolderlesslistsOperationProps](ctx)

	folderlessList, _ := getSpace(accessToken, input.SpaceID)

	return folderlessList, nil
}

func (c *GetFolderlesslistOperation) Test(ctx *sdk.RunContext) (sdk2.JSON, error) {
	return c.Run(ctx)
}

func (c *GetFolderlesslistOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
