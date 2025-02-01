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
	"fmt"

	"github.com/wakflo/extensions/internal/integrations/zoom/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/integration"
)

type createMeetingRegistrantActionProps struct {
	MeetingID string `json:"meeting_id"`
	shared.MeetingRegistrant
}

type CreateMeetingRegistrantAction struct{}

func (c CreateMeetingRegistrantAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (c CreateMeetingRegistrantAction) Name() string {
	return "Create Zoom Meeting Registrant"
}

func (c CreateMeetingRegistrantAction) Description() string {
	return "Create and submit a user's registration to a meeting."
}

func (c CreateMeetingRegistrantAction) Documentation() *integration.OperationDocumentation {
	return &integration.OperationDocumentation{
		Documentation: &createMeetingRegistrationDocs,
	}
}

func (c CreateMeetingRegistrantAction) Icon() *string {
	return nil
}

func (c CreateMeetingRegistrantAction) SampleData() (sdkcore.JSON, error) {
	return nil, nil
}

func (c CreateMeetingRegistrantAction) Properties() map[string]*sdkcore.AutoFormSchema {
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

func (c CreateMeetingRegistrantAction) Auth() *integration.Auth {
	return &integration.Auth{
		Inherit: true,
	}
}

func (c CreateMeetingRegistrantAction) Perform(ctx integration.PerformContext) (sdkcore.JSON, error) {
	input, err := integration.InputToTypeSafely[createMeetingRegistrantActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	meetingRegistrant, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	reqURL := fmt.Sprintf("%s/v2/meetings/%s/registrants", shared.ZoomBaseURL, input.MeetingID)

	resp, err := shared.ZoomRequest(ctx.Auth.AccessToken, reqURL, meetingRegistrant)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NewCreateMeetingRegistrantAction() integration.Action {
	return &CreateMeetingRegistrantAction{}
}
