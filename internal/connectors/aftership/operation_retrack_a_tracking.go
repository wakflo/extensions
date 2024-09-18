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

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type getRetrackTrackingOperationProps struct {
	TrackingID string `json:"tracking_id"`
}

type RetrackATrackingOperation struct {
	options *sdk.OperationInfo
}

func NewRetrackATrackingOperation() *GetATrackingOperation {
	return &GetATrackingOperation{
		options: &sdk.OperationInfo{
			Name:        "Retrack a specific tracking",
			Description: "Retrack an expired tracking. Max 3 times per tracking.",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"tracking_id": autoform.NewShortTextField().
					SetDisplayName("Tracking ID").
					SetDescription("tracking ID").
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

func (c *RetrackATrackingOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing aftership api key")
	}

	input := sdk.InputToType[getRetrackTrackingOperationProps](ctx)

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

func (c *RetrackATrackingOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *RetrackATrackingOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
