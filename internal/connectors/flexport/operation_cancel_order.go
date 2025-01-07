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
	"fmt"
	"net/http"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type cancelOrderOperationProps struct {
	OrderID string `json:"Id"`
}

type CancelOrderOperation struct {
	options *sdk.OperationInfo
}

func NewCancelOrderOperation() *CancelOrderOperation {
	return &CancelOrderOperation{
		options: &sdk.OperationInfo{
			Name:        "Cancel Order",
			Description: "Cancel a specific order",
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"Id": autoform.NewShortTextField().
					SetDisplayName("Order ID").
					SetDescription("Order ID").
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

func (c *CancelOrderOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	input := sdk.InputToType[cancelOrderOperationProps](ctx)

	reqURL := fmt.Sprintf("/api/2024-07/orders/%s/cancel", input.OrderID)
	resp, err := flexportRequest(ctx.Auth.Extra["api-key"], reqURL, http.MethodPost, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *CancelOrderOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CancelOrderOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
