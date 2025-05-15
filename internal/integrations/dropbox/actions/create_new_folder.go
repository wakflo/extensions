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

type createFolderActionProps struct {
	Path       string `json:"path"`
	Autorename bool   `json:"autorename"`
}

type CreateFolderAction struct{}

// Metadata returns metadata about the action
func (a *CreateFolderAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_folder",
		DisplayName:   "Create new Folder",
		Description:   "Create folder",
		Type:          core.ActionTypeAction,
		Documentation: createFolderDocs,
		SampleOutput:  map[string]any{},
		Settings:      core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *CreateFolderAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_folder", "Create new Folder")

	form.TextField("path", "Path").
		Required(true).
		HelpText("The path of the new folder e.g. /Homework/math")

	form.CheckboxField("autorename", "Auto Rename").
		Required(false).
		HelpText("If there's a conflict, have the Dropbox server try to autorename the file to avoid conflict.")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *CreateFolderAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *CreateFolderAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[createFolderActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	newFolder, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	reqURL := "/2/files/create_folder_v2"
	resp, err := shared.DropBoxClient(reqURL, authCtx.Token.AccessToken, newFolder)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NewCreateFolderAction() sdk.Action {
	return &CreateFolderAction{}
}
