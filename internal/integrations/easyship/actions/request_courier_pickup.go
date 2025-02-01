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

	"github.com/wakflo/extensions/internal/integrations/easyship/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/integration"
)

type createCourierPickupActionProps struct {
	CourierID    string     `json:"courier-id"`
	ShipmentIDs  []string   `json:"shipment_ids"`
	SelectedDate *time.Time `json:"dueDate"`
}

type CreateCourierPickupAction struct{}

func (c CreateCourierPickupAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (c CreateCourierPickupAction) Name() string {
	return "Create Courier Pickup"
}

func (c CreateCourierPickupAction) Description() string {
	return "creates a courier pickup"
}

func (c CreateCourierPickupAction) Documentation() *integration.OperationDocumentation {
	return &integration.OperationDocumentation{
		Documentation: &createCourierDocs,
	}
}

func (c CreateCourierPickupAction) Icon() *string {
	return nil
}

func (c CreateCourierPickupAction) SampleData() (sdkcore.JSON, error) {
	return nil, nil
}

func (c CreateCourierPickupAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
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
	}
}

func (c CreateCourierPickupAction) Auth() *integration.Auth {
	return &integration.Auth{
		Inherit: true,
	}
}

func (c CreateCourierPickupAction) Perform(ctx integration.PerformContext) (sdkcore.JSON, error) {
	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing easyship api key")
	}

	input, err := integration.InputToTypeSafely[createCourierPickupActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	endpoint := "/pickups"

	shipmentData := map[string]interface{}{
		"courier_id":            input.CourierID,
		"easyship_shipment_ids": input.ShipmentIDs,
		"selected_date":         input.SelectedDate,
	}

	response, err := shared.PostRequest(endpoint, ctx.Auth.Extra["api-key"], shipmentData)
	if err != nil {
		return nil, fmt.Errorf("error requesting pickup:  %v", err)
	}

	return response, nil
}

func NewCreateCourierPickupAction() integration.Action {
	return &CreateCourierPickupAction{}
}
