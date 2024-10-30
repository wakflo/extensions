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

package prisync

import (
	"net/http"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type getProductsSmartPricesOperationProps struct {
	Startfrom string `json:"start-from"`
}

type GetProductsSmartPricesOperation struct {
	options *sdk.OperationInfo
}

func NewGetProductsSmartPricesOperation() *GetProductsSmartPricesOperation {
	return &GetProductsSmartPricesOperation{
		options: &sdk.OperationInfo{
			Name:        "List Smart Prices for All Products ",
			Description: "Get Products smart prices",
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"start-from": autoform.NewShortTextField().
					SetDisplayName("Start From (Optional)").
					SetDescription("Offset for pagination. It can take 0 and exact multiples of 100 as a value.").
					SetRequired(false).Build(),
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

func (c *GetProductsSmartPricesOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	_ = sdk.InputToType[getProductsSmartPricesOperationProps](ctx)

	reqURL := "https://prisync.com/api/v2/list/smartprice/startFrom/0"
	resp, err := prisyncRequest(ctx.Auth.Extra["api-key"], ctx.Auth.Extra["api-token"], reqURL, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GetProductsSmartPricesOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *GetProductsSmartPricesOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
