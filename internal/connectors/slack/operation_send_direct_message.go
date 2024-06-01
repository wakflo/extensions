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

package slack

import (
	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type SendDirectMessageOperationProps struct {
	User    string `json:"user"`
	Message string `json:"message"`
}

type SendDirectMessageOperation struct {
	options *sdk.OperationInfo
}

func NewSendDirectMessageOperation() *SendDirectMessageOperation {
	getSlackUsers := func(ctx *sdkcore.DynamicOptionsContext) (interface{}, error) {
		client := getSlackClient(ctx.Auth.AccessToken)

		users, err := getUsers(client)
		if err != nil {
			return nil, err
		}

		return users, nil
	}

	return &SendDirectMessageOperation{
		options: &sdk.OperationInfo{
			Name:        "Send direct message",
			Description: "Operation that sends a direct slack message to a specified user",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"user": autoform.NewDynamicField(sdkcore.String).
					SetDisplayName("User").
					SetDescription("Select user to send message").
					SetDynamicOptions(&getSlackUsers).
					SetDependsOn([]string{"connection"}).
					SetRequired(true).
					Build(),
				"message": sharedLongMessageAutoform,
			},
			SampleOutput: map[string]interface{}{
				"name":       "slack-send-direct-message",
				"usage_mode": "operation",
				"message":    "Hey you, did you know Australia is wider than the moon?",
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c SendDirectMessageOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	input := sdk.InputToType[SendDirectMessageOperationProps](ctx)
	message := input.Message
	userID := input.User

	client := getSlackClient(ctx.Auth.AccessToken)

	err := sendMessage(client, message, userID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"name":       "slack-send-direct-message",
		"usage_mode": "operation",
		"message":    message,
	}, nil
}

func (c SendDirectMessageOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c SendDirectMessageOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
