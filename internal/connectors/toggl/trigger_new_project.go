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

package toggl

import (
	"errors"
	"log"

	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type listProjectProps struct {
	WorkspaceID string `json:"workspace_id"`
}

type ProjectTrigger struct {
	options *sdk.TriggerInfo
}

func NewProjectTrigger() *ProjectTrigger {
	return &ProjectTrigger{
		options: &sdk.TriggerInfo{
			Name:        "New Project",
			Description: "triggers workflow when a new project is created, modified or deleted",
			RequireAuth: true,
			Auth:        sharedAuth,
			Strategy:    sdkcore.TriggerStrategyPolling,
			Input: map[string]*sdkcore.AutoFormSchema{
				"workspace_id": getWorkSpaceInput(),
			},
			Settings: &sdkcore.TriggerSettings{},
			ErrorSettings: &sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (t *ProjectTrigger) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing toggl api key")
	}
	apiKey := ctx.Auth.Extra["api-key"]
	input := sdk.InputToType[listProjectProps](ctx)
	lastRunTime := ctx.Metadata.LastRun

	var updatedTime int64
	if lastRunTime != nil {
		updatedTime = lastRunTime.UTC().Unix()
	}

	response, err := getProjects(apiKey, input.WorkspaceID, updatedTime)
	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}

	return response, nil
}

func (t *ProjectTrigger) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return t.Run(ctx)
}

func (t *ProjectTrigger) OnEnabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *ProjectTrigger) OnDisabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *ProjectTrigger) GetInfo() *sdk.TriggerInfo {
	return t.options
}
