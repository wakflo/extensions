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

type sendPrivateChannelMessageProps struct {
	Channel string `json:"channel"`
	Message string `json:"message"`
}

type SendPrivateChannelMessage struct {
	options *sdk.OperationInfo
}

func NewSendPrivateChannelMessageOperation() *SendPrivateChannelMessage {
	getPrivateChannels := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		client := getSlackClient(ctx.Auth.AccessToken)

		publicChannels, err := getChannels(client, "private_channel")
		if err != nil {
			return nil, err
		}

		return ctx.Respond(publicChannels, len(publicChannels))
	}

	return &SendPrivateChannelMessage{
		options: &sdk.OperationInfo{
			Name:        "Send private channel message",
			Description: "Operation that sends a message to specified private channel",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"channel": autoform.NewDynamicField(sdkcore.String).
					SetDisplayName("Private Channel").
					SetDescription("Select private channel where message will be sent").
					SetDynamicOptions(&getPrivateChannels).
					SetDependsOn([]string{"connection"}).
					SetRequired(true).
					Build(),
				"message": sharedLongMessageAutoform,
			},
			SampleOutput: map[string]interface{}{
				"name":       "slack-send-private-channel-message",
				"usage_mode": "operation",
				"message":    "Hello people in the private channel!",
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c SendPrivateChannelMessage) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	client := getSlackClient(ctx.Auth.AccessToken)
	input := sdk.InputToType[sendPrivateChannelMessageProps](ctx)
	message := input.Message
	channelID := input.Channel

	err := sendMessage(client, message, channelID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"name":       "slack-send-private-channel-message",
		"usage_mode": "operation",
	}, nil
}

func (c SendPrivateChannelMessage) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c SendPrivateChannelMessage) GetInfo() *sdk.OperationInfo {
	return c.options
}
