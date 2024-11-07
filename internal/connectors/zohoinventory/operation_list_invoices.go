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
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type getInvoiceListOperationProps struct {
	OrganizationID string `json:"organization_id"`
}

type GetInvoiceListOperation struct {
	options *sdk.OperationInfo
}

func NewGetInvoiceListOperation() sdk.IOperation {
	return &GetInvoiceListOperation{
		options: &sdk.OperationInfo{
			Name:        "Get Invoice List",
			Description: "Retrieve a list of Invoices",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"organization_id": getOrganizationsInput(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *GetInvoiceListOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Token == nil {
		return nil, errors.New("missing Zoho auth token")
	}

	input := sdk.InputToType[getInvoiceListOperationProps](ctx)

	url := baseURL + "/v1/invoices/?organization_id=" + input.OrganizationID

	invoices, err := getZohoClient(ctx.Auth.AccessToken, url)
	if err != nil {
		return nil, err
	}
	return invoices, nil
}

func (c *GetInvoiceListOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *GetInvoiceListOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
