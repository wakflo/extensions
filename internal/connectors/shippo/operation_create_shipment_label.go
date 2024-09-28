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

package shippo

import (
	"errors"
	"fmt"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type createShipmentLabelOperationProps struct {
	Rate string `json:"rate"`
}

type CreateShipmentLabelOperation struct {
	options *sdk.OperationInfo
}

func NewCreateShipmentLabelOperation() *CreateShipmentLabelOperation {
	return &CreateShipmentLabelOperation{
		options: &sdk.OperationInfo{
			Name:        "Create Shipment Label",
			Description: "creates a shipment label",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"rate": autoform.NewShortTextField().
					SetDisplayName("Rate").
					SetDescription("ID of the Rate object for which a Label has to be obtained.").
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

func (c *CreateShipmentLabelOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing shippo api key")
	}

	input := sdk.InputToType[createShipmentLabelOperationProps](ctx)

	endpoint := "/transactions"

	shipmentData := map[string]interface{}{
		"rate":            input.Rate,
		"async":           false,
		"label_file_type": "PDF",
	}

	response, err := CreateAShipment(endpoint, ctx.Auth.Extra["api-key"], shipmentData)
	if err != nil {
		return nil, fmt.Errorf("error creating shipment:  %v", err)
	}

	return response, nil
}

func (c *CreateShipmentLabelOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CreateShipmentLabelOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
