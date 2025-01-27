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
	"time"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"

	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type TriggerNewFolder struct {
	options *sdk.TriggerInfo
}

func NewTriggerNewFolder() *TriggerNewFolder {
	return &TriggerNewFolder{
		options: &sdk.TriggerInfo{
			Name:        "New Folder",
			Description: "triggers workflow when a new folder is produced",
			RequireAuth: true,
			Auth:        sharedAuth,
			Type:        sdkcore.TriggerTypePolling,
			Input: map[string]*sdkcore.AutoFormSchema{
				"parentFolder":      getParentFoldersInput(),
				"includeTeamDrives": includeTeamFieldInput,
			},
			Settings: &sdkcore.TriggerSettings{},
			ErrorSettings: &sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (t *TriggerNewFolder) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Token == nil {
		return nil, errors.New("missing google auth token")
	}

	input := sdk.InputToType[triggerNewFileProps](ctx)
	driveService, err := drive.NewService(context.Background(), option.WithTokenSource(*ctx.Auth.TokenSource))
	if err != nil {
		return nil, err
	}

	fmt.Printf("ctx.Metadata.LastRun %+v \n", ctx.Metadata.LastRun)

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

		qarr = append(qarr, fmt.Sprintf(`createdTime %v '%v'`, op, input.CreatedTime.UTC().Format(time.RFC3339)))
	}

	qarr = append(qarr, "trashed = false")
	q := fmt.Sprintf("%v %v", "mimeType='application/vnd.google-apps.folder'  and ", strings.Join(qarr, " and "))

	req := driveService.Files.List().
		IncludeItemsFromAllDrives(input.IncludeTeamDrives).
		SupportsAllDrives(input.IncludeTeamDrives).
		Fields("files(id, name, mimeType, webViewLink, kind, createdTime)").
		Q(q)

	files, err := req.Do()
	if err != nil {
		return nil, err
	}

	return files.Files, nil
}

func (t *TriggerNewFolder) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return t.Run(ctx)
}

func (t *TriggerNewFolder) OnEnabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *TriggerNewFolder) OnDisabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *TriggerNewFolder) GetInfo() *sdk.TriggerInfo {
	return t.options
}
