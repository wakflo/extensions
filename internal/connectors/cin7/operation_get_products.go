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

package cin7

import (
	"log"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
	sdk2 "github.com/wakflo/go-sdk/sdk"
)

type getProductsOperationProps struct {
	PageLimit int `json:"page-limit"`
}

type GetProductsOperation struct {
	options *sdk.OperationInfo
}

func NewGetProductsOperation() *GetProductsOperation {
	return &GetProductsOperation{
		options: &sdk.OperationInfo{
			Name:        "Get Products",
			Description: "Retrieves a list of products from Cin7",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"page-limit": autoform.NewNumberField().
					SetDisplayName("Page limit").
					SetDescription(" Specifies the page limit for products.").
					Build(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *GetProductsOperation) Run(ctx *sdk.RunContext) (sdk2.JSON, error) {
	input := sdk.InputToType[getProductsOperationProps](ctx)

	endpoint := "/ExternalApi/Products"
	accountID := ctx.Auth.Extra["account_id"]
	applicationKey := ctx.Auth.Extra["key"]

	if input.PageLimit == 0 {
		input.PageLimit = 100
	}

	queryParams := map[string]interface{}{
		"Page":  1,
		"Limit": input.PageLimit,
	}

	response, err := fetchData(endpoint, accountID, applicationKey, queryParams)
	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}

	return response, nil
}

func (c *GetProductsOperation) Test(ctx *sdk.RunContext) (sdk2.JSON, error) {
	return c.Run(ctx)
}

func (c *GetProductsOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
