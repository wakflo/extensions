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
)

type getSalesListOperationProps struct {
	PageLimit int `json:"page-limit"`
}

type GetSalesListOperation struct {
	options *sdk.OperationInfo
}

func NewGetSalesListOperation() *GetSalesListOperation {
	return &GetSalesListOperation{
		options: &sdk.OperationInfo{
			Name:        "Get Sale List",
			Description: "Retrieves a list of sales",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"page-limit": autoform.NewNumberField().
					SetDisplayName("Page limit").
					SetDescription(" Specifies the page limit for getting sales list.").
					Build(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *GetSalesListOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	input := sdk.InputToType[getSalesListOperationProps](ctx)

	endpoint := "/ExternalApi/SaleList"
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

func (c *GetSalesListOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *GetSalesListOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
