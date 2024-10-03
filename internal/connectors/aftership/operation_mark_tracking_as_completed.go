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
	"fmt"

	"github.com/aftership/tracking-sdk-go/v5"
	"github.com/aftership/tracking-sdk-go/v5/model"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type MarkTrackingAsCompletedOperation struct {
	options *sdk.OperationInfo
}

func NewMarkTrackingAsCompletedOperation() *MarkTrackingAsCompletedOperation {
	return &MarkTrackingAsCompletedOperation{
		options: &sdk.OperationInfo{
			Name:        "Mark a tracking as completed",
			Description: "Mark a tracking as completed. The tracking won't auto update until retrack it.",
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

func (c *MarkTrackingAsCompletedOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing aftership api key")
	}

	input := sdk.InputToType[getATrackingOperationProps](ctx)

	afterShipSdk, err := tracking.New(tracking.WithApiKey(ctx.Auth.Extra["api-key"]))
	if err != nil {
		fmt.Println(err)
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

func (c *MarkTrackingAsCompletedOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *MarkTrackingAsCompletedOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}