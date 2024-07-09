// Copyright 2022-present Wakflo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package shopify

import (
	"context"
	"errors"
	// "fmt"
	// "strings"
	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type updateCustomerOperationProps struct {
	CustomerID uint64 `json:"customerId"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	Tags       string `json:"tags"`
}
type UpdateCustomerOperation struct {
	options *sdk.OperationInfo
}

func NewUpdateCustomerOperation() *UpdateCustomerOperation {
	return &UpdateCustomerOperation{
		options: &sdk.OperationInfo{
			Name:        "Update Customer",
			Description: "Update an existing customer.",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"customerId": autoform.NewNumberField().
					SetDisplayName("Customer ID").
					SetDescription("The id of the customer to update.").
					SetRequired(true).
					Build(),
				"firstName": autoform.NewShortTextField().
					SetDisplayName("First name").
					SetDescription("Customer first name.").
					SetRequired(false).
					Build(),
				"lastName": autoform.NewShortTextField().
					SetDisplayName("Last name").
					SetDescription("Customer last name.").
					SetRequired(false).
					Build(),
				"phone": autoform.NewShortTextField().
					SetDisplayName("Phone").
					SetDescription("Customer phone number.").
					SetRequired(false).
					Build(),
				"email": autoform.NewShortTextField().
					SetDisplayName("Email").
					SetDescription("Customer email address.").
					SetRequired(false).
					Build(),
				"tags": autoform.NewLongTextField().
					SetDisplayName("Tags").
					SetDescription("A string of comma-separated tags for filtering and search").
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
func (c *UpdateCustomerOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["token"] == "" {
		return nil, errors.New("missing shopify auth token")
	}
	input := sdk.InputToType[updateCustomerOperationProps](ctx)
	domain := ctx.Auth.Extra["domain"]
	shopName := domain + ".myshopify.com"
	client := getShopifyClient(shopName, ctx.Auth.Extra["token"])
	existingCustomer, err := client.Customer.Get(context.Background(), input.CustomerID, nil)
	if err != nil {
		return nil, err
	}
	if input.FirstName != "" {
		existingCustomer.FirstName = input.FirstName
	}
	if input.LastName != "" {
		existingCustomer.LastName = input.LastName
	}
	if input.Email != "" {
		existingCustomer.Email = input.Email
	}
	if input.Phone != "" {
		existingCustomer.Phone = input.Phone
	}
	if input.Tags != "" {
		existingCustomer.Tags = input.Tags
	}
	updatedCustomer, err := client.Customer.Update(context.Background(), *existingCustomer)
	if err != nil {
		return nil, errors.New("failed to update customer")
	}
	return map[string]interface{}{
		"updated_customer": updatedCustomer,
	}, nil
}
func (c *UpdateCustomerOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}
func (c *UpdateCustomerOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
