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

type MoveFolderOperation struct {
	options *sdk.OperationInfo
}

func NewMoveFolderOperation() *MoveFolderOperation {
	return &MoveFolderOperation{
		options: &sdk.OperationInfo{
			Name:        "Move a Folder",
			Description: "move folder",
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"from_path": autoform.NewShortTextField().
					SetDisplayName("From Path").
					SetDescription("The current path of the folder (e.g. /folder1/sourceFolder)").
					SetRequired(true).Build(),
				"to_path": autoform.NewShortTextField().
					SetDisplayName("To Path").
					SetDescription("The new path for the folder (e.g. /folder2/destinationFolder)").
					SetRequired(true).Build(),
				"autorename": autoform.NewBooleanField().
					SetDisplayName("Auto Rename").
					SetDescription("If there's a conflict, have the Dropbox server try to autorename the file to avoid conflict.").
					SetRequired(false).
					SetDefaultValue(false).
					Build(),
				"allow_ownership_transfer": autoform.NewBooleanField().
					SetDisplayName("Auto Ownership Transfer").
					SetDescription(" Allows copy by owner even if it would result in an ownership transfer.").
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

func (c *MoveFolderOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	input := sdk.InputToType[FileMove](ctx)

	folders, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	reqURL := "https://api.dropboxapi.com/2/files/move_v2"
	resp, err := dropBoxClient(reqURL, ctx.Auth.AccessToken, folders)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *MoveFolderOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *MoveFolderOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
