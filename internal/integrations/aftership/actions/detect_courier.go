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
	"github.com/aftership/tracking-sdk-go/v5"
	"github.com/aftership/tracking-sdk-go/v5/model"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/integration"
)

type detectCourierActionProps struct {
	TrackingNumber string `json:"tracking_number,omitempty"`
}

type DetectCourierAction struct{}

func (c DetectCourierAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (c DetectCourierAction) Name() string {
	return "Detect Courier"
}

func (c DetectCourierAction) Description() string {
	return "Returns a list of matched couriers based on tracking number format"
}

func (c DetectCourierAction) Documentation() *integration.OperationDocumentation {
	return &integration.OperationDocumentation{
		Documentation: &detectCourierDocs,
	}
}

func (c DetectCourierAction) Icon() *string {
	return nil
}

func (c DetectCourierAction) SampleData() (sdkcore.JSON, error) {
	return nil, nil
}

func (c DetectCourierAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"tracking_number": autoform.NewShortTextField().
			SetDisplayName("Tracking Number").
			SetDescription("tracking number of the shipment").
			SetRequired(true).
			Build(),
	}
}

func (c DetectCourierAction) Auth() *integration.Auth {
	return &integration.Auth{
		Inherit: true,
	}
}

func (c DetectCourierAction) Perform(ctx integration.PerformContext) (sdkcore.JSON, error) {
	input, err := integration.InputToTypeSafely[detectCourierActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	afterShipSdk, err := tracking.New(tracking.WithApiKey(ctx.Auth.Extra["api-key"]))
	if err != nil {
		return nil, err
	}

	result, err := afterShipSdk.Courier.
		DetectCourier().
		BuildBody(model.DetectCourierRequest{
			TrackingNumber: input.TrackingNumber,
		}).
		Execute()
	if err != nil {
		return nil, err
	}

	return result.Couriers, nil
}

func NewDetectCourierAction() integration.Action {
	return &DetectCourierAction{}
}
