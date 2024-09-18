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

type createATrackingOperationProps struct {
	TrackingNumber string `json:"tracking_number,omitempty"`
	Slug           string `json:"slug"`
}

type CreateATrackingOperation struct {
	options *sdk.OperationInfo
}

func NewCreateATrackingOperation() *GetATrackingOperation {
	return &GetATrackingOperation{
		options: &sdk.OperationInfo{
			Name:        "Create a Tracking",
			Description: "create a new tracking",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"tracking_number": autoform.NewShortTextField().
					SetDisplayName("Tracking Number").
					SetDescription("tracking number of the shipment").
					SetRequired(true).
					Build(),
				"slug": autoform.NewShortTextField().
					SetDisplayName("Slug").
					SetDescription("Unique courier code.").
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

func (c *CreateATrackingOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing aftership api key")
	}

	input := sdk.InputToType[createATrackingOperationProps](ctx)

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

func (c *CreateATrackingOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CreateATrackingOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
