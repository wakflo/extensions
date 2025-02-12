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

type getPurchaseInvoiceOperationProps struct {
	TaskID string `json:"id"`
}

type GetPurchaseInvoiceOperation struct {
	options *sdk.OperationInfo
}

func NewGetPurchaseInvoiceOperation() *GetPurchaseInvoiceOperation {
	return &GetPurchaseInvoiceOperation{
		options: &sdk.OperationInfo{
			Name:        "Get Purchase Invoice",
			Description: "Retrieves the purchase invoice",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"id": autoform.NewShortTextField().
					SetDisplayName("Task ID").
					SetDescription("The ID of the invoice to retrieve.").
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

func (c *GetPurchaseInvoiceOperation) Run(ctx *sdk.RunContext) (sdk2.JSON, error) {
	input := sdk.InputToType[getPurchaseInvoiceOperationProps](ctx)

	endpoint := "/ExternalApi/v2/purchase/invoice"
	accountID := ctx.Auth.Extra["account_id"]
	applicationKey := ctx.Auth.Extra["key"]
	queryParams := map[string]interface{}{
		"TaskID":                   input.TaskID,
		"CombineAdditionalCharges": true,
	}

	response, err := fetchData(endpoint, accountID, applicationKey, queryParams)
	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}

	return response, nil
}

func (c *GetPurchaseInvoiceOperation) Test(ctx *sdk.RunContext) (sdk2.JSON, error) {
	return c.Run(ctx)
}

func (c *GetPurchaseInvoiceOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
