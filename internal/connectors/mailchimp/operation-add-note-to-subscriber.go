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

type addNoteToSubscriberOperationProps struct {
	ListID string `json:"list-id"`
	Email  string `json:"email"`
	Note   string `json:"note"`
}

type AddNoteToSubscriberOperation struct {
	options *sdk.OperationInfo
}

func NewAddNoteToSubscriberOperation() sdk.IOperation {
	return &AddNoteToSubscriberOperation{
		options: &sdk.OperationInfo{
			Name:        "Add note to Subscriber",
			Description: "Add note to a Subscriber",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"list-id": getListInput(),
				"note": autoform.NewLongTextField().
					SetDisplayName(" Note").
					SetDescription("Note to add to the subscriber").
					SetRequired(true).
					Build(),
				"email": autoform.NewShortTextField().
					SetDisplayName("Email").
					SetDescription("Email of the subscriber").
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

func (c *AddNoteToSubscriberOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing mailchimp auth token")
	}
	accessToken := ctx.Auth.AccessToken

	input := sdk.InputToType[addNoteToSubscriberOperationProps](ctx)

	dc, err := getMailChimpServerPrefix(accessToken)
	if err != nil {
		log.Fatalf("Error getting MailChimp server prefix: %v", err)
	}

	err = addMemberNote(accessToken, dc, input.ListID, input.Email, input.Note)
	if err != nil {
		return nil, err
	}

	return sdk.JSON(map[string]interface{}{
		"result": "note added!",
	}), err
}

func (c *AddNoteToSubscriberOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *AddNoteToSubscriberOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
