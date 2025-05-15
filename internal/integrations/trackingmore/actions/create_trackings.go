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
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/trackingmore/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type createTrackingsActionProps struct {
	TrackingNumber string                   `json:"tracking_number"`
	CourierCode    string                   `json:"courier_code"`
	Body           []map[string]interface{} `json:"body"`
}

type CreateTrackingsAction struct{}

func (c *CreateTrackingsAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_batch_tracking",
		DisplayName:   "Create batch tracking",
		Description:   "Create batch tracking",
		Type:          core.ActionTypeAction,
		Documentation: createTrackingDocs,
		SampleOutput:  nil,
		Settings:      core.ActionSettings{},
	}
}

func (c *CreateTrackingsAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_batch_tracking", "Create batch tracking")

	form.SelectField("courier_code", "Courier Code").
		Required(true).
		AddOptions(shared.CourierCodes...).
		Placeholder("Courier code").
		HelpText("Courier code")

	form.TextField("tracking_number", "Tracking Number").
		Placeholder("Tracking number of a package.").
		Required(true).
		HelpText("Tracking number of a package.")

	schema := form.Build()
	return schema
}

func (c *CreateTrackingsAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[createTrackingsActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	applicationKey := authCtx.Extra["key"]

	//  payload := []map[string]interface{}{
	//	{
	//		"courier_code":    "usps",
	//		"tracking_number": "9269990312443844954410",
	//	},
	//	{
	//		"courier_code":    "fedex",
	//		"tracking_number": "608285157867",
	//	},
	//}

	response, err := shared.CreateBatchTracking(applicationKey, input.Body)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (c *CreateTrackingsAction) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Inherit: true,
	}
}

func NewCreateTrackingsAction() sdk.Action {
	return &CreateTrackingsAction{}
}
