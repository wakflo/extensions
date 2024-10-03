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

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type createNewLabelOperationProps struct {
	ShipmentID string `json:"shipment_id"`
	CourierID  string `json:"courier_id"`
}

type CreateNewLabelOperation struct {
	options *sdk.OperationInfo
}

func NewCreateNewLabelOperation() *CreateNewLabelOperation {
	return &CreateNewLabelOperation{
		options: &sdk.OperationInfo{
			Name:        "Create A New label",
			Description: "creates a new label",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"shipment_id": autoform.NewShortTextField().
					SetDisplayName("Easyship Shipment ID").
					SetDescription("Readable identifier prefixed with ES (Easyship) and destination country code").
					SetRequired(true).
					Build(),
				"courier_id": autoform.NewShortTextField().
					SetDisplayName("Courier ID").
					SetDescription("Courier ID in case you need to overwrite the one suggested by default").
					SetRequired(false).
					Build(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *CreateNewLabelOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing easyship api key")
	}

	endpoint := "/2023-01/labels"

	input := sdk.InputToType[createNewLabelOperationProps](ctx)

	newLabel := map[string]interface{}{
		"shipments": []map[string]interface{}{
			{
				"easyship_shipment_id": input.ShipmentID,
				"courier_id":           input.CourierID,
			},
		},
	}

	response, err := PostRequest(endpoint, ctx.Auth.Extra["api-key"], newLabel)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *CreateNewLabelOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CreateNewLabelOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
