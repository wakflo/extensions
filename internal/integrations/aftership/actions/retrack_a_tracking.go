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

	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type retrackATrackingActionProps struct {
	TrackingID string `json:"tracking_id"`
}

type RetrackATrackingAction struct{}

// Metadata returns metadata about the action
func (c RetrackATrackingAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "retrack_a_tracking",
		DisplayName:   "Retrack a specific tracking",
		Description:   "retrack an expired tracking. Max 3 times per tracking.",
		Type:          core.ActionTypeAction,
		Documentation: retrackATracking,
		SampleOutput:  nil,
		Settings:      core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (c RetrackATrackingAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("retrack_a_tracking", "Retrack a specific tracking")

	form.TextField("tracking_id", "Tracking ID").
		Placeholder("Enter tracking ID").
		Required(true).
		HelpText("tracking ID")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (c RetrackATrackingAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (c RetrackATrackingAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	if authCtx.Extra["api-key"] == "" {
		return nil, errors.New("missing aftership api key")
	}

	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[retrackATrackingActionProps](ctx)
	if err != nil {
		return nil, err
	}

	afterShipSdk, err := tracking.New(tracking.WithApiKey(authCtx.Extra["api-key"]))
	if err != nil {
		return nil, err
	}

	result, err := afterShipSdk.Tracking.
		RetrackTrackingById().
		BuildPath(input.TrackingID).
		Execute()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func NewRetrackATrackingAction() sdk.Action {
	return &RetrackATrackingAction{}
}
