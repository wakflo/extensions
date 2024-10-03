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

package dropbox

import (
	"encoding/json"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type createNewFolderProps struct {
	Autorename bool   `json:"autorename"`
	Path       string `json:"path"`
}

type CreateFolderOperation struct {
	options *sdk.OperationInfo
}

func NewCreateFolderOperation() *CreateFolderOperation {
	return &CreateFolderOperation{
		options: &sdk.OperationInfo{
			Name:        "Create new Folder",
			Description: "Create folder",
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"path": autoform.NewShortTextField().
					SetDisplayName("Path").
					SetDescription("The path of the new folder e.g. /Homework/math").
					SetRequired(true).Build(),
				"autorename": autoform.NewBooleanField().
					SetDisplayName("Auto Rename").
					SetDescription("If there's a conflict, have the Dropbox server try to autorename the file to avoid conflict.").
					SetRequired(false).
					SetDefaultValue(false).
					Build(),
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

func (c *CreateFolderOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	input := sdk.InputToType[createNewFolderProps](ctx)

	newFolder, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	reqURL := "https://api.dropboxapi.com/2/files/create_folder_v2"

	resp, err := dropBoxClient(reqURL, ctx.Auth.AccessToken, newFolder)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *CreateFolderOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CreateFolderOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
