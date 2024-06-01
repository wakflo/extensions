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

type readFileOperationProps struct {
	FileID   string  `json:"fileId"`
	FileName *string `json:"fileName"`
}

type ReadFileOperation struct {
	options *sdk.OperationInfo
}

func NewReadFileOperation() *ReadFileOperation {
	return &ReadFileOperation{
		options: &sdk.OperationInfo{
			Name:        "Read file",
			Description: "Read a selected file from google drive file",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"fileId": autoform.NewShortTextField().
					SetDisplayName("File ID").
					SetDescription("File ID coming from | New File -> id |").
					SetRequired(true).
					Build(),
				"fileName": autoform.NewShortTextField().
					SetDisplayName("File Name").
					SetDescription("Destination File name").
					SetRequired(true).
					Build(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *ReadFileOperation) Run(ctx *sdk.RunContext) (sdk.Json, error) {
	if ctx.Auth.Token == nil {
		return nil, errors.New("missing google auth token")
	}

	input := sdk.InputToType[readFileOperationProps](ctx)
	driveService, err := drive.NewService(context.Background(), option.WithTokenSource(*ctx.Auth.TokenSource))
	if err != nil {
		return nil, err
	}

	return downloadFile(ctx, driveService, input.FileID, input.FileName)
}

func (c *ReadFileOperation) Test(ctx *sdk.RunContext) (sdk.Json, error) {
	return c.Run(ctx)
}

func (c *ReadFileOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
