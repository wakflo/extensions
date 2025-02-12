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
	sdk2 "github.com/wakflo/go-sdk/sdk"
)

type CopyFileOperation struct {
	options *sdk.OperationInfo
}

func NewCopyFileOperation() *CopyFileOperation {
	return &CopyFileOperation{
		options: &sdk.OperationInfo{
			Name:        "Copy File",
			Description: "Create file",
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"from_path": autoform.NewShortTextField().
					SetDisplayName("From Path").
					SetDescription("The source path of the file (e.g. /folder1/sourcefile.txt)").
					SetRequired(true).Build(),
				"to_path": autoform.NewShortTextField().
					SetDisplayName("To Path").
					SetDescription("The destination path for the copied (e.g. /folder2/destinationfile.txt)").
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

func (c *CopyFileOperation) Run(ctx *sdk.RunContext) (sdk2.JSON, error) {
	input := sdk.InputToType[FileMove](ctx)

	file, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	reqURL := "/2/files/copy_v2"

	resp, err := dropBoxClient(reqURL, ctx.Auth.AccessToken, file)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *CopyFileOperation) Test(ctx *sdk.RunContext) (sdk2.JSON, error) {
	return c.Run(ctx)
}

func (c *CopyFileOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
