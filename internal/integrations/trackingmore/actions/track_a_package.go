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

	"github.com/wakflo/extensions/internal/integrations/trackingmore/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type trackAPackageActionProps struct {
	TrackingNumber string `json:"tracking_number"`
	CourierCode    string `json:"courier_code"`
	Title          string `json:"title"`
	Note           string `json:"note"`
}

type TrackAPackageAction struct{}

func (c *TrackAPackageAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func (c TrackAPackageAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (c TrackAPackageAction) Name() string {
	return "Create A tracking"
}

func (c TrackAPackageAction) Description() string {
	return "create a tracking for a package"
}

func (c TrackAPackageAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createTrackingDocs,
	}
}

func (c TrackAPackageAction) Icon() *string {
	return nil
}

func (c TrackAPackageAction) SampleData() sdkcore.JSON {
	return nil
}

func (c TrackAPackageAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"tracking_number": autoform.NewShortTextField().
			SetDisplayName("Tracking Number").
			SetDescription("Tracking number of a package.").
			SetRequired(true).
			Build(),
		"courier_code": autoform.NewSelectField().
			SetDisplayName("Courier Code").
			SetDescription("Courier code").
			SetOptions(shared.CourierCodes).
			SetRequired(true).
			Build(),
		"title": autoform.NewShortTextField().
			SetDisplayName("Title").
			SetDescription("Title of the package.").
			SetRequired(true).
			Build(),
		"note": autoform.NewLongTextField().
			SetDisplayName("Note").
			SetDescription("Note about the package").
			SetRequired(false).
			Build(),
	}
}

func (c TrackAPackageAction) Auth() *sdk.Auth {
	return &sdk.Auth{
		Inherit: true,
	}
}

func (c TrackAPackageAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[trackAPackageActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	endpoint := "/v4/trackings/create"
	applicationKey := ctx.Auth.Extra["key"]

	payload := map[string]interface{}{
		"courier_code":    input.CourierCode,
		"tracking_number": input.TrackingNumber,
	}

	if input.Title != "" {
		payload["title"] = input.Title
	}
	if input.Note != "" {
		payload["note"] = input.Note
	}

	response, err := shared.CreateATracking(endpoint, applicationKey, payload)
	if err != nil {
		return nil, err
	}

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid response format: data field is not an array")
	}

	return sdk.JSON(data), nil
}

func NewTrackAPackageAction() sdk.Action {
	return &TrackAPackageAction{}
}
