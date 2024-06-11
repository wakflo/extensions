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

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type uploadFileOperationProps struct {
	FileName          string  `json:"fileName"`
	File              string  `json:"file"`
	ParentFolder      *string `json:"parentFolder"`
	IncludeTeamDrives bool    `json:"includeTeamDrives"`
}

type UploadFileOperation struct {
	options *sdk.OperationInfo
}

func NewUploadFileOperation() *UploadFileOperation {
	return &UploadFileOperation{
		options: &sdk.OperationInfo{
			Name:        "Upload file",
			Description: "Upload a file in your Google Drive",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"fileName": autoform.NewShortTextField().
					SetDisplayName("Name").
					SetDescription("The name of the new file").
					SetRequired(true).
					Build(),
				"file": autoform.NewFileField().
					SetDisplayName("File").
					SetDescription("The file URL or base64 to upload").
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

func (c *UploadFileOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Token == nil {
		return nil, errors.New("missing google auth token")
	}

	input := sdk.InputToType[uploadFileOperationProps](ctx)
	driveService, err := drive.NewService(context.Background(), option.WithTokenSource(*ctx.Auth.TokenSource))
	if err != nil {
		return nil, err
	}

	fileData, err := sdk.StringToFile(input.File)
	if err != nil {
		return nil, err
	}

	var parents []string
	if input.ParentFolder != nil {
		parents = append(parents, *input.ParentFolder)
	}

	in := &drive.File{
		MimeType: fileData.Mime,
		Name:     input.FileName,
		Parents:  parents,
	}

	return driveService.Files.Create(in).
		Media(fileData.Data).
		SupportsAllDrives(input.IncludeTeamDrives).
		Do()
}

func (c *UploadFileOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *UploadFileOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
