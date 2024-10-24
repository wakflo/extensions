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

type createMeetingOperationProps struct{}

type GetProductsOperation struct {
	options *sdk.OperationInfo
}

func NewGetProductsOperation() *GetProductsOperation {
	return &GetProductsOperation{
		options: &sdk.OperationInfo{
			Name:        "Get all products",
			Description: "Get all products",
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"topic": autoform.NewShortTextField().
					SetDisplayName("").
					SetDescription("").
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

func (c *GetProductsOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	_ = sdk.InputToType[createMeetingOperationProps](ctx)

	reqURL := "https://logistics-api.flexport.com/logistics/api/2024-07/products"

	resp, err := flexportRequest(ctx.Auth.Extra["api-key"], reqURL, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GetProductsOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *GetProductsOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
