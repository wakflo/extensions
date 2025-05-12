package actions

import (
	"encoding/json"
	"fmt"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/zoom/shared"
	"github.com/wakflo/go-sdk/autoform"
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
		DisplayName: "Create Meeting Registration",
		Description: "Crea",
	}
}

func (a *CreateMeetingRegistrationAction) Properties() *smartform.FormSchema {
	// TODO implement me
	panic("implement me")
}

func (a *CreateMeetingRegistrationAction) Auth() *sdkcore.AuthMetadata {
	// TODO implement me
	panic("implement me")
}

func (a *CreateMeetingRegistrationAction) Name() string {
	return "Create Meeting Registration"
}

func (a *CreateMeetingRegistrationAction) Description() string {
	return "Create Meeting Registration: Automatically generates meeting registrations for attendees, including details such as name, email, and RSVP status. This integration action streamlines the process of tracking attendee information and reduces manual errors."
}

func (a *CreateMeetingRegistrationAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *CreateMeetingRegistrationAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createMeetingRegistrationDocs,
	}
}

func (a *CreateMeetingRegistrationAction) Icon() *string {
	return nil
}

func (a *CreateMeetingRegistrationAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"meeting_id": autoform.NewShortTextField().
			SetDisplayName("Meeting ID").
			SetDescription("The meeting's ID").
			SetRequired(true).Build(),
		"first_name": autoform.NewShortTextField().
			SetDisplayName("First Name").
			SetDescription("First name").
			SetRequired(true).Build(),
		"last_name": autoform.NewShortTextField().
			SetDisplayName("Last Name").
			SetDescription("Last name").
			SetRequired(false).Build(),
		"email": autoform.NewShortTextField().
			SetDisplayName("Email").
			SetDescription("The registrant's email address.").
			SetRequired(true).Build(),
		"address": autoform.NewShortTextField().
			SetDisplayName("Address").
			SetDescription("The registrant's address.").
			SetRequired(false).Build(),
		"city": autoform.NewShortTextField().
			SetDisplayName("City").
			SetDescription("The registrant's city.").
			SetRequired(false).Build(),
		"state": autoform.NewShortTextField().
			SetDisplayName("State").
			SetDescription("The registrant's state or province.").
			SetRequired(false).Build(),
		"zip": autoform.NewShortTextField().
			SetDisplayName("Zip").
			SetDescription("The registrant's zip or postal code.").
			SetRequired(false).Build(),
		"country": autoform.NewShortTextField().
			SetDisplayName("Country").
			SetDescription("The registrant's two-letter country code.").
			SetRequired(false).Build(),
		"phone": autoform.NewShortTextField().
			SetDisplayName("Phone").
			SetDescription("The registrant's phone number.").
			SetRequired(false).Build(),
		"comments": autoform.NewLongTextField().
			SetDisplayName("Phone").
			SetDescription("The registrant's questions and comments.").
			SetRequired(false).Build(),
		"industry": autoform.NewShortTextField().
			SetDisplayName("Industry").
			SetDescription("The registrant's industry").
			SetRequired(false).Build(),
		"job_title": autoform.NewShortTextField().
			SetDisplayName("Job Title").
			SetDescription("The registrant's job title").
			SetRequired(false).Build(),
		"no_of_employees": autoform.NewShortTextField().
			SetDisplayName("Number of Employees").
			SetDescription("The registrant's number of employees.").
			SetRequired(false).Build(),
		"org": autoform.NewShortTextField().
			SetDisplayName("Organization").
			SetDescription("The registrant's organization.").
			SetRequired(false).Build(),
		"purchasing_time_frame": autoform.NewShortTextField().
			SetDisplayName("Purchasing time frame").
			SetDescription("The registrant's purchasing time frame.").
			SetRequired(false).Build(),
		"role_in_purchase_process": autoform.NewShortTextField().
			SetDisplayName("Role in purchase process").
			SetDescription("The registrant's role in purchase process").
			SetRequired(false).Build(),
	}
}

func (a *CreateMeetingRegistrationAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createMeetingRegistrationActionProps](ctx)
	if err != nil {
		return nil, err
	}

	meetingRegistrant, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	reqURL := fmt.Sprintf("%s/v2/meetings/%s/registrants", shared.ZoomBaseURL, input.MeetingID)

	resp, err := shared.ZoomRequest(ctx.Auth().AccessToken, reqURL, meetingRegistrant)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (a *CreateMeetingRegistrationAction) Auth() *sdk.Auth {
	return nil
}

func (a *CreateMeetingRegistrationAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *CreateMeetingRegistrationAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreateMeetingRegistrationAction() sdk.Action {
	return &CreateMeetingRegistrationAction{}
}
