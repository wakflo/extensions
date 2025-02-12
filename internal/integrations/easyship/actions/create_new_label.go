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

	"github.com/wakflo/extensions/internal/integrations/easyship/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type createNewLabelActionProps struct {
	ShipmentID string `json:"shipment_id"`
	CourierID  string `json:"courier_id"`
}

type CreateNewLabelAction struct{}

func (c *CreateNewLabelAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func (c CreateNewLabelAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (c CreateNewLabelAction) Name() string {
	return "Create A New label"
}

func (c CreateNewLabelAction) Description() string {
	return "creates a new label"
}

func (c CreateNewLabelAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &newLabelDocs,
	}
}

func (c CreateNewLabelAction) Icon() *string {
	return nil
}

func (c CreateNewLabelAction) SampleData() sdkcore.JSON {
	return nil
}

func (c CreateNewLabelAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"shipment_id": autoform.NewShortTextField().
			SetDisplayName("EasyShip Shipment ID").
			SetDescription("Readable identifier prefixed with ES (Easyship) and destination country code").
			SetRequired(true).
			Build(),
		"courier_id": autoform.NewShortTextField().
			SetDisplayName("Courier ID").
			SetDescription("Courier ID in case you need to overwrite the one suggested by default").
			SetRequired(false).
			Build(),
	}
}

func (c CreateNewLabelAction) Auth() *sdk.Auth {
	return &sdk.Auth{
		Inherit: true,
	}
}

func (c CreateNewLabelAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing easyship api key")
	}

	endpoint := "/2023-01/labels"

	input, err := sdk.InputToTypeSafely[createNewLabelActionProps](ctx.BaseContext)
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

	response, err := shared.PostRequest(endpoint, ctx.Auth.Extra["api-key"], newLabel)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func NewCreateNewLabelAction() sdk.Action {
	return &CreateNewLabelAction{}
}
