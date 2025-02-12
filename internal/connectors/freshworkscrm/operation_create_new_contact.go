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

package freshworkscrm

import (
	"errors"
	"fmt"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
	sdk2 "github.com/wakflo/go-sdk/sdk"
)

type createNewContactOperationProps struct {
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	MobileNumber string  `json:"mobile_number"`
	Email        *string `json:"email"`
}

type CreateNewContactOperation struct {
	options *sdk.OperationInfo
}

func NewCreateNewContactOperation() *CreateNewContactOperation {
	return &CreateNewContactOperation{
		options: &sdk.OperationInfo{
			Name:        "Create A New Contact",
			Description: "creates a new contact",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"first_name": autoform.NewShortTextField().
					SetDisplayName("First Name").
					SetDescription("Contact's first name").
					SetRequired(true).
					Build(),
				"last_name": autoform.NewShortTextField().
					SetDisplayName("Last Name").
					SetDescription("Contact's last name").
					SetRequired(true).
					Build(),
				"email": autoform.NewShortTextField().
					SetDisplayName("Email").
					SetDescription("Contact's email").
					SetRequired(true).
					Build(),
				"mobile_number": autoform.NewShortTextField().
					SetDisplayName("Mobile Number").
					SetDescription("Contact's mobile number").
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

func (c *CreateNewContactOperation) Run(ctx *sdk.RunContext) (sdk2.JSON, error) {
	if ctx.Auth.Extra["api-key"] == "" || ctx.Auth.Extra["domain"] == "" {
		return nil, errors.New("missing freshworks auth parameters")
	}

	input := sdk.InputToType[createNewContactOperationProps](ctx)

	domain := ctx.Auth.Extra["domain"]
	freshworksDomain := "https://" + domain + ".myfreshworks.com"

	contactData := map[string]interface{}{
		"contact": map[string]interface{}{
			"first_name":    input.FirstName,
			"last_name":     input.LastName,
			"email":         input.Email,
			"mobile_number": input.MobileNumber,
		},
	}

	response, err := CreateContact(freshworksDomain, ctx.Auth.Extra["api-key"], contactData)
	if err != nil {
		return nil, fmt.Errorf("error creating contact:  %v", err)
	}

	return sdk2.JSON(map[string]interface{}{
		"Status": response,
	}), nil
}

func (c *CreateNewContactOperation) Test(ctx *sdk.RunContext) (sdk2.JSON, error) {
	return c.Run(ctx)
}

func (c *CreateNewContactOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
