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

type getSalesOrderOperationProps struct {
	SaleID string `json:"id"`
}

type GetSalesOrderOperation struct {
	options *sdk.OperationInfo
}

func NewGetSalesOrderOperation() *GetSalesOrderOperation {
	return &GetSalesOrderOperation{
		options: &sdk.OperationInfo{
			Name:        "Get Sale Order",
			Description: "Retrieves sales order",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"id": autoform.NewShortTextField().
					SetDisplayName("Sales ID").
					SetDescription("The ID of the sale order to retrieve.").
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

func (c *GetSalesOrderOperation) Run(ctx *sdk.RunContext) (sdk2.JSON, error) {
	input := sdk.InputToType[getSalesOrderOperationProps](ctx)

	endpoint := "/ExternalApi/Sale"
	accountID := ctx.Auth.Extra["account_id"]
	applicationKey := ctx.Auth.Extra["key"]
	queryParams := map[string]interface{}{
		"ID": input.SaleID,
	}

	response, err := fetchData(endpoint, accountID, applicationKey, queryParams)
	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}

	return response, nil
}

func (c *GetSalesOrderOperation) Test(ctx *sdk.RunContext) (sdk2.JSON, error) {
	return c.Run(ctx)
}

func (c *GetSalesOrderOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
