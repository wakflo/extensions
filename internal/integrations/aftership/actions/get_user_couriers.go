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

type GetUserCouriersAction struct{}

// Metadata returns metadata about the action
func (c GetUserCouriersAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_user_couriers",
		DisplayName:   "Get User Activated Couriers",
		Description:   "get couriers activated by user",
		Type:          core.ActionTypeAction,
		Documentation: getUserCouriersDocs,
		SampleOutput:  nil,
		Settings:      core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (c GetUserCouriersAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_user_couriers", "Get User Activated Couriers")

	// No input properties needed for this action
	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (c GetUserCouriersAction) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Inherit: true,
	}
}

// Perform executes the action with the given context and input
func (c GetUserCouriersAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
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

	result, err := afterShipSdk.Courier.GetUserCouriers().Execute()
	if err != nil {
		return nil, err
	}

	return result.Couriers, nil
}

func NewGetUserCouriersAction() sdk.Action {
	return &GetUserCouriersAction{}
}
