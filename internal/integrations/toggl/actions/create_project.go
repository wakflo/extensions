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

	"github.com/wakflo/extensions/internal/integrations/toggl/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type createProjectActionProps struct {
	WorkspaceID string `json:"workspace_id"`
	Name        string `json:"name"`
	Active      bool   `json:"active"`
}

type CreateProjectAction struct{}

func (c *CreateProjectAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func (c CreateProjectAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (c CreateProjectAction) Name() string {
	return "Create Project"
}

func (c CreateProjectAction) Description() string {
	return "Create a new project"
}

func (c CreateProjectAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createProjectDocs,
	}
}

func (c CreateProjectAction) Icon() *string {
	return nil
}

func (c CreateProjectAction) SampleData() sdkcore.JSON {
	return nil
}

func (c CreateProjectAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"workspace_id": shared.GetWorkSpaceInput(),
		"name": autoform.NewShortTextField().
			SetDisplayName("Project Name").
			SetDescription("Project Name").
			SetRequired(true).
			Build(),
		"active": autoform.NewBooleanField().
			SetDisplayName("Active").
			SetDescription("make project active").
			SetDefaultValue(true).
			Build(),
	}
}

func (c CreateProjectAction) Auth() *sdk.Auth {
	return &sdk.Auth{
		Inherit: true,
	}
}

func (c CreateProjectAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing toggl api key")
	}
	apiKey := ctx.Auth.Extra["api-key"]
	input, err := sdk.InputToTypeSafely[createProjectActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	response, err := shared.CreateProjects(apiKey, input.WorkspaceID, input.Name, input.Active)
	if err != nil {
		log.Fatalf("error fetching data: %v", err)
	}

	return response, nil
}

func NewCreateProjectAction() sdk.Action {
	return &CreateProjectAction{}
}
