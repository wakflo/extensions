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

package easyship

import (
	"errors"
	"fmt"
	"time"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type createCourierPickupOperationProps struct {
	CourierID    string     `json:"courier-id"`
	ShipmentIDs  []string   `json:"shipment_ids"`
	SelectedDate *time.Time `json:"dueDate"`
}

type CreateCourierPickupOperation struct {
	options *sdk.OperationInfo
}

func NewCreateCourierPickupOperation() *CreateCourierPickupOperation {
	return &CreateCourierPickupOperation{
		options: &sdk.OperationInfo{
			Name:        "Create Courier Pickup Request",
			Description: "Creates a courier pickup request",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"courier-id": autoform.NewShortTextField().
					SetDisplayName("Courier ID").
					SetDescription("Courier ID in case you need to overwrite the one suggested by default").
					SetRequired(true).
					Build(),
				"shipment_ids": autoform.NewArrayField().
					SetDisplayName("Shipment IDs").
					SetDescription("All shipments to be requested. Shipments must have the same courier and their labels must be pending or generated.").
					SetItems(
						autoform.NewShortTextField().
							SetDisplayName("Shipment ID").
							SetDescription("shipment Id").
							SetRequired(true).
							Build(),
					).
					SetRequired(false).Build(),
				"selected-date": autoform.NewDateTimeField().
					SetDisplayName("Selected date").
					SetDescription("Selected date for pickup").
					SetRequired(true).
					Build(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *CreateCourierPickupOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing easyship api key")
	}

	input := sdk.InputToType[createCourierPickupOperationProps](ctx)

	endpoint := "/pickups"

	shipmentData := map[string]interface{}{
		"courier_id":            input.CourierID,
		"easyship_shipment_ids": input.ShipmentIDs,
		"selected_date":         input.SelectedDate,
	}

	response, err := PostRequest(endpoint, ctx.Auth.Extra["api-key"], shipmentData)
	if err != nil {
		return nil, fmt.Errorf("error requesting pickup:  %v", err)
	}

	return response, nil
}

func (c *CreateCourierPickupOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CreateCourierPickupOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
