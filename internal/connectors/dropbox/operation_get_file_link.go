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

type getFileLinkProps struct {
	Path string `json:"path"`
}

type GetFileLinkOperation struct {
	options *sdk.OperationInfo
}

func NewGetFileLinkOperation() *GetFileLinkOperation {
	return &GetFileLinkOperation{
		options: &sdk.OperationInfo{
			Name:        "Get a temporary file link",
			Description: "Get temporary file link",
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"path": autoform.NewShortTextField().
					SetDisplayName("Path").
					SetDescription("The path of the file (e.g. /folder1/file.txt)").
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

func (c *GetFileLinkOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	input := sdk.InputToType[getFileLinkProps](ctx)

	fileLink, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	reqURL := "https://api.dropboxapi.com/2/files/get_temporary_link"

	resp, err := dropBoxClient(reqURL, ctx.Auth.AccessToken, fileLink)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GetFileLinkOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *GetFileLinkOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
