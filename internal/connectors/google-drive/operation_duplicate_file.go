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

package google_drive

import (
	"context"
	"errors"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type duplicateFileOperationProps struct {
	FileID            string  `json:"fileId"`
	FileName          string  `json:"fileName"`
	FolderId          *string `json:"folderId"`
	IncludeTeamDrives bool    `json:"includeTeamDrives"`
}

type DuplicateFileOperation struct {
	options *sdk.OperationInfo
}

func NewDuplicateFileOperation() *DuplicateFileOperation {
	return &DuplicateFileOperation{
		options: &sdk.OperationInfo{
			Name:        "Duplicate File",
			Description: "Duplicate a file from Google Drive. Returns the new file ID.",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"fileId": autoform.NewShortTextField().
					SetDisplayName("File ID").
					SetDescription("The ID of the file to duplicate").
					SetRequired(true).
					Build(),
				"fileName": autoform.NewShortTextField().
					SetDisplayName("Name").
					SetDescription("he name of the new file").
					SetRequired(true).
					Build(),
				"folderId":          getFoldersInput("Folder Id", "The ID of the folder where the file will be duplicated", false),
				"includeTeamDrives": includeTeamFieldInput,
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *DuplicateFileOperation) Run(ctx *sdk.RunContext) (sdk.Json, error) {
	if ctx.Auth.Token == nil {
		return nil, errors.New("missing google auth token")
	}

	input := sdk.InputToType[duplicateFileOperationProps](ctx)
	driveService, err := drive.NewService(context.Background(), option.WithTokenSource(*ctx.Auth.TokenSource))
	if err != nil {
		return nil, err
	}

	var parents []string
	if input.FolderId != nil {
		parents = append(parents, *input.FolderId)
	}

	in := &drive.File{
		Name:    input.FileName,
		Parents: parents,
	}

	return driveService.Files.Copy(input.FileID, in).SupportsAllDrives(input.IncludeTeamDrives).Do()
}

func (c *DuplicateFileOperation) Test(ctx *sdk.RunContext) (sdk.Json, error) {
	return c.Run(ctx)
}

func (c *DuplicateFileOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
