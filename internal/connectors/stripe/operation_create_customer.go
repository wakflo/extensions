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

type createCustomerOperationProps struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	City        string `json:"city,omitempty"`
	Country     string `json:"country,omitempty"`
	Address     string `json:"line1,omitempty"`
	Phone       string `json:"phone,omitempty"`
	Description string `json:"description,omitempty"`
	State       string `json:"state,omitempty"`
	PostalCode  string `json:"postal_code,omitempty"`
}

type CreateCustomerOperation struct {
	options *sdk.OperationInfo
}

func NewCreateCustomerOperation() *CreateCustomerOperation {
	return &CreateCustomerOperation{
		options: &sdk.OperationInfo{
			Name:        "Create Customer",
			Description: "Create a customer",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"email": autoform.NewShortTextField().
					SetDisplayName("Email").
					SetDescription("email").
					SetRequired(true).
					Build(),
				"name": autoform.NewShortTextField().
					SetDisplayName("Name").
					SetDescription("name").
					SetRequired(true).
					Build(),
				"description": autoform.NewLongTextField().
					SetDisplayName("Description").
					SetDescription("description").
					Build(),
				"phone": autoform.NewShortTextField().
					SetDisplayName("Phone").
					SetDescription("phone").
					Build(),
				"line1": autoform.NewLongTextField().
					SetDisplayName("Address Line 1").
					SetRequired(false).
					Build(),
				"postal_code": autoform.NewShortTextField().
					SetDisplayName("Postal Code").
					SetRequired(false).
					Build(),
				"city": autoform.NewShortTextField().
					SetDisplayName("City").
					SetRequired(false).
					Build(),
				"state": autoform.NewShortTextField().
					SetDisplayName("State").
					SetRequired(false).
					Build(),
				"country": autoform.NewShortTextField().
					SetDisplayName("Country").
					SetRequired(false).
					Build(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *CreateCustomerOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing stripe secret api-key")
	}
	apiKey := ctx.Auth.Extra["api-key"]

	input := sdk.InputToType[createCustomerOperationProps](ctx)

	data := url.Values{}
	data.Set("name", input.Name)
	data.Set("email", input.Email)
	data.Set("description", input.Description)
	if input.Phone != "" {
		data.Set("phone", input.Phone)
	}
	if input.City != "" {
		data.Set("address[city]", input.City)
	}
	if input.Country != "" {
		data.Set("address[country]", input.Country)
	}
	if input.Address != "" {
		data.Set("address[line1]", input.Address)
	}
	if input.State != "" {
		data.Set("address[state]", input.State)
	}
	if input.PostalCode != "" {
		data.Set("address[postal_code]", input.PostalCode)
	}

	payload := []byte(data.Encode())

	reqURL := "/v1/customers"

	resp, err := stripClient(apiKey, reqURL, http.MethodPost, payload, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *CreateCustomerOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CreateCustomerOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
