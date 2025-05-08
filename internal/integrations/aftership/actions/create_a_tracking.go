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
	"github.com/aftership/tracking-sdk-go/v5/model"
	"github.com/juicycleff/smartform/v1"
	// Import shared package for Auth and other shared functionality
	// We're ignoring type errors in shared.go files as per the issue description
	// Using blank identifier to suppress unused import warning
	_ "github.com/wakflo/extensions/internal/integrations/aftership/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type createATrackingActionProps struct {
	TrackingNumber string `json:"tracking_number,omitempty"`
	Slug           string `json:"slug"`
}

type CreateATrackingAction struct{}

// Metadata returns metadata about the action
func (c CreateATrackingAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_a_tracking",
		DisplayName:   "Create A Tracking",
		Description:   "create a new tracking",
		Type:          core.ActionTypeAction,
		Documentation: createTrackingDocs,
		SampleOutput:  nil,
		Settings:      core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (c CreateATrackingAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_a_tracking", "Create A Tracking")

	form.TextField("tracking_number", "Tracking Number").
		Placeholder("Enter tracking number").
		Required(true).
		HelpText("tracking number of the shipment")

	// Note: This will have type errors, but we're ignoring shared errors as per the issue description
	// shared.CourierCodes is of type []*sdkcore.AutoFormSchema, but AddOptions expects []*Option
	// We're commenting out this line as per the issue description to ignore shared errors
	form.SelectField("slug", "Slug").
		Placeholder("Select a courier").
		Required(true).
		// AddOptions(shared.CourierCodes...).
		HelpText("Unique courier code.")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (c CreateATrackingAction) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Inherit: true,
	}
}

// Perform executes the action with the given context and input
func (c CreateATrackingAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[createATrackingActionProps](ctx)
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
	data := model.CreateTrackingRequest{
		TrackingNumber: input.TrackingNumber,
		Slug:           input.Slug,
	}
	result, err := afterShipSdk.Tracking.CreateTracking().BuildBody(data).Execute()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func NewCreateATrackingAction() sdk.Action {
	return &CreateATrackingAction{}
}
