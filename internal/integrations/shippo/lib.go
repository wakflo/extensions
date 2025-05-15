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

package shippo

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/shippo/actions"
	"github.com/wakflo/extensions/internal/integrations/shippo/shared"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewShippo())

type Shippo struct{}

func (n *Shippo) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *Shippo) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.ShippoSharedAuth,
	}
}

func (n *Shippo) Triggers() []sdk.Trigger {
	return []sdk.Trigger{}
}

func (n *Shippo) Actions() []sdk.Action {
	return []sdk.Action{
		// actions.NewCreateNewShipmentAction(),
		actions.NewCreateShipmentLabelAction(),
	}
}

func NewShippo() sdk.Integration {
	return &Shippo{}
}
