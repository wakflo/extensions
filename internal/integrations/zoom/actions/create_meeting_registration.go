package actions

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/zoom/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type createMeetingRegistrationActionProps struct {
	MeetingID string `json:"meeting_id"`
	shared.MeetingRegistrant
}

type CreateMeetingRegistrationAction struct{}

func (a *CreateMeetingRegistrationAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		DisplayName:   "Create Meeting Registration",
		Description:   "Create Meeting Registration: Automatically generates meeting registrations for attendees, including details such as name, email, and RSVP status. This integration action streamlines the process of tracking attendee information and reduces manual errors.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: createMeetingRegistrationDocs,
		SampleOutput: map[string]any{
			"message": "hello world",
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *CreateMeetingRegistrationAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_meeting_registration", "Create Meeting Registration")

	form.TextField("meeting_id", "Meeting ID").
		Placeholder("Enter a value for Meeting ID.").
		Required(true).
		HelpText("The meeting ID")

	form.TextField("first_name", "First Name").
		Placeholder("Enter a value for First Name.").
		Required(true).
		HelpText("The registrant's first name")

	form.TextField("last_name", "Last Name").
		Placeholder("Enter a value for Last Name.").
		Required(false).
		HelpText("The registrant's last name")

	form.TextField("email", "Email").
		Placeholder("Enter a value for Email.").
		Required(true).
		HelpText("The registrant's email address")

	form.TextField("address", "Address").
		Placeholder("Enter a value for Address.").
		Required(false).
		HelpText("The registrant's address")

	form.TextField("city", "City").
		Placeholder("Enter a value for City.").
		Required(false).
		HelpText("The registrant's city")

	form.TextField("state", "State").
		Placeholder("Enter a value for State.").
		Required(false).
		HelpText("The registrant's state")

	form.TextField("zip", "Zip").
		Placeholder("Enter a value for Zip.").
		Required(false).
		HelpText("The registrant's zip code")

	form.TextField("country", "Country").
		Placeholder("Enter a value for Country.").
		Required(false).
		HelpText("The registrant's country")

	form.TextField("phone", "Phone").
		Placeholder("Enter a value for Phone.").
		Required(false).
		HelpText("The registrant's phone number")

	form.TextareaField("comments", "Comments").
		Placeholder("The registrant's questions and comments.").
		Required(false).
		HelpText("The registrant's questions and comments")

	form.TextField("industry", "Industry").
		Placeholder("Enter a value for Industry.").
		Required(false).
		HelpText("The registrant's industry")

	form.TextField("job_title", "Job Title").
		Placeholder("Enter a value for Job Title.").
		Required(false).
		HelpText("The registrant's job title")

	form.TextField("no_of_employees", "Number of Employees").
		Placeholder("Enter a value for Number of Employees.").
		Required(false).
		HelpText("The registrant's number of employees")

	form.TextField("org", "Organization").
		Placeholder("Enter a value for Organization.").
		Required(false).
		HelpText("The registrant's organization")

	form.TextField("purchasing_time_frame", "Purchasing Time Frame").
		Placeholder("Enter a value for Purchasing Time Frame.").
		Required(false).
		HelpText("The registrant's purchasing time frame")

	form.TextField("role_in_purchase_process", "Role in Purchase Process").
		Placeholder("Enter a value for Role in Purchase Process.").
		Required(false).
		HelpText("The registrant's role in the purchase process")

	schema := form.Build()

	return schema
}

func (a *CreateMeetingRegistrationAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func (a *CreateMeetingRegistrationAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createMeetingRegistrationActionProps](ctx)
	if err != nil {
		return nil, err
	}

	tokenSource := ctx.Auth().Token
	if tokenSource == nil {
		return nil, errors.New("missing authentication token")
	}
	token := tokenSource.AccessToken

	meetingRegistrant, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	reqURL := fmt.Sprintf("%s/v2/meetings/%s/registrants", shared.ZoomBaseURL, input.MeetingID)

	resp, err := shared.ZoomRequest(token, reqURL, meetingRegistrant)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NewCreateMeetingRegistrationAction() sdk.Action {
	return &CreateMeetingRegistrationAction{}
}
