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

package zohoinventory

import (
	"errors"
	"fmt"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type getInvoiceOperationProps struct {
	OrganizationID string `json:"organization_id"`
	InvoiceID      string `json:"invoice_id"`
}

type GetInvoiceOperation struct {
	options *sdk.OperationInfo
}

func NewGetInvoiceOperation() sdk.IOperation {
	return &GetInvoiceOperation{
		options: &sdk.OperationInfo{
			Name:        "Get Invoice",
			Description: "Retrieve a specific invoice.",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"organization_id": getOrganizationsInput(),
				"invoice_id": autoform.NewShortTextField().
					SetDisplayName("Invoice ID").
					SetDescription("The ID of the invoice to retrieve to retrieve").
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

func (c *GetInvoiceOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing Zoho auth token")
	}

	input := sdk.InputToType[getInvoiceOperationProps](ctx)

	url := fmt.Sprintf("/v1/invoices/%s?organization_id=%s",
		input.InvoiceID, input.OrganizationID)
	invoice, err := getZohoClient(ctx.Auth.AccessToken, url)
	if err != nil {
		return nil, err
	}

	return invoice, nil
}

func (c *GetInvoiceOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *GetInvoiceOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
