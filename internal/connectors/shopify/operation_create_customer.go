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

package shopify

import (
	"context"
	"errors"

	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type createCustomerOperationProps struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Tags      string `json:"tags"`
}
type CreateCustomerOperation struct {
	options *sdk.OperationInfo
}

func NewCreateCustomerOperation() *CreateCustomerOperation {
	return &CreateCustomerOperation{
		options: &sdk.OperationInfo{
			Name:        "Create Customer",
			Description: "Create a new customer.",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
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

func (c *CreateCustomerOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["token"] == "" {
		return nil, errors.New("missing shopify auth token")
	}

	input := sdk.InputToType[createCustomerOperationProps](ctx)

	domain := ctx.Auth.Extra["domain"]
	shopName := domain + ".myshopify.com"
	client := getShopifyClient(shopName, ctx.Auth.Extra["token"])

	customer := goshopify.Customer{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Tags:      input.Tags,
		Phone:     input.Phone,
	}
	newCustomer, err := client.Customer.Create(context.Background(), customer)
	if err != nil {
		return nil, err
	}
	if newCustomer == nil {
		return nil, errors.New("customer not created! ")
	}
	return map[string]interface{}{
		"new_customer": newCustomer,
	}, nil
}

func (c *CreateCustomerOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CreateCustomerOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
