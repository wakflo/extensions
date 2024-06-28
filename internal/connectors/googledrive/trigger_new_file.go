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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	fastshot "github.com/opus-domini/fast-shot"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type triggerNewFileProps struct {
	ParentFolder       *string    `json:"parentFolder"`
	IncludeTeamDrives  bool       `json:"includeTeamDrives"`
	IncludeFileContent bool       `json:"includeFileContent"`
	CreatedTime        *time.Time `json:"createdTime"`
	CreatedTimeOp      *string    `json:"createdTimeOp"`
}

type TriggerNewFile struct {
	options *sdk.TriggerInfo
}

func NewTriggerNewFile() *TriggerNewFile {
	getParentFolders := func(ctx *sdkcore.DynamicFieldContext) (interface{}, error) {
		client := fastshot.NewClient("https://www.googleapis.com/drive/v3").
			Auth().BearerToken(ctx.Auth.AccessToken).
			Header().
			AddAccept("application/json").
			Build()

		rsp, err := client.GET("/files").Query().
			AddParams(map[string]string{
				"q": "mimeType='application/vnd.google-apps.folder' and trashed = false",
				/*"supportsTeamDrives": "true",
				"supportsAllDrives":  "true",*/
			}).Send()
		if err != nil {
			return nil, err
		}

		if rsp.IsError() {
			return nil, errors.New(rsp.StatusText())
		}

		bytes, err := io.ReadAll(rsp.RawBody())
		if err != nil {
			return nil, err
		}

		var body ListFileResponse
		err = json.Unmarshal(bytes, &body)
		if err != nil {
			return nil, err
		}

		return body.Files, nil
	}

	return &TriggerNewFile{
		options: &sdk.TriggerInfo{
			Name:        "New File",
			Description: "triggers workflow when a new file is produced",
			RequireAuth: true,
			Auth:        sharedAuth,
			Type:        sdkcore.TriggerTypeCron,
			Settings:    &sdkcore.TriggerSettings{},
			Input: map[string]*sdkcore.AutoFormSchema{
				"parentFolder": autoform.NewDynamicField(sdkcore.String).
					SetDisplayName("Parent Folder").
					SetDescription("select parent folder").
					SetDynamicOptions(&getParentFolders).
					SetDependsOn([]string{"connection"}).
					SetRequired(false).Build(),
				"includeTeamDrives": includeTeamFieldInput,
				"includeFileContent": autoform.NewBooleanField().
					SetDisplayName("Include File Content").
					SetDescription("Include the file content in the output. This will increase the time taken to fetch the files and might cause issues with large files.").
					SetDefaultValue(false).
					Build(),
			},
			SampleOutput: map[string]interface{}{
				"kind":     "drive#file",
				"mimeType": "image/jpeg",
				"id":       "1dpv4-sKJfKRwI9qx1vWqQhEGEn3EpbI5",
				"name":     "example.jpg",
			},
			ErrorSettings: &sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (t TriggerNewFile) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Token == nil {
		return nil, errors.New("missing google auth token")
	}

	input := sdk.InputToType[triggerNewFileProps](ctx)
	driveService, err := drive.NewService(context.Background(), option.WithTokenSource(*ctx.Auth.TokenSource))
	if err != nil {
		return nil, err
	}

	var qarr []string
	if input.ParentFolder != nil {
		qarr = append(qarr, fmt.Sprintf("'%v' in parents", *input.ParentFolder))
	}
	if input.CreatedTime == nil {
		input.CreatedTime = ctx.Metadata.LastRun
	}
	if input.CreatedTime != nil {
		op := ">"
		if input.CreatedTimeOp != nil {
			op = *input.CreatedTimeOp
		}
		qarr = append(qarr, fmt.Sprintf(`createdTime %v '%v'`, op, input.CreatedTime.UTC().Format("2006-01-02T15:04:05Z")))
	}

	qarr = append(qarr, "trashed = false")
	q := fmt.Sprintf("%v %v", "mimeType!='application/vnd.google-apps.folder'  and ", strings.Join(qarr, " and "))

	req := driveService.Files.List().
		IncludeTeamDriveItems(input.IncludeTeamDrives).
		SupportsAllDrives(input.IncludeTeamDrives).
		Fields("files(id, name, mimeType, webViewLink, kind, createdTime)").
		Q(q)

	files, err := req.Do()
	if err != nil {
		return nil, err
	}

	if input.IncludeFileContent {
		return handleFileContent(ctx, files.Files, driveService)
	}
	return files.Files, nil
}

func (t TriggerNewFile) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return t.Run(ctx)
}

func (t TriggerNewFile) OnEnabled(ctx *sdk.RunContext) error {
	return nil
}

func (t TriggerNewFile) OnDisabled(ctx *sdk.RunContext) error {
	return nil
}

func (t TriggerNewFile) GetInfo() *sdk.TriggerInfo {
	return t.options
}
