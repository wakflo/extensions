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
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type detectCourierActionProps struct {
	TrackingNumber string `json:"tracking_number,omitempty"`
}

type DetectCourierAction struct{}

// Metadata returns metadata about the action
func (c DetectCourierAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "detect_courier",
		DisplayName:   "Detect Courier",
		Description:   "Returns a list of matched couriers based on tracking number format",
		Type:          core.ActionTypeAction,
		Documentation: detectCourierDocs,
		SampleOutput:  nil,
		Settings:      core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (c DetectCourierAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("detect_courier", "Detect Courier")

	form.TextField("tracking_number", "Tracking Number").
		Placeholder("Enter tracking number").
		Required(true).
		HelpText("tracking number of the shipment")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (c DetectCourierAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (c DetectCourierAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[detectCourierActionProps](ctx)
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

	result, err := afterShipSdk.Courier.
		DetectCourier().
		BuildBody(model.DetectCourierRequest{
			TrackingNumber: input.TrackingNumber,
		}).
		Execute()
	if err != nil {
		return nil, err
	}

	return result.Couriers, nil
}

func NewDetectCourierAction() sdk.Action {
	return &DetectCourierAction{}
}
