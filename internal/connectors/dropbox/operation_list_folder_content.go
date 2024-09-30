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

type listFolderContentProps struct {
	Path      string `json:"path"`
	Recursive bool   `json:"recursive"`
	Limit     int    `json:"limit"`
}

type ListFolderOperation struct {
	options *sdk.OperationInfo
}

func NewListFolderOperation() *ListFolderOperation {
	limit := 2000
	return &ListFolderOperation{
		options: &sdk.OperationInfo{
			Name:        "List Folder Content",
			Description: "List the contents of a folder",
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"path": autoform.NewShortTextField().
					SetDisplayName("From Path").
					SetDescription("The path of the folder to be listed (e.g. /folder1). Use an empty string for the root folder.").
					SetRequired(true).Build(),
				"limit": autoform.NewNumberField().
					SetDisplayName("Limit").
					SetDescription("The maximum number of results to return (between 1 and 2000). Default is 2000 if not specified.").
					SetRequired(false).
					SetDefaultValue(limit).
					Build(),
				"recursive": autoform.NewBooleanField().
					SetDisplayName("Recursive").
					SetDescription("If set to true, the list folder operation will be applied recursively to all subfolders and the response will contain contents of all subfolders.").
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

func (c *ListFolderOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	input := sdk.InputToType[listFolderContentProps](ctx)

	folderContent, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	reqURL := "https://api.dropboxapi.com/2/files/list_folder"

	resp, err := listFolderContent(reqURL, ctx.Auth.AccessToken, folderContent)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *ListFolderOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *ListFolderOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
