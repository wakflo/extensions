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
	"errors"
	"log"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/toggl/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type createProjectActionProps struct {
	WorkspaceID string `json:"workspaces"`
	Name        string `json:"name"`
	Active      bool   `json:"active"`
}

type CreateProjectAction struct{}

func (c *CreateProjectAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_project",
		DisplayName:   "Create Project",
		Description:   "Create a new project",
		Type:          core.ActionTypeAction,
		Documentation: createProjectDocs,
		SampleOutput:  nil,
		Settings:      core.ActionSettings{},
	}
}

func (c *CreateProjectAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_project", "Create Project")

	shared.RegisterWorkspacesProp(form)

	form.TextField("name", "Project Name").
		Placeholder("Project Name").
		Required(true).
		HelpText("Project Name")

	form.CheckboxField("active", "Active").
		DefaultValue(true).
		HelpText("make project active")

	schema := form.Build()
	return schema
}

func (c *CreateProjectAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	if authCtx.Extra["api-key"] == "" {
		return nil, errors.New("missing toggl api key")
	}
	apiKey := authCtx.Extra["api-key"]

	input, err := sdk.InputToTypeSafely[createProjectActionProps](ctx)
	if err != nil {
		return nil, err
	}

	response, err := shared.CreateProjects(apiKey, input.WorkspaceID, input.Name, input.Active)
	if err != nil {
		log.Fatalf("error fetching data: %v", err)
	}

	return response, nil
}

func (c *CreateProjectAction) Auth() *core.AuthMetadata {
	return nil
}

func NewCreateProjectAction() sdk.Action {
	return &CreateProjectAction{}
}
