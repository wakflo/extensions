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

type SlackGenericRespose struct {
	Ok        bool   `json:"ok"`
	TimeStamp string `json:"ts"`
}

type SlackMetadataResponse struct {
	NextCursor string `json:"next_cursor"`
}

type SlackUserResponse struct {
	ID       string `json:"id"`
	TeamID   string `json:"team_id"`
	Name     string `json:"name"`
	Deleted  bool   `json:"deleted"`
	RealName string `json:"real_name"`
	IsBot    bool   `json:"is_bot"`
}

type SlackListUsersResponse struct {
	Ok        bool                  `json:"ok"`
	Error     string                `json:"error"` // Has value when Ok is false
	Members   []SlackUserResponse   `json:"members"`
	Timestamp int64                 `json:"cache_ts"` // Confirm if timestamp (seconds) is int64
	Metadata  SlackMetadataResponse `json:"response_metadata"`
}

type SlackUser struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	RealName string `json:"real_name"`
}

type SlackChannelResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	IsChannel  bool   `json:"is_channel"`
	IsGroup    bool   `json:"is_group"`
	IsMpim     bool   `json:"is_mpim"`
	IsIM       bool   `json:"is_im"`
	IsPrivate  bool   `json:"is_private"`
	IsMember   bool   `json:"is_member"`
	IsArchived bool   `json:"is_archived"`
}

type SlackChannelsListResponse struct {
	Ok       bool                   `json:"ok"`
	Error    string                 `json:"error"`
	Channels []SlackChannelResponse `json:"channels"`
	Metadata SlackMetadataResponse  `json:"response_metadata"`
}

type SlackChannel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SlackMessage struct {
	Text      string `json:"text"`
	BotID     string `json:"bot_id"`
	Timestamp string `json:"ts"`
}

type SlackPostMessageResponse struct {
	Ok        bool         `json:"ok"`
	Error     string       `json:"error"`
	Channel   string       `json:"channel"`
	Timestamp string       `json:"ts"`
	Message   SlackMessage `json:"message"`
}
