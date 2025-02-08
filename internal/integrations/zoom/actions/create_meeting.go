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

package actions

import (
	"encoding/json"

	"github.com/wakflo/extensions/internal/integrations/zoom/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/integration"
)

type createMeetingActionProps struct {
	Topic       string `json:"topic"`
	StartTime   string `json:"start_time"`
	ScheduleFor string `json:"schedule_for"`
	Duration    int    `json:"duration"`
	PreSchedule string `json:"pre_schedule"`
	Password    string `json:"password"`
	Agenda      string `json:"agenda"`
	JoinURL     string `json:"join_url"`
}

type CreateMeetingAction struct{}

func (c *CreateMeetingAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func (c CreateMeetingAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (c CreateMeetingAction) Name() string {
	return "Create Zoom Meeting"
}

func (c CreateMeetingAction) Description() string {
	return "create a new zoom meeting"
}

func (c CreateMeetingAction) Documentation() *integration.OperationDocumentation {
	return &integration.OperationDocumentation{
		Documentation: &createMeetingDocs,
	}
}

func (c CreateMeetingAction) Icon() *string {
	return nil
}

func (c CreateMeetingAction) SampleData() sdkcore.JSON {
	return nil
}

func (c CreateMeetingAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
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
	}
}

func (c CreateMeetingAction) Auth() *integration.Auth {
	return &integration.Auth{
		Inherit: true,
	}
}

func (c CreateMeetingAction) Perform(ctx integration.PerformContext) (sdkcore.JSON, error) {
	input, err := integration.InputToTypeSafely[createMeetingActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

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

	resp, err := shared.ZoomRequest(ctx.Auth.AccessToken, reqURL, meeting)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NewCreateMeetingAction() integration.Action {
	return &CreateMeetingAction{}
}
