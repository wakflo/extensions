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

package flexport

import (
	"net/http"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
	sdk2 "github.com/wakflo/go-sdk/sdk"
)

type getShipmentOperationProps struct {
	ShipmentID string `json:"shipmentId"`
}

type GetShipmentOperation struct {
	options *sdk.OperationInfo
}

func NewGetShipmentOperation() *GetShipmentOperation {
	return &GetShipmentOperation{
		options: &sdk.OperationInfo{
			Name:        "Get a shipment",
			Description: "Get a shipment",
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"shipmentId": autoform.NewShortTextField().
					SetDisplayName("Shipment ID").
					SetDescription("Shipment ID").
					SetRequired(true).Build(),
			},
			SampleOutput: map[string]interface{}{},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
			RequireAuth: true,
		},
	}
}

func (c *GetShipmentOperation) Run(ctx *sdk.RunContext) (sdk2.JSON, error) {
	input := sdk.InputToType[getShipmentOperationProps](ctx)

	reqURL := "/api/2024-07/inbounds/shipments/" + input.ShipmentID

	resp, err := flexportRequest(ctx.Auth.Extra["api-key"], reqURL, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GetShipmentOperation) Test(ctx *sdk.RunContext) (sdk2.JSON, error) {
	return c.Run(ctx)
}

func (c *GetShipmentOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
