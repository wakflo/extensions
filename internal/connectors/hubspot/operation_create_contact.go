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

package hubspot

import (
	"encoding/json"
	"net/http"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
	sdk2 "github.com/wakflo/go-sdk/sdk"
)

type createNewContactProps struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Zip       string `json:"zip"`
}

type CreateContactOperation struct {
	options *sdk.OperationInfo
}

func NewCreateContactOperation() *CreateContactOperation {
	return &CreateContactOperation{
		options: &sdk.OperationInfo{
			Name:        "Create new Contact",
			Description: "Create new contact",
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"firstname": autoform.NewShortTextField().
					SetDisplayName("First Name").
					SetDescription("First Name").
					SetRequired(true).Build(),
				"lastname": autoform.NewShortTextField().
					SetDisplayName("Last Name").
					SetDescription("Last name").
					SetRequired(true).Build(),
				"email": autoform.NewShortTextField().
					SetDisplayName("Email").
					SetDescription("Email").
					SetRequired(false).Build(),
				"zip": autoform.NewShortTextField().
					SetDisplayName("Zip").
					SetDescription("Zip Code").
					SetRequired(false).Build(),
			},
			SampleOutput: map[string]interface{}{},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
			RequireAuth: true,
		},
	}
}

func (c *CreateContactOperation) Run(ctx *sdk.RunContext) (sdk2.JSON, error) {
	input := sdk.InputToType[createNewContactProps](ctx)

	contact := contactRequest{
		Properties: map[string]interface{}{
			"firstname": input.FirstName,
			"lastname":  input.LastName,
			"email":     input.Email,
			"zip":       input.Zip,
		},
	}

	newContact, err := json.Marshal(contact)
	if err != nil {
		return nil, err
	}

	reqURL := "/crm/v3/objects/contacts"

	resp, err := hubspotClient(reqURL, ctx.Auth.AccessToken, http.MethodPost, newContact)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *CreateContactOperation) Test(ctx *sdk.RunContext) (sdk2.JSON, error) {
	return c.Run(ctx)
}

func (c *CreateContactOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
