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
	"encoding/json"
	"errors"
	"fmt"
	"io"

	fastshot "github.com/opus-domini/fast-shot"

	"github.com/wakflo/go-sdk/autoform"
)

const SlackAPIURL = "https://slack.com/api"

var sharedLongMessageAutoform = autoform.
	NewLongTextField().
	SetDisplayName("Message").
	SetDescription("Message that will be sent").
	SetMinLength(1).
	//nolint:mnd
	SetMaxLength(4000). // https://api.slack.com/apis/rate-limits#rtm-posting-messages
	SetRequired(false).
	Build()

var (
	// #nosec
	authURL = "https://slack.com/oauth/v2/authorize"
	// #nosec
	tokenURL = "https://slack.com/api/oauth.v2.access"
)

var sharedAuth = autoform.NewOAuthField(authURL, &tokenURL, []string{
	"channels:read",
	"chat:write",
	"chat:write.public",
	"groups:read",
	"users:read",
}).Build()

func getSlackClient(accessToken string) fastshot.ClientHttpMethods {
	return fastshot.NewClient(SlackAPIURL).
		Auth().BearerToken(accessToken).
		Build()
}

// Function that requests Slack members to the API.
//
// Bots and SlackBot are ignored.
func getUsers(client fastshot.ClientHttpMethods) ([]SlackUser, error) {
	resp, err := client.GET("/users.list").Send()
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("bad HTTP request - error: %s", resp.StatusText())
	}

	respBytes, err := io.ReadAll(resp.RawBody())
	if err != nil {
		return nil, err
	}

	var respJSON SlackListUsersResponse
	err = json.Unmarshal(respBytes, &respJSON)
	if err != nil {
		return nil, err
	}

	if !respJSON.Ok {
		return nil, errors.New(fmt.Sprintf("Bad Slack request, error: '%s'", respJSON.Error))
	}

	var users []SlackUser
	for _, member := range respJSON.Members {
		// SlackBot needs to be filtered out with id since IsBot flag
		if !member.IsBot && !member.Deleted && member.ID != "USLACKBOT" {
			users = append(users, SlackUser{ID: member.ID, Name: member.Name, RealName: member.RealName})
		}
	}

	return users, nil
}

// Function that requests Slack channels of given type(s) to the API and returns an array of SlackChannel.
// You can request multiple, e.g: "public_channel,private_channel" to get public and private channels.
//
//	Types: public_channel, private_channel, mpim (group direct message), im (direct messages).
//	INFO: private channels are only showed if Bot is invited to the channel
func getChannels(client fastshot.ClientHttpMethods, channelTypes string) ([]SlackChannel, error) {
	payload := map[string]string{
		"types": channelTypes,
	}

	resp, err := client.GET("/conversations.list").
		Query().AddParams(payload).
		Send()
	if err != nil {
		fmt.Println("Error on the response")
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("bad request - error: %s", resp.StatusText())
	}

	respBytes, err := io.ReadAll(resp.RawBody())
	if err != nil {
		return nil, err
	}

	var respJSON SlackChannelsListResponse
	err = json.Unmarshal(respBytes, &respJSON)
	if err != nil {
		fmt.Println("Error on JSON marshal")
		return nil, err
	}

	if !respJSON.Ok {
		return nil, errors.New(fmt.Sprintf("Bad request, error: '%s'", respJSON.Error))
	}

	var channels []SlackChannel
	for _, channel := range respJSON.Channels {
		if !channel.IsArchived {
			channels = append(channels, SlackChannel{ID: channel.ID, Name: channel.Name})
		}
	}

	return channels, nil
}

// Function that sends given message to a given Channel ID.
//
//	Channel ID - can be channel id or a user id.
func sendMessage(client fastshot.ClientHttpMethods, message string, channelID string) error {
	payload := map[string]interface{}{
		"channel": channelID,
		"text":    message,
	}

	resp, err := client.POST("/chat.postMessage").
		Header().AddContentType("application/json").
		Body().AsJSON(payload).
		Send()
	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("bad request - error: %s", resp.StatusText())
	}

	respBytes, err := io.ReadAll(resp.RawBody())
	if err != nil {
		return err
	}

	var respJSON SlackPostMessageResponse
	err = json.Unmarshal(respBytes, &respJSON)
	if err != nil {
		return err
	}

	if !respJSON.Ok {
		return fmt.Errorf("bad request, error: '%s'", respJSON.Error)
	}

	return nil
}
