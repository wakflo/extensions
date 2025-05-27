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

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/easyship/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type createNewLabelActionProps struct {
	ShipmentID string `json:"shipment_id"`
	CourierID  string `json:"courier_id"`
}

type CreateNewLabelAction struct{}

// Metadata returns metadata about the action
func (a *CreateNewLabelAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_new_label",
		DisplayName:   "Create A New label",
		Description:   "creates a new label",
		Type:          core.ActionTypeAction,
		Documentation: newLabelDocs,
		SampleOutput:  nil,
		Settings:      core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *CreateNewLabelAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_new_label", "Create A New label")

	// Add shipment_id field
	form.TextField("shipment_id", "EasyShip Shipment ID").
		Placeholder("Enter a shipment ID").
		Required(true).
		HelpText("Readable identifier prefixed with ES (Easyship) and destination country code")

	shared.RegisterCourierProps(form)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *CreateNewLabelAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *CreateNewLabelAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	if authCtx.Extra["api-key"] == "" {
		return nil, errors.New("missing easyship api key")
	}

	endpoint := "/2023-01/labels"

	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[createNewLabelActionProps](ctx)
	if err != nil {
		return nil, err
	}

	newLabel := map[string]interface{}{
		"shipments": []map[string]interface{}{
			{
				"easyship_shipment_id": input.ShipmentID,
				"courier_id":           input.CourierID,
			},
		},
	}

	response, err := shared.PostRequest(endpoint, authCtx.Extra["api-key"], newLabel)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func NewCreateNewLabelAction() sdk.Action {
	return &CreateNewLabelAction{}
}
