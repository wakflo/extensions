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

package trackingmore

import (
	"errors"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type trackAPackageProps struct {
	TrackingNumber string `json:"tracking_number"`
	CourierCode    string `json:"courier_code"`
	Title          string `json:"title"`
	Note           string `json:"note"`
}

type TrackAPackageOperation struct {
	options *sdk.OperationInfo
}

func NewTrackAPackageOperation() *TrackAPackageOperation {
	return &TrackAPackageOperation{
		options: &sdk.OperationInfo{
			Name:        "Create a tracking",
			Description: "create a tracking for a package",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"tracking_number": autoform.NewShortTextField().
					SetDisplayName("Tracking Number").
					SetDescription("Tracking number of a package.").
					SetRequired(true).
					Build(),
				"courier_code": autoform.NewSelectField().
					SetDisplayName("Courier Code").
					SetDescription("Courier code").
					SetOptions(courierCodes).
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
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *TrackAPackageOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	input := sdk.InputToType[trackAPackageProps](ctx)

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

	response, err := createTracking(endpoint, applicationKey, payload)
	if err != nil {
		return nil, err
	}

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid response format: data field is not an array")
	}

	return sdk.JSON(data), nil
}

func (c *TrackAPackageOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *TrackAPackageOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
