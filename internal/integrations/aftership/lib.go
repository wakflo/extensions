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

package aftership

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/aftership/actions"
	"github.com/wakflo/extensions/internal/integrations/aftership/shared"
	"github.com/wakflo/go-sdk/sdk"
)

type AfterShip struct{}

func (n *AfterShip) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: true,
		Schema:   *shared.AfterShipSharedAuth,
	}
}

func (n *AfterShip) Triggers() []sdk.Trigger {
	return []sdk.Trigger{}
}

func (n *AfterShip) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewCreateATrackingAction(),
		actions.NewGetATrackingAction(),
		actions.NewDetectCourierAction(),
		actions.NewGetAllTrackingsAction(),
		actions.NewGetCouriersAction(),
		actions.NewMarkTrackingAsCompletedAction(),
		actions.NewRetrackATrackingAction(),
		actions.NewGetUserCouriersAction(),
	}
}

func NewAfterShip() sdk.Integration {
	return &AfterShip{}
}

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewAfterShip(), Flow, ReadME)
