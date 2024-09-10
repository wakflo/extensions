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

package xero

import (
	"errors"
	"fmt"

	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type emailInvoiceOperationProps struct {
	TenantID  string `json:"tenant_id"`
	InvoiceID string `json:"invoice_id"`
}

type EmailInvoiceOperation struct {
	options *sdk.OperationInfo
}

func NewEmailInvoiceOperation() sdk.IOperation {
	return &EmailInvoiceOperation{
		options: &sdk.OperationInfo{
			Name:        "Email Invoice",
			Description: "Sends a copy of a specific invoice to related contact via email",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"tenant_id":  getTenantInput("Organization", "select organization", true),
				"invoice_id": getInvoiceInput("Invoice", "The ID of the invoice to send.", true),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *EmailInvoiceOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing Xero access token")
	}

	input := sdk.InputToType[emailInvoiceOperationProps](ctx)

	endpoint := fmt.Sprintf("/Invoices/%s/Email", input.InvoiceID)

	err := sendInvoiceEmail(ctx.Auth.AccessToken, endpoint, input.TenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch invoice: %v", err)
	}

	return sdk.JSON(map[string]interface{}{
		"Report": "Invoice sent successfully",
	}), nil
}

func (c *EmailInvoiceOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *EmailInvoiceOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
