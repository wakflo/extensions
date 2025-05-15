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

type getFileLinkActionProps struct {
	Path string `json:"path"`
}

type GetFileLinkAction struct{}

// Metadata returns metadata about the action
func (a *GetFileLinkAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_file_link",
		DisplayName:   "Get a temporary file link",
		Description:   "Get temporary file link",
		Type:          core.ActionTypeAction,
		Documentation: getFileLinkDocs,
		SampleOutput:  map[string]any{},
		Settings:      core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *GetFileLinkAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_file_link", "Get a temporary file link")

	form.TextField("path", "Path").
		Required(true).
		HelpText("The path of the file (e.g. /folder1/file.txt)")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *GetFileLinkAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *GetFileLinkAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[getFileLinkActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	fileLink, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	reqURL := "/2/files/get_temporary_link"
	resp, err := shared.DropBoxClient(reqURL, authCtx.Token.AccessToken, fileLink)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NewGetFileLinkAction() sdk.Action {
	return &GetFileLinkAction{}
}
