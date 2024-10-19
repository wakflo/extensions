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
)

type getOrderOperationProps struct {
	OrderID string `json:"Id"`
}

type GetOrderOperation struct {
	options *sdk.OperationInfo
}

func NewGetOrderOperation() *GetOrderOperation {
	return &GetOrderOperation{
		options: &sdk.OperationInfo{
			Name:        "Get order by ID",
			Description: "Get order by ID",
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"Id": autoform.NewShortTextField().
					SetDisplayName("Order ID").
					SetDescription("Order ID from flexport").
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

func (c *GetOrderOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	input := sdk.InputToType[getOrderOperationProps](ctx)

	reqURL := "https://logistics-api.flexport.com/logistics/api/2024-07/orders/" + input.OrderID
	resp, err := flexportRequest(ctx.Auth.Extra["api-key"], reqURL, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GetOrderOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *GetOrderOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
