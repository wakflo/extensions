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

	"github.com/aftership/tracking-sdk-go/v5"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/integration"
)

type RetrackATrackingAction struct{}

func (c *RetrackATrackingAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func (c RetrackATrackingAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (c RetrackATrackingAction) Name() string {
	return "Retrack a specific tracking"
}

func (c RetrackATrackingAction) Description() string {
	return "retrack an expired tracking. Max 3 times per tracking."
}

func (c RetrackATrackingAction) Documentation() *integration.OperationDocumentation {
	return &integration.OperationDocumentation{
		Documentation: &retrackATracking,
	}
}

func (c RetrackATrackingAction) Icon() *string {
	return nil
}

func (c RetrackATrackingAction) SampleData() sdkcore.JSON {
	return nil
}

func (c RetrackATrackingAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"tracking_id": autoform.NewShortTextField().
			SetDisplayName("Tracking ID").
			SetDescription("tracking ID").
			SetRequired(true).
			Build(),
	}
}

func (c RetrackATrackingAction) Auth() *integration.Auth {
	return &integration.Auth{
		Inherit: true,
	}
}

func (c RetrackATrackingAction) Perform(ctx integration.PerformContext) (sdkcore.JSON, error) {
	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing aftership api key")
	}
	input, err := integration.InputToTypeSafely[getATrackingActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	afterShipSdk, err := tracking.New(tracking.WithApiKey(ctx.Auth.Extra["api-key"]))
	if err != nil {
		return nil, err
	}
	result, err := afterShipSdk.Tracking.
		RetrackTrackingById().
		BuildPath(input.TrackingID).
		Execute()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func NewRetrackATrackingAction() integration.Action {
	return &RetrackATrackingAction{}
}
