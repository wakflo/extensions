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

const SlackApiUrl = "https://slack.com/api"

var sharedLongMessageAutoform = autoform.
	NewLongTextField().
	SetDisplayName("Message").
	SetDescription("Message that will be sent").
	SetMinLength(1).
	SetMaxLength(4000). // https://api.slack.com/apis/rate-limits#rtm-posting-messages
	SetRequired(false).
	Build()

var (
	authUrl  = fmt.Sprintf("https://slack.com/oauth/v2/authorize")
	tokenUrl = fmt.Sprintf("https://slack.com/api/oauth.v2.access")
)

var sharedAuth = autoform.NewOAuthField(authUrl, &tokenUrl, []string{
	"channels:read",
	"chat:write",
	"chat:write.public",
	"groups:read",
	"users:read",
}).Build()

func getSlackClient(accessToken string) fastshot.ClientHttpMethods {
	return fastshot.NewClient(SlackApiUrl).
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
		return nil, errors.New(fmt.Sprintf("Bad HTTP request - error: %s", resp.StatusText()))
	}

	respBytes, err := io.ReadAll(resp.RawBody())

	var respJson SlackListUsersResponse
	err = json.Unmarshal(respBytes, &respJson)
	if err != nil {
		return nil, err
	}

	if !respJson.Ok {
		return nil, errors.New(fmt.Sprintf("Bad Slack request, error: '%s'", respJson.Error))
	}

	var users []SlackUser
	for _, member := range respJson.Members {
		// SlackBot needs to be filtered out with id since IsBot flag
		if !member.IsBot && !member.Deleted && member.Id != "USLACKBOT" {
			users = append(users, SlackUser{Id: member.Id, Name: member.Name, RealName: member.RealName})
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
		return nil, errors.New(fmt.Sprintf("Bad request - error: %s", resp.StatusText()))
	}

	respBytes, err := io.ReadAll(resp.RawBody())

	var respJson SlackChannelsListResponse
	err = json.Unmarshal(respBytes, &respJson)
	if err != nil {
		fmt.Println("Error on JSON marshal")
		return nil, err
	}

	if !respJson.Ok {
		return nil, errors.New(fmt.Sprintf("Bad request, error: '%s'", respJson.Error))
	}

	var channels []SlackChannel
	for _, channel := range respJson.Channels {
		if !channel.IsArchived {
			channels = append(channels, SlackChannel{Id: channel.Id, Name: channel.Name})
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
		return errors.New(fmt.Sprintf("Bad request - error: %s", resp.StatusText()))
	}

	respBytes, err := io.ReadAll(resp.RawBody())
	if err != nil {
		return err
	}

	var respJson SlackPostMessageResponse
	err = json.Unmarshal(respBytes, &respJson)
	if err != nil {
		return err
	}

	if !respJson.Ok {
		return errors.New(fmt.Sprintf("Bad request, error: '%s'", respJson.Error))
	}

	return nil
}
