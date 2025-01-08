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
)

type updateContactOperationProps struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	MobileNumber  string `json:"mobile_number"`
	Email         string `json:"email"`
	ContactViewID string `json:"contact_view_id"`
	ContactID     string `json:"contact_id"`
}

type UpdateContactOperation struct {
	options *sdk.OperationInfo
}

func NewUpdateContactOperation() *UpdateContactOperation {
	return &UpdateContactOperation{
		options: &sdk.OperationInfo{
			Name:        "Update Contact",
			Description: "Updates a contact",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"first_name": autoform.NewShortTextField().
					SetDisplayName("First Name").
					SetDescription("Contact's first name").
					SetRequired(false).
					Build(),
				"last_name": autoform.NewShortTextField().
					SetDisplayName("Last Name").
					SetDescription("Contact's last name").
					SetRequired(false).
					Build(),
				"email": autoform.NewShortTextField().
					SetDisplayName("Email").
					SetDescription("Contact's email").
					SetRequired(false).
					Build(),
				"mobile_number": autoform.NewShortTextField().
					SetDisplayName("Mobile Number").
					SetDescription("Contact's mobile number").
					SetRequired(false).
					Build(),
				"contact_view_id": getContactViewInput(),
				"contact_id":      getContactsInput(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *UpdateContactOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["api-key"] == "" || ctx.Auth.Extra["domain"] == "" {
		return nil, errors.New("missing freshworks auth parameters")
	}

	input := sdk.InputToType[updateContactOperationProps](ctx)

	domain := ctx.Auth.Extra["domain"]
	freshworksDomain := "https://" + domain + ".myfreshworks.com"

	contact := make(map[string]interface{})

	updateField(contact, "first_name", input.FirstName)
	updateField(contact, "last_name", input.LastName)
	updateField(contact, "email", input.Email)
	updateField(contact, "mobile_number", input.MobileNumber)

	contactData := map[string]interface{}{
		"contact": contact,
	}

	response, err := updateContact(freshworksDomain, ctx.Auth.Extra["api-key"], input.ContactID, contactData)
	if err != nil {
		return nil, fmt.Errorf("error updating contact:  %v", err)
	}

	return sdk.JSON(map[string]interface{}{
		"Status": response,
	}), nil
}

func (c *UpdateContactOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *UpdateContactOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
