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

package aftership

import (
	"errors"

	"github.com/aftership/tracking-sdk-go/v5"
	"github.com/aftership/tracking-sdk-go/v5/model"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type detectCourierOperationProps struct {
	TrackingNumber string `json:"tracking_number,omitempty"`
}

type DetectCourierOperation struct {
	options *sdk.OperationInfo
}

func NewDetectCourierOperation() *DetectCourierOperation {
	return &DetectCourierOperation{
		options: &sdk.OperationInfo{
			Name:        "Detect Courier",
			Description: "Returns a list of matched couriers based on tracking number format",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"tracking_number": autoform.NewShortTextField().
					SetDisplayName("Tracking Number").
					SetDescription("tracking number of the shipment").
					SetRequired(true).
					Build(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *DetectCourierOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing aftership api key")
	}

	input := sdk.InputToType[detectCourierOperationProps](ctx)

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

func (c *DetectCourierOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *DetectCourierOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
