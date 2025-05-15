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

type deleteFolderActionProps struct {
	Path string `json:"path"`
}

type DeleteFolderAction struct{}

// Metadata returns metadata about the action
func (a *DeleteFolderAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "delete_folder",
		DisplayName:   "Delete A Folder",
		Description:   "delete an existing folder",
		Type:          core.ActionTypeAction,
		Documentation: deleteFolderDocs,
		SampleOutput:  map[string]any{},
		Settings:      core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *DeleteFolderAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("delete_folder", "Delete A Folder")

	form.TextField("path", "Path").
		Required(true).
		HelpText("The path of the folder to be deleted (e.g. /folder1)")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *DeleteFolderAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *DeleteFolderAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[deleteFolderActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	folder, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	reqURL := "/2/files/delete_v2"
	resp, err := shared.DropBoxClient(reqURL, authCtx.Token.AccessToken, folder)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NewDeleteFolderAction() sdk.Action {
	return &DeleteFolderAction{}
}
