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
	"github.com/wakflo/go-sdk/sdk"
)

type MarkTrackingAsCompletedAction struct{}

func (c *MarkTrackingAsCompletedAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func (c MarkTrackingAsCompletedAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (c MarkTrackingAsCompletedAction) Name() string {
	return "Mark a tracking as completed"
}

func (c MarkTrackingAsCompletedAction) Description() string {
	return "mark a tracking as completed. The tracking won't auto update until retrack it."
}

func (c MarkTrackingAsCompletedAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &markTrackingAsCompletedDocs,
	}
}

func (c MarkTrackingAsCompletedAction) Icon() *string {
	return nil
}

func (c MarkTrackingAsCompletedAction) SampleData() sdkcore.JSON {
	return nil
}

func (c MarkTrackingAsCompletedAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"tracking_id": autoform.NewShortTextField().
			SetDisplayName("Tracking ID").
			SetDescription("tracking ID").
			SetRequired(true).
			Build(),
	}
}

func (c MarkTrackingAsCompletedAction) Auth() *sdk.Auth {
	return &sdk.Auth{
		Inherit: true,
	}
}

func (c MarkTrackingAsCompletedAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getATrackingActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	afterShipSdk, err := tracking.New(tracking.WithApiKey(ctx.Auth.Extra["api-key"]))
	if err != nil {
		return nil, err
	}
	result, err := afterShipSdk.Tracking.
		MarkTrackingCompletedById().
		BuildPath(input.TrackingID).
		BuildBody(model.MarkTrackingCompletedByIdRequest{Reason: "DELIVERED"}).
		Execute()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func NewMarkTrackingAsCompletedAction() sdk.Action {
	return &MarkTrackingAsCompletedAction{}
}
