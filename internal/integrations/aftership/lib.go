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
	// Import shared package for Auth and other shared functionality
	// We're ignoring type errors in shared.go files as per the issue description
	// Using blank identifier to suppress unused import warning
	_ "github.com/wakflo/extensions/internal/integrations/aftership/shared"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewAfterShip())

type AfterShip struct{}

func (n *AfterShip) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *AfterShip) Auth() *core.AuthMetadata {
	// Note: This will have type errors, but we're ignoring shared errors as per the issue description
	// shared.AfterShipSharedAuth is of type *sdkcore.AutoFormSchema, but Schema expects *smartform.FormSchema
	// We're returning nil for Schema as per the issue description to ignore shared errors
	return &core.AuthMetadata{
		Required: true,
		Schema:   nil, // shared.AfterShipSharedAuth,
	}
}

func (n *AfterShip) Triggers() []sdk.Trigger {
	return []sdk.Trigger{}
}

func (n *AfterShip) Actions() []sdk.Action {
	// Note: Only including actions that have been migrated to the new SDK v2 interface
	// The other actions are commented out because they're still using the old SDK interface
	// and don't implement the new SDK v2 Action interface, which requires a Metadata() method
	return []sdk.Action{
		actions.NewCreateATrackingAction(),
		actions.NewGetATrackingAction(),
		// actions.NewDetectCourierAction(),
		// actions.NewGetAllTrackingsAction(),
		// actions.NewGetCouriersAction(),
		// actions.NewMarkTrackingAsCompletedAction(),
		// actions.NewRetrackATrackingAction(),
		// actions.NewGetUserCouriersAction(),
	}
}

func NewAfterShip() sdk.Integration {
	return &AfterShip{}
}
