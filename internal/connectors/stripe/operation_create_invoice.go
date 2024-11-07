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

package stripe

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type createInvoiceOperationProps struct {
	CustomerID  string `json:"customer"`
	Currency    string `json:"currency"`
	Description string `json:"description"`
}

type CreateInvoiceOperation struct {
	options *sdk.OperationInfo
}

func NewCreateInvoiceOperation() *CreateInvoiceOperation {
	return &CreateInvoiceOperation{
		options: &sdk.OperationInfo{
			Name:        "Create Invoice",
			Description: "Create an invoice",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"customer": autoform.NewShortTextField().
					SetDisplayName("Customer ID").
					SetDescription("Stripe customer ID").
					SetRequired(true).
					Build(),
				"currency": autoform.NewShortTextField().
					SetDisplayName("Currency").
					SetDescription("Currency for the invoice (e.g., USD)").
					SetRequired(true).
					Build(),
				"description": autoform.NewLongTextField().
					SetDisplayName("Description").
					SetDescription("description").
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
	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing stripe secret api-key")
	}
	apiKey := ctx.Auth.Extra["api-key"]

	input := sdk.InputToType[createInvoiceOperationProps](ctx)

	data := url.Values{}
	data.Set("customer", input.CustomerID)
	data.Set("currency", input.Currency)

	if input.Description != "" {
		data.Set("description", input.Description)
	}

	payload := []byte(data.Encode())

	reqURL := "/v1/invoices"

	resp, err := stripClient(apiKey, reqURL, http.MethodPost, payload, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *CreateInvoiceOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CreateInvoiceOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
