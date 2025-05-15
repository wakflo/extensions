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

type copyFolderActionProps struct {
	FromPath               string `json:"from_path"`
	ToPath                 string `json:"to_path"`
	AutoRename             bool   `json:"autorename"`
	AllowOwnershipTransfer bool   `json:"allow_ownership_transfer"`
}

type CopyFolderAction struct{}

// Metadata returns metadata about the action
func (a *CopyFolderAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "copy_folder",
		DisplayName:   "Copy a Folder",
		Description:   "copy folder",
		Type:          core.ActionTypeAction,
		Documentation: copyFolderDocs,
		SampleOutput:  map[string]any{},
		Settings:      core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *CopyFolderAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("copy_folder", "Copy a Folder")

	form.TextField("from_path", "From Path").
		Required(true).
		HelpText("The source path of the folder (e.g. /folder1/sourceFolder)")

	form.TextField("to_path", "To Path").
		Required(true).
		HelpText("The destination path for the copied folder (e.g. /folder2/destinationFolder)")

	form.CheckboxField("autorename", "Auto Rename").
		Required(false).
		HelpText("If there's a conflict, have the Dropbox server try to autorename the file to avoid conflict.")

	form.CheckboxField("allow_ownership_transfer", "Auto Ownership Transfer").
		Required(false).
		HelpText("Allows copy by owner even if it would result in an ownership transfer.")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *CopyFolderAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *CopyFolderAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[shared.FileMove](ctx)
	if err != nil {
		return nil, err
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	folders, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	reqURL := "/2/files/copy_v2"
	resp, err := shared.DropBoxClient(reqURL, authCtx.Token.AccessToken, folders)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NewCopyFolderAction() sdk.Action {
	return &CopyFolderAction{}
}
