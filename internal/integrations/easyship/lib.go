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

package easyship

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/easyship/actions"
	"github.com/wakflo/extensions/internal/integrations/easyship/shared"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewEasyShip())

type EasyShip struct{}

func (n *EasyShip) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *EasyShip) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.EasyShipSharedAuth,
	}
}

func (n *EasyShip) Triggers() []sdk.Trigger {
	return []sdk.Trigger{}
}

func (n *EasyShip) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewCreateNewLabelAction(),
		actions.NewCreateCourierPickupAction(),
	}
}

func NewEasyShip() sdk.Integration {
	return &EasyShip{}
}
