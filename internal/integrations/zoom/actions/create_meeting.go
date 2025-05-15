package actions

import (
	"encoding/json"
	"errors"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/zoom/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
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

func (a *CreateMeetingAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		DisplayName:   "Create Meeting",
		Description:   "Create a new meeting with specified details, including title, start and end dates, duration, and attendees. This action allows you to schedule meetings with ease, streamlining your workflow and reducing manual errors.",
		HelpText:      "Create a new meeting with specified details, including title, start and end dates, duration, and attendees. This action allows you to schedule meetings with ease, streamlining your workflow and reducing manual errors.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: createMeetingDocs,
		SampleOutput:  nil,
		Tags:          nil,
	}
}

func (a *CreateMeetingAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create-meeting", "Creat Meeting")

	form.TextField("topic", "Topic").
		Placeholder("Meeting's topic").
		HelpText("The meeting topic").
		Required(true)

	form.TextField("start_time", "Start Time").
		Placeholder("Start Time").
		HelpText("Meeting start date-time").
		Required(false)

	form.TextField("schedule_for", "Schedule For").
		Placeholder("Schedule For").
		HelpText("The email address or user ID of the user to schedule a meeting for.").
		Required(false)

	form.NumberField("duration", "Duration (in Minutes)").
		Placeholder("Duration (in Minutes)").
		HelpText("Duration of the meeting").
		Required(false)

	form.CheckboxField("pre_schedule", "Pre Schedule").
		Placeholder("Pre Schedule").
		HelpText("Whether the prescheduled meeting was created via the GSuite app.").
		Required(false)

	form.TextField("password", "Password").
		Placeholder("Password").
		HelpText("The password required to join the meeting. By default, a password can only have a maximum length of 10 characters and only contain alphanumeric characters and the @, -, _, and * characters.").
		Required(false)

	form.TextField("agenda", "Agenda").
		Placeholder("Agenda").
		HelpText("The meeting's agenda").
		Required(false)

	form.TextField("join_url", "Join URL").
		Placeholder("Join URL").
		HelpText("URL for participants to join the meeting.").
		Required(false)

	return form.Build()
}

func (a *CreateMeetingAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func (a *CreateMeetingAction) Icon() *string {
	return nil
}

func (a *CreateMeetingAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createMeetingActionProps](ctx)
	if err != nil {
		return nil, err
	}

	tokenSource := ctx.Auth().Token
	if tokenSource == nil {
		return nil, errors.New("missing authentication token")
	}
	token := tokenSource.AccessToken

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

	resp, err := shared.ZoomRequest(token, reqURL, meeting)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (a *CreateMeetingAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreateMeetingAction() sdk.Action {
	return &CreateMeetingAction{}
}
