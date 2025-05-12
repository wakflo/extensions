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

package asana

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/asana/actions"
	"github.com/wakflo/extensions/internal/integrations/asana/shared"
	"github.com/wakflo/extensions/internal/integrations/asana/triggers"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewAsana())

type Asana struct{}

func (n *Asana) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *Asana) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.AsanaSharedAuth,
	}
}

func (n *Asana) Triggers() []sdk.Trigger {
	return []sdk.Trigger{
		triggers.NewTaskCompletedTrigger(),
		triggers.NewTaskCreatedTrigger(),
		triggers.NewTaskUpdatedTrigger(),
	}
}

func (n *Asana) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewCreateTaskAction(),
		actions.NewGetTaskAction(),
		actions.NewUpdateTaskAction(),
		actions.NewListTasksAction(),
	}
}

func NewAsana() sdk.Integration {
	return &Asana{}
}
