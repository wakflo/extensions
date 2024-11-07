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

type deleteFolderProps struct {
	Path string `json:"path"`
}

type DeleteFolderOperation struct {
	options *sdk.OperationInfo
}

func NewDeleteFolderOperation() *DeleteFolderOperation {
	return &DeleteFolderOperation{
		options: &sdk.OperationInfo{
			Name:        "Delete A Folder",
			Description: "delete an existing folder",
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"path": autoform.NewShortTextField().
					SetDisplayName("Path").
					SetDescription("The path of the folder to be deleted (e.g. /folder1)").
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

func (c *DeleteFolderOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	input := sdk.InputToType[deleteFolderProps](ctx)

	folder, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	reqURL := "/2/files/delete_v2"

	resp, err := dropBoxClient(reqURL, ctx.Auth.AccessToken, folder)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *DeleteFolderOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *DeleteFolderOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
