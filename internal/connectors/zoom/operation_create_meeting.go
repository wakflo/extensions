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

package zoom

import (
	"encoding/json"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type createMeetingOperationProps struct {
	Topic       string `json:"topic"`
	StartTime   string `json:"start_time"`
	ScheduleFor string `json:"schedule_for"`
	Duration    int    `json:"duration"`
	PreSchedule string `json:"pre_schedule"`
	Password    string `json:"password"`
	Agenda      string `json:"agenda"`
	JoinURL     string `json:"join_url"`
}

type CreateMeetingOperation struct {
	options *sdk.OperationInfo
}

func NewCreateMeetingOperation() *CreateMeetingOperation {
	return &CreateMeetingOperation{
		options: &sdk.OperationInfo{
			Name:        "Create Zoom Meeting",
			Description: "Create a new Zoom Meeting",
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"topic": autoform.NewShortTextField().
					SetDisplayName("Meeting's topic").
					SetDescription("The meeting topic").
					SetRequired(true).Build(),
				"start_time": autoform.NewShortTextField().
					SetDisplayName("Start Time").
					SetDescription("Meeting start date-time").
					SetRequired(false).Build(),
				"schedule_for": autoform.NewShortTextField().
					SetDisplayName("Schedule for").
					SetDescription("The email address or user ID of the user to schedule a meeting for.").
					SetRequired(false).Build(),
				"duration": autoform.NewNumberField().
					SetDisplayName("Duration (in Minutes)").
					SetDescription("Duration of the meeting").
					SetRequired(false).Build(),
				"pre_schedule": autoform.NewBooleanField().
					SetDisplayName("Pre Schedule").
					SetDescription("Whether the prescheduled meeting was created via the GSuite app.").
					SetRequired(false).Build(),
				"password": autoform.NewShortTextField().
					SetDisplayName("Password").
					SetDescription("The password required to join the meeting. By default, a password can only have a maximum length of 10 characters and only contain alphanumeric characters and the @, -, _, and * characters.").
					SetRequired(false).Build(),
				"agenda": autoform.NewLongTextField().
					SetDisplayName("Agenda").
					SetDescription("The meeting's agenda").
					SetRequired(false).Build(),
				"join_url": autoform.NewLongTextField().
					SetDisplayName("Join URL").
					SetDescription("URL for participants to join the meeting.").
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

func (c *CreateMeetingOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	input := sdk.InputToType[createMeetingOperationProps](ctx)

	data := map[string]interface{}{
		"topic":            input.Topic,
		"agenda":           "My Meeting",
		"default_password": false,
		//nolint:mnd
		"duration":     30,
		"pre_schedule": false,
		"timezone":     "UTC",
		//nolint:mnd
		"type": 2,
		"settings": map[string]interface{}{
			"allow_multiple_devices": true,
			//nolint:mnd
			"approval_type":                  2,
			"audio":                          "telephony",
			"calendar_type":                  1,
			"close_registration":             false,
			"email_notification":             true,
			"host_video":                     true,
			"join_before_host":               false,
			"meeting_authentication":         true,
			"mute_upon_entry":                false,
			"participant_video":              false,
			"private_meeting":                false,
			"registrants_confirmation_email": true,
			"registrants_email_notification": true,
			"registration_type":              1,
			"show_share_button":              true,
			"host_save_video_order":          true,
		},
	}

	if input.Duration != 0 {
		data["duration"] = input.Duration
	}
	if input.PreSchedule != "" {
		data["pre_schedule"] = input.PreSchedule
	}
	if input.StartTime != "" {
		data["start_time"] = input.StartTime
	}
	if input.ScheduleFor != "" {
		data["schedule_for"] = input.ScheduleFor
	}
	if input.Agenda != "" {
		data["agenda"] = input.Agenda
	}
	if input.JoinURL != "" {
		data["join_url"] = input.JoinURL
	}
	if input.Password != "" {
		data["[password]"] = input.Password
	}

	meeting, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	reqURL := "/v2/users/me/meetings"

	resp, err := zoomRequest(ctx.Auth.AccessToken, reqURL, meeting)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *CreateMeetingOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CreateMeetingOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
