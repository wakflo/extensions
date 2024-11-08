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
	"fmt"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type createMeetingRegistrantOperationProps struct {
	MeetingID string `json:"meeting_id"`
	MeetingRegistrant
}

type CreateMeetingRegistrantOperation struct {
	options *sdk.OperationInfo
}

func NewCreateMeetingRegistrantOperation() *CreateMeetingRegistrantOperation {
	return &CreateMeetingRegistrantOperation{
		options: &sdk.OperationInfo{
			Name:        "Create Zoom Meeting Registrant",
			Description: "Create and submit a user's registration to a meeting.",
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
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

func (c *CreateMeetingRegistrantOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	input := sdk.InputToType[createMeetingRegistrantOperationProps](ctx)

	meetingRegistrant, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	reqURL := fmt.Sprintf("%s/v2/meetings/%s/registrants", baseURL, input.MeetingID)

	resp, err := zoomRequest(ctx.Auth.AccessToken, reqURL, meetingRegistrant)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *CreateMeetingRegistrantOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CreateMeetingRegistrantOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
