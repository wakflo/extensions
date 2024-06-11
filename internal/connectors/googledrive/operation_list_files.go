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
	"fmt"
	"strings"

	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type listFilesOperationProps struct {
	FolderID          *string `json:"folderId"`
	IncludeTeamDrives bool    `json:"includeTeamDrives"`
}

type ListFilesOperation struct {
	options *sdk.OperationInfo
}

func NewListFilesOperation() *ListFilesOperation {
	return &ListFilesOperation{
		options: &sdk.OperationInfo{
			Name:        "List files",
			Description: "List files from a Google Drive folder",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"folderId":          getFoldersInput("Folder ID", "Folder ID coming from | New Folder -> id | (or any other source)", false),
				"includeTeamDrives": includeTeamFieldInput,
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *ListFilesOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Token == nil {
		return nil, errors.New("missing google auth token")
	}
	input := sdk.InputToType[listFilesOperationProps](ctx)

	driveService, err := drive.NewService(context.Background(), option.WithTokenSource(*ctx.Auth.TokenSource))
	if err != nil {
		return nil, err
	}
	var qarr []string
	if input.FolderID != nil {
		qarr = append(qarr, fmt.Sprintf("'%v' in parents", *input.FolderID))
	}
	qarr = append(qarr, "trashed = false")
	q := fmt.Sprintf("%v %v", "mimeType!='application/vnd.google-apps.folder'  and ", strings.Join(qarr, " and "))

	req := driveService.Files.List().
		Fields("files(id, name, mimeType, webViewLink, kind, createdTime)").
		SupportsAllDrives(input.IncludeTeamDrives).
		Q(q)

	file, err := req.Do()
	return file, err
}

func (c *ListFilesOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *ListFilesOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
