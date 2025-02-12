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

type getFolderOperationProps struct {
	FolderID string `json:"folder-id"`
}

type GetFolderOperation struct {
	options *sdk.OperationInfo
}

func NewGetFolderOperation() *GetFolderOperation {
	return &GetFolderOperation{
		options: &sdk.OperationInfo{
			Name:        "Get Folder",
			Description: "Gets a specific folder",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"workspace-id": getWorkSpaceInput("Workspaces", "select a workspace", true),
				"space-id":     getSpacesInput("Spaces", "select a space", true),
				"folder-id":    getFoldersInput("Folders", "select a folder", true),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *GetFolderOperation) Run(ctx *sdk.RunContext) (sdk2.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing clickup auth token")
	}
	accessToken := ctx.Auth.AccessToken

	input := sdk.InputToType[getFolderOperationProps](ctx)
	url := "/v2/folder/" + input.FolderID

	folder, _ := getData(accessToken, url)

	return folder, nil
}

func (c *GetFolderOperation) Test(ctx *sdk.RunContext) (sdk2.JSON, error) {
	return c.Run(ctx)
}

func (c *GetFolderOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
