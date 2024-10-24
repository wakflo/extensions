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

package mailchimp

import (
	"errors"
	"log"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type addMemberToListOperationProps struct {
	ListID    string `json:"list-id"`
	To        string `json:"to"`
	FirstName string `json:"first-name"`
	LastName  string `json:"last-name"`
	Email     string `json:"email"`
	Status    string `json:"status"`
}

type AddMemberToListOperation struct {
	options *sdk.OperationInfo
}

func NewAddMemberToListOperation() sdk.IOperation {
	return &AddMemberToListOperation{
		options: &sdk.OperationInfo{
			Name:        "Add Member to an Audience (List)",
			Description: "Add a member to an existing Mailchimp audience (list)",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"list-id": autoform.NewShortTextField().
					SetDisplayName(" List ID").
					SetDescription("").
					SetRequired(true).
					Build(),
				"first-name": autoform.NewShortTextField().
					SetDisplayName(" First Name").
					SetDescription("First name of the new contact").
					SetRequired(true).
					Build(),
				"last-name": autoform.NewShortTextField().
					SetDisplayName(" Last Name").
					SetDescription("Last name of the new contact").
					SetRequired(true).
					Build(),
				"email": autoform.NewShortTextField().
					SetDisplayName("Email").
					SetDescription("Email of the new contact").
					SetRequired(true).
					Build(),
				"status": autoform.NewSelectField().
					SetDisplayName("Status").
					SetOptions(mailchimpStatusType).
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

func (c *AddMemberToListOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing mailchimp auth token")
	}
	accessToken := ctx.Auth.AccessToken

	input := sdk.InputToType[addMemberToListOperationProps](ctx)
	dc, err := getMailChimpServerPrefix(accessToken)
	if err != nil {
		log.Fatalf("Error getting MailChimp server prefix: %v", err)
	}

	err = addContactToList(accessToken, dc, input.ListID, input.Email, input.FirstName, input.Status, input.LastName)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success": "Contact Added!",
	}, err
}

func (c *AddMemberToListOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *AddMemberToListOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
