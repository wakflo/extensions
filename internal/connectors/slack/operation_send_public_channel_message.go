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

type SendPublicChannelMessage struct {
	options *sdk.OperationInfo
}

func NewSendPublicChannelMessageOperation() *SendPublicChannelMessage {
	getPublicChannels := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		client := getSlackClient(ctx.Auth.AccessToken)

		publicChannels, err := getChannels(client, "public_channel")
		if err != nil {
			return nil, err
		}

		return ctx.Respond(publicChannels, len(publicChannels))
	}

	return &SendPublicChannelMessage{
		options: &sdk.OperationInfo{
			Name:        "Send public channel message",
			Description: "Action that sends a message to specified public channel",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"channel": autoform.NewDynamicField(sdkcore.String).
					SetDisplayName("Public Channel").
					SetDescription("Select public channel where message will be sent").
					SetDynamicOptions(&getPublicChannels).
					SetDependsOn([]string{"connection"}).
					SetRequired(true).
					Build(),
				"message": sharedLongMessageAutoform,
			},
			SampleOutput: map[string]interface{}{
				"name":       "slack-send-public-channel-message",
				"usage_mode": "operation",
				"message":    "Hello people in the public channel!",
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c SendPublicChannelMessage) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	client := getSlackClient(ctx.Auth.AccessToken)
	input := sdk.InputToType[sendPrivateChannelMessageProps](ctx)

	message := input.Message
	channelID := input.Channel

	err := sendMessage(client, message, channelID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"name":       "slack-send-public-channel-message",
		"usage_mode": "operation",
		"message":    message,
	}, nil
}

func (c SendPublicChannelMessage) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c SendPublicChannelMessage) GetInfo() *sdk.OperationInfo {
	return c.options
}
