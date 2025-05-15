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

type listFolderActionProps struct {
	Path      string `json:"path"`
	Limit     int    `json:"limit"`
	Recursive bool   `json:"recursive"`
}

type ListFolderAction struct{}

// Metadata returns metadata about the action
func (a *ListFolderAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "list_folder",
		DisplayName:   "List Folder Content",
		Description:   "List the contents of a folder",
		Type:          core.ActionTypeAction,
		Documentation: listFolderDocs,
		SampleOutput:  map[string]any{},
		Settings:      core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *ListFolderAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("list_folder", "List Folder Content")

	form.TextField("path", "From Path").
		Required(true).
		HelpText("The path of the folder to be listed (e.g. /folder1). Use an empty string for the root folder.")

	form.NumberField("limit", "Limit").
		Required(false).
		DefaultValue(2000).
		HelpText("The maximum number of results to return (between 1 and 2000). Default is 2000 if not specified.")

	form.CheckboxField("recursive", "Recursive").
		Required(false).
		HelpText("If set to true, the list folder operation will be applied recursively to all subfolders and the response will contain contents of all subfolders.")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *ListFolderAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *ListFolderAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[listFolderActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	folderContent, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	reqURL := "/2/files/list_folder"
	resp, err := shared.ListFolderContent(reqURL, authCtx.Token.AccessToken, folderContent)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NewListFolderAction() sdk.Action {
	return &ListFolderAction{}
}
