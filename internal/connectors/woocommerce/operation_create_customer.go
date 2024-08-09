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

package woocommerce

import (
	"errors"

	"github.com/hiscaler/woocommerce-go"
	"github.com/hiscaler/woocommerce-go/entity"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type CreateCustomerOperation struct {
	options *sdk.OperationInfo
}

type createCustomerOperationProps struct {
	Email         string `json:"email,omitempty"`
	FirstName     string `json:"first_name,omitempty"`
	LastName      string `json:"last_name,omitempty"`
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	Phone         string `json:"phone,omitempty"`
	City          string `json:"city,omitempty"`
	Country       string `json:"country,omitempty"`
	State         string `json:"state,omitempty"`
	StreetAddress string `json:"street_address,omitempty"`
	PostalCode    string `json:"post_code,omitempty"`
}

func NewCreateCustomerOperation() *CreateCustomerOperation {
	return &CreateCustomerOperation{
		options: &sdk.OperationInfo{
			Name:        "Create a customer",
			Description: "create a customer",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"first_name": autoform.NewShortTextField().
					SetDisplayName("First Name").
					SetDescription("Enter first name").
					SetRequired(true).
					Build(),
				"last_name": autoform.NewShortTextField().
					SetDisplayName("Last Name").
					SetRequired(true).
					Build(),
				"email": autoform.NewShortTextField().
					SetDisplayName(" Email").
					SetDescription("Enter email address").
					SetRequired(true).
					Build(),
				"username": autoform.NewShortTextField().
					SetDisplayName("Username").
					SetDescription("Enter username").
					SetRequired(true).
					Build(),
				"password": autoform.NewShortTextField().
					SetDisplayName("Password").
					SetDescription("Enter Password").
					SetRequired(true).
					Build(),
				"phone": autoform.NewShortTextField().
					SetDisplayName("Phone").
					SetDescription("Enter Phone number").
					Build(),
				"country": autoform.NewShortTextField().
					SetDisplayName("Country").
					SetDescription("Enter Country").
					SetRequired(true).
					Build(),
				"city": autoform.NewShortTextField().
					SetDisplayName("City").
					SetRequired(true).
					SetDescription("Enter City").
					Build(),
				"state": autoform.NewShortTextField().
					SetDisplayName("State").
					SetRequired(true).
					SetDescription("Enter State").
					Build(),
				"street_address": autoform.NewLongTextField().
					SetDisplayName("Address").
					SetDescription("Enter the street address").
					SetRequired(true).
					Build(),
				"post_code": autoform.NewShortTextField().
					SetDisplayName("Postal Code").
					SetDescription("Enter State").
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
	baseURL := ctx.Auth.Extra["shop-url"]
	consumerKey := ctx.Auth.Extra["consumer-key"]
	consumerSecret := ctx.Auth.Extra["consumer-secret"]

	if baseURL == "" || consumerKey == "" || consumerSecret == "" {
		return nil, errors.New("missing WooCommerce authentication credentials")
	}

	input := sdk.InputToType[createCustomerOperationProps](ctx)

	wooClient := initializeWooCommerceClient(baseURL, consumerKey, consumerSecret)

	billing := entity.Billing{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Address1:  input.StreetAddress,
		City:      input.City,
		State:     input.State,
		Postcode:  input.PostalCode,
		Country:   input.Country,
		Email:     input.Email,
		Phone:     input.Phone,
	}

	// Create a query parameters struct
	params := woocommerce.CreateCustomerRequest{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Username:  input.Username,
		Password:  input.Password,
		Billing:   &billing,
	}

	newCustomer, err := wooClient.Services.Customer.Create(params)
	if err != nil {
		return nil, err
	}

	return newCustomer, nil
}

func (c *CreateCustomerOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CreateCustomerOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
