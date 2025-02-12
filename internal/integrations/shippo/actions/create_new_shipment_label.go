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
	"errors"
	"fmt"

	"github.com/wakflo/extensions/internal/integrations/shippo/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type createShipmentLabelActionProps struct {
	Rate string `json:"rate"`
}

type CreateShipmentLabelAction struct{}

func (c *CreateShipmentLabelAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func (c CreateShipmentLabelAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (c CreateShipmentLabelAction) Name() string {
	return "Create Zoom Meeting Registrant"
}

func (c CreateShipmentLabelAction) Description() string {
	return "Create and submit a user's registration to a meeting."
}

func (c CreateShipmentLabelAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &newShipmentLabelDocs,
	}
}

func (c CreateShipmentLabelAction) Icon() *string {
	return nil
}

func (c CreateShipmentLabelAction) SampleData() sdkcore.JSON {
	return nil
}

func (c CreateShipmentLabelAction) Properties() map[string]*sdkcore.AutoFormSchema {
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

func (c CreateShipmentLabelAction) Auth() *sdk.Auth {
	return &sdk.Auth{
		Inherit: true,
	}
}

func (c CreateShipmentLabelAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing shippo api key")
	}

	input, err := sdk.InputToTypeSafely[createShipmentLabelActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}
	endpoint := "/transactions"

	shipmentData := map[string]interface{}{
		"rate":            input.Rate,
		"async":           false,
		"label_file_type": "PDF",
	}

	response, err := shared.CreateAShipment(endpoint, ctx.Auth.Extra["api-key"], shipmentData)
	if err != nil {
		return nil, fmt.Errorf("error creating shipment:  %v", err)
	}

	return response, nil
}

func NewCreateShipmentLabelAction() sdk.Action {
	return &CreateShipmentLabelAction{}
}
