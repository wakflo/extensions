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

type addSubscriberToTagOperationProps struct {
	ListID string `json:"list-id"`
	Tags   string `json:"tag-names"`
	Email  string `json:"email"`
}

type AddSubscriberToTagOperation struct {
	options *sdk.OperationInfo
}

func NewAddSubscriberToTagOperation() sdk.IOperation {
	return &AddSubscriberToTagOperation{
		options: &sdk.OperationInfo{
			Name:        "Add Subscriber to a tag",
			Description: "Adds a subscriber to a tag. This will fail if the user is not subscribed to the audience.",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"list-id": getListInput(),
				"tag-names": autoform.NewLongTextField().
					SetDisplayName(" Tag Name").
					SetDescription("Tag name to add to the subscriber").
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

func (c *AddSubscriberToTagOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing mailchimp auth token")
	}
	accessToken := ctx.Auth.AccessToken

	input := sdk.InputToType[addSubscriberToTagOperationProps](ctx)
	dc, err := getMailChimpServerPrefix(accessToken)
	if err != nil {
		log.Fatalf("Error getting MailChimp server prefix: %v", err)
	}

	tags := processTagNamesInput(input.Tags)

	err = modifySubscriberTags(accessToken, dc, input.ListID, input.Email, tags, "active")
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"status": "Tag added!",
	}, nil
}

func (c *AddSubscriberToTagOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *AddSubscriberToTagOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
