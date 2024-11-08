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

type getOrderByExternalIDOperationProps struct {
	ExternalOrderID string `json:"Id"`
}

type GetOrderByExternalIDOperation struct {
	options *sdk.OperationInfo
}

func NewGetOrderByExternalIDOperation() *GetOrderByExternalIDOperation {
	return &GetOrderByExternalIDOperation{
		options: &sdk.OperationInfo{
			Name:        "Get order by External ID",
			Description: "Get order using the external order ID given during order creation.",
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"Id": autoform.NewShortTextField().
					SetDisplayName("External Order ID").
					SetDescription("External Order ID from store(woocommerce, shopify etc.)").
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

func (c *GetOrderByExternalIDOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	input := sdk.InputToType[getOrderByExternalIDOperationProps](ctx)

	reqURL := "/api/2024-07/orders/external_id/" + input.ExternalOrderID
	resp, err := flexportRequest(ctx.Auth.Extra["api-key"], reqURL, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GetOrderByExternalIDOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *GetOrderByExternalIDOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
