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
	"github.com/wakflo/extensions/internal/integrations/aftership/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type createATrackingActionProps struct {
	TrackingNumber string `json:"tracking_number,omitempty"`
	Slug           string `json:"slug"`
}

type CreateATrackingAction struct{}

func (c CreateATrackingAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func (c CreateATrackingAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (c CreateATrackingAction) Name() string {
	return "Create A Tracking"
}

func (c CreateATrackingAction) Description() string {
	return "create a new tracking"
}

func (c CreateATrackingAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createTrackingDocs,
	}
}

func (c CreateATrackingAction) Icon() *string {
	return nil
}

func (c CreateATrackingAction) SampleData() sdkcore.JSON {
	return nil
}

func (c CreateATrackingAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"tracking_number": autoform.NewShortTextField().
			SetDisplayName("Tracking Number").
			SetDescription("tracking number of the shipment").
			SetRequired(true).
			Build(),
		"slug": autoform.NewSelectField().
			SetDisplayName("Slug").
			SetDescription("Unique courier code.").
			SetOptions(shared.CourierCodes).
			SetRequired(true).
			Build(),
	}
}

func (c CreateATrackingAction) Auth() *sdk.Auth {
	return &sdk.Auth{
		Inherit: true,
	}
}

func (c CreateATrackingAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createATrackingActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	afterShipSdk, err := tracking.New(tracking.WithApiKey(ctx.Auth.Extra["api-key"]))
	if err != nil {
		return nil, err
	}
	data := model.CreateTrackingRequest{
		TrackingNumber: input.TrackingNumber,
		Slug:           input.Slug,
	}
	result, err := afterShipSdk.Tracking.CreateTracking().BuildBody(data).Execute()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func NewCreateATrackingAction() sdk.Action {
	return &CreateATrackingAction{}
}
