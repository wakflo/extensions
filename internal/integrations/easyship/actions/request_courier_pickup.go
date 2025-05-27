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
	"fmt"
	"time"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/easyship/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type createCourierPickupActionProps struct {
	CourierID    string     `json:"courier-id"`
	ShipmentIDs  []string   `json:"shipment_ids"`
	SelectedDate *time.Time `json:"dueDate"`
}

type CreateCourierPickupAction struct{}

// Metadata returns metadata about the action
func (a *CreateCourierPickupAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_courier_pickup",
		DisplayName:   "Create Courier Pickup",
		Description:   "creates a courier pickup",
		Type:          core.ActionTypeAction,
		Documentation: createCourierDocs,
		SampleOutput:  nil,
		Settings:      core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *CreateCourierPickupAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_courier_pickup", "Create Courier Pickup")

	shared.RegisterCourierProps(form)

	// Add shipment_ids field
	shipmentsArray := form.ArrayField("shipment_ids", "Shipment IDs")

	shipmentGroup := shipmentsArray.ObjectTemplate("shipment", "")

	shipmentGroup.TextField("shipment_id", "Shipment ID")

	form.DateTimeField("dueDate", "Selected date").
		Required(true).
		HelpText("Selected date for pickup")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *CreateCourierPickupAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *CreateCourierPickupAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	if authCtx.Extra["api-key"] == "" {
		return nil, errors.New("missing easyship api key")
	}

	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[createCourierPickupActionProps](ctx)
	if err != nil {
		return nil, err
	}

	endpoint := "/pickups"

	shipmentData := map[string]interface{}{
		"courier_id":            input.CourierID,
		"easyship_shipment_ids": input.ShipmentIDs,
		"selected_date":         input.SelectedDate,
	}

	response, err := shared.PostRequest(endpoint, authCtx.Extra["api-key"], shipmentData)
	if err != nil {
		return nil, fmt.Errorf("error requesting pickup:  %v", err)
	}

	return response, nil
}

func NewCreateCourierPickupAction() sdk.Action {
	return &CreateCourierPickupAction{}
}
