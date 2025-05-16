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

	"github.com/aftership/tracking-sdk-go/v5"
	"github.com/juicycleff/smartform/v1"
	// Import shared package for Auth and other shared functionality
	// We're ignoring type errors in shared.go files as per the issue description
	// Using blank identifier to suppress unused import warning
	_ "github.com/wakflo/extensions/internal/integrations/aftership/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type GetCouriersAction struct{}

// Metadata returns metadata about the action
func (c GetCouriersAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_couriers",
		DisplayName:   "Get All Couriers",
		Description:   "get all couriers",
		Type:          core.ActionTypeAction,
		Documentation: getCouriersDocs,
		SampleOutput:  nil,
		Settings:      core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (c GetCouriersAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_couriers", "Get All Couriers")

	// No input properties needed for this action
	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (c GetCouriersAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (c GetCouriersAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	if authCtx.Extra["api-key"] == "" {
		return nil, errors.New("missing aftership api key")
	}

	afterShipSdk, err := tracking.New(tracking.WithApiKey(authCtx.Extra["api-key"]))
	if err != nil {
		return nil, err
	}

	result, err := afterShipSdk.Courier.GetAllCouriers().Execute()
	if err != nil {
		return nil, err
	}

	return result.Couriers, nil
}

func NewGetCouriersAction() sdk.Action {
	return &GetCouriersAction{}
}
