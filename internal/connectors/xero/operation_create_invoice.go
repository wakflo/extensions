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
	"time"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type createInvoiceOperationProps struct {
	ContactID string                   `json:"contact-id"`
	Contact   string                   `json:"contact"`
	LineItems []map[string]interface{} `json:"line_items"`
	DueDate   string                   `json:"due_date"`
	Status    string                   `json:"status"`
	Date      string                   `json:"date"`
	Email     string                   `json:"email"`
	Reference string                   `json:"reference"`
}

type CreateInvoiceOperation struct {
	options *sdk.OperationInfo
}

func NewCreateInvoiceOperation() sdk.IOperation {
	return &CreateInvoiceOperation{
		options: &sdk.OperationInfo{
			Name:        "Create Xero Invoice'",
			Description: "Create or Update Invoice'",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"contact-id": autoform.NewShortTextField().
					SetDisplayName("Contact ID").
					SetDescription("Contact ID").
					SetRequired(false).
					Build(),
				"email": autoform.NewShortTextField().
					SetDisplayName("Contact Email").
					SetDescription("Contact Email").
					SetRequired(false).
					Build(),
				"due_date": autoform.NewShortTextField().
					SetDisplayName("Due Date").
					SetDescription("Due date of the invoice. Format example: 2019-03-11").
					SetDefaultValue(time.Now().Format("2006-01-02")).
					SetRequired(true).
					Build(),
				"date": autoform.NewShortTextField().
					SetDisplayName("Date").
					SetDescription("Date the invoice was created. Format example: 2019-03-11").
					SetRequired(false).
					Build(),
				"reference": autoform.NewShortTextField().
					SetDisplayName("Invoice Reference").
					SetDescription("Reference number of the Invoice").
					SetRequired(false).
					Build(),
				"contact": autoform.NewShortTextField().
					SetDisplayName("Contact Full Name").
					SetDescription("Contact Name").
					SetRequired(true).
					Build(),
				"line_items": autoform.NewArrayField().
					SetDisplayName("Line Items").
					SetDescription("List of line items for the invoice.").
					SetRequired(true).
					SetItems(
						autoform.NewShortTextField().
							SetDisplayName("Label").
							SetDescription("Label").
							SetRequired(true).
							Build(),
					).
					SetDefaultValue([]map[string]interface{}{
						{
							"Description": "Default item description",
							"Quantity":    0,
							"UnitAmount":  0,
							"AccountCode": "200",
							"TaxType":     "NONE",
							"LineAmount":  0,
						},
					}).
					Build(),
				"status": autoform.NewSelectField().
					SetDisplayName("Status").
					SetOptions(xeroInvoiceStatus).
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

func (c *CreateInvoiceOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing Xero access token")
	}

	input := sdk.InputToType[createInvoiceOperationProps](ctx)

	body := map[string]interface{}{
		"Invoices": []map[string]interface{}{
			{
				"Type": "ACCREC",
				"Contact": map[string]interface{}{
					"Name": input.Contact,
				},
				"LineItems": func() []map[string]interface{} {
					if len(input.LineItems) > 0 {
						return input.LineItems
					}
					return []map[string]interface{}{}
				}(),
				"Date":      input.Date,
				"DueDate":   input.DueDate,
				"Reference": input.Reference,
				"Status":    input.Status,
			},
		},
	}

	response, err := createDraftInvoice(ctx.Auth.AccessToken, body)
	fmt.Println(response, err)

	return map[string]interface{}{
		"Report": "Invoice created Successfully",
	}, nil
}

func (c *CreateInvoiceOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CreateInvoiceOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
