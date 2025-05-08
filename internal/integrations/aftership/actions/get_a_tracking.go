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

type getATrackingActionProps struct {
	TrackingID string `json:"tracking_id"`
}

type GetATrackingAction struct{}

// Metadata returns metadata about the action
func (c GetATrackingAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_a_tracking",
		DisplayName:   "Get A Tracking",
		Description:   "get a specific tracking",
		Type:          core.ActionTypeAction,
		Documentation: getATrackingDocs,
		SampleOutput:  nil,
		Settings:      core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (c GetATrackingAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_a_tracking", "Get A Tracking")

	form.TextField("tracking_id", "Tracking ID").
		Placeholder("Enter tracking ID").
		Required(true).
		HelpText("tracking ID")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (c GetATrackingAction) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Inherit: true,
	}
}

// Perform executes the action with the given context and input
func (c GetATrackingAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[getATrackingActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	afterShipSdk, err := tracking.New(tracking.WithApiKey(authCtx.Extra["api-key"]))
	if err != nil {
		return nil, err
	}
	result, err := afterShipSdk.Tracking.
		GetTrackingById().
		BuildPath(input.TrackingID).
		Execute()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func NewGetATrackingAction() sdk.Action {
	return &GetATrackingAction{}
}
