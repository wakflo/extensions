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
	"strings"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/easyship/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type createCourierPickupActionProps struct {
	CourierServiceID string `json:"courier_service_id"`
	ShipmentIDs      string `json:"easyship_shipment_ids"`
	SelectedDate     string `json:"selected_date"`
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

func (a *CreateCourierPickupAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_courier_pickup", "Create Courier Pickup")

	form.TextField("courier_service_id", "Courier Service ID").
		Placeholder("01563646-58c1-4607-8fe0-cae3e33c0001").
		Required(true).
		HelpText("Courier service ID (UUID format)")

	form.TextareaField("easyship_shipment_ids", "Shipment IDs").
		Required(true).
		HelpText("Enter shipment IDs separated by commas or new lines").
		Placeholder("ESSG10006001, ESSG10006002, ESSG10006003")

	form.DateField("selected_date", "Pickup Date").
		Required(true).
		HelpText("Selected date for pickup")

	schema := form.Build()
	return schema
}

func (a *CreateCourierPickupAction) Auth() *core.AuthMetadata {
	return nil
}

func parseShipmentIDs(text string) ([]string, error) {
	if strings.TrimSpace(text) == "" {
		return nil, errors.New("shipment IDs text cannot be empty")
	}

	var result []string
	parts := strings.FieldsFunc(text, func(c rune) bool {
		return c == ',' || c == '\n' || c == '\r' || c == ';'
	})

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	if len(result) == 0 {
		return nil, errors.New("no valid shipment IDs found")
	}

	return result, nil
}

func (a *CreateCourierPickupAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	if authCtx.Extra["api-key"] == "" {
		return nil, errors.New("missing easyship api key")
	}

	input, err := sdk.InputToTypeSafely[createCourierPickupActionProps](ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input: %v", err)
	}

	shipmentIDStrings, err := parseShipmentIDs(input.ShipmentIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse shipment IDs: %v", err)
	}

	endpoint := "/pickups"

	shipmentData := map[string]interface{}{
		"courier_service_id":    input.CourierServiceID,
		"easyship_shipment_ids": shipmentIDStrings,
		"selected_date":         input.SelectedDate,
	}

	response, err := shared.PostRequest(endpoint, authCtx.Extra["api-key"], shipmentData)
	if err != nil {
		return nil, fmt.Errorf("error requesting pickup: %v", err)
	}

	return response, nil
}

func NewCreateCourierPickupAction() sdk.Action {
	return &CreateCourierPickupAction{}
}
