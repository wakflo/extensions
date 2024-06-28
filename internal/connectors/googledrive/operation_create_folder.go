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

package googledrive

import (
	"context"
	"errors"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type createFileOperationProps struct {
	FolderName        string  `json:"folderName"`
	ParentFolder      *string `json:"parentFolder"`
	IncludeTeamDrives bool    `json:"includeTeamDrives"`
}

type CreateFolderOperation struct {
	options *sdk.OperationInfo
}

func NewCreateFolderOperation() *CreateFolderOperation {
	return &CreateFolderOperation{
		options: &sdk.OperationInfo{
			Name:        "Create New Folder",
			Description: "Create a new empty folder in your Google Drive",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"folderName": autoform.NewShortTextField().
					SetDisplayName("Folder name").
					SetDescription("The name of the new folder.").
					SetRequired(true).
					Build(),
				"parentFolder":      getParentFoldersInput(),
				"includeTeamDrives": includeTeamFieldInput,
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *CreateFolderOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Token == nil {
		return nil, errors.New("missing google auth token")
	}

	input := sdk.InputToType[createFileOperationProps](ctx)
	driveService, err := drive.NewService(context.Background(), option.WithTokenSource(*ctx.Auth.TokenSource))
	if err != nil {
		return nil, err
	}

	var parents []string
	if input.ParentFolder != nil {
		parents = append(parents, *input.ParentFolder)
	}

	folder, err := driveService.Files.Create(&drive.File{
		MimeType: "application/vnd.google-apps.folder",
		Name:     input.FolderName,
		Parents:  parents,
	}).
		Fields("id, name, mimeType, webViewLink, kind, createdTime").
		SupportsAllDrives(input.IncludeTeamDrives).
		Do()
	return folder, err
}

func (c *CreateFolderOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CreateFolderOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
