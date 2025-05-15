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

package actions

import (
	"encoding/json"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/dropbox/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type deleteFileActionProps struct {
	Path string `json:"path"`
}

type DeleteFileAction struct{}

// Metadata returns metadata about the action
func (a *DeleteFileAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "delete_file",
		DisplayName:   "Delete A file",
		Description:   "delete an existing file",
		Type:          core.ActionTypeAction,
		Documentation: deleteFileDocs,
		SampleOutput:  map[string]any{},
		Settings:      core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *DeleteFileAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("delete_file", "Delete A file")

	form.TextField("path", "Path").
		Required(true).
		HelpText("The path of the file to be deleted (e.g. /folder1/file.txt)")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *DeleteFileAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *DeleteFileAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[deleteFileActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	deletedFile, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	reqURL := "/2/files/delete_v2"
	resp, err := shared.DropBoxClient(reqURL, authCtx.Token.AccessToken, deletedFile)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NewDeleteFileAction() sdk.Action {
	return &DeleteFileAction{}
}
