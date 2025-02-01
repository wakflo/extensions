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
	"github.com/wakflo/extensions/internal/integrations/asana/actions"
	"github.com/wakflo/extensions/internal/integrations/asana/shared"
	"github.com/wakflo/go-sdk/integration"
)

type Asana struct{}

func (n *Asana) Auth() *integration.Auth {
	return &integration.Auth{
		Required: true,
		Schema:   *shared.AsanaSharedAuth,
	}
}

func (n *Asana) Triggers() []integration.Trigger {
	return []integration.Trigger{}
}

func (n *Asana) Actions() []integration.Action {
	return []integration.Action{
		actions.NewCreateTaskAction(),
	}
}

func NewAsana() integration.Integration {
	return &Asana{}
}

var Integration = integration.Register(NewAsana())
