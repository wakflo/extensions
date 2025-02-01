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
	"github.com/wakflo/extensions/internal/integrations/trackingmore/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/integration"
)

type createTrackingsActionProps struct {
	TrackingNumber string                   `json:"tracking_number"`
	CourierCode    string                   `json:"courier_code"`
	Body           []map[string]interface{} `json:"body"`
}

type CreateTrackingsAction struct{}

func (c CreateTrackingsAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (c CreateTrackingsAction) Name() string {
	return "Create batch tracking"
}

func (c CreateTrackingsAction) Description() string {
	return "Create batch tracking"
}

func (c CreateTrackingsAction) Documentation() *integration.OperationDocumentation {
	return &integration.OperationDocumentation{
		Documentation: &createTrackingDocs,
	}
}

func (c CreateTrackingsAction) Icon() *string {
	return nil
}

func (c CreateTrackingsAction) SampleData() (sdkcore.JSON, error) {
	return nil, nil
}

func (c CreateTrackingsAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"courier_code": autoform.NewSelectField().
			SetDisplayName("Courier Code").
			SetDescription("Courier code").
			SetOptions(shared.CourierCodes).
			SetRequired(true).
			Build(),
		"tracking_number": autoform.NewShortTextField().
			SetDisplayName("Tracking Number").
			SetDescription("Tracking number of a package.").
			SetRequired(true).
			Build(),
	}
}

func (c CreateTrackingsAction) Auth() *integration.Auth {
	return &integration.Auth{
		Inherit: true,
	}
}

func (c CreateTrackingsAction) Perform(ctx integration.PerformContext) (sdkcore.JSON, error) {
	input, err := integration.InputToTypeSafely[createTrackingsActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	applicationKey := ctx.Auth.Extra["key"]

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

func NewCreateTrackingsAction() integration.Action {
	return &CreateTrackingsAction{}
}
