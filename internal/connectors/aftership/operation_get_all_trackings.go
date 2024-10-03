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

type getAllTrackingsOperationProps struct {
	Keyword string `json:"keyword"`
}

type GetAllTrackingsOperation struct {
	options *sdk.OperationInfo
}

func NewGetAllTrackingsOperation() *GetAllTrackingsOperation {
	return &GetAllTrackingsOperation{
		options: &sdk.OperationInfo{
			Name:        "Get All Trackings",
			Description: "get all available trackings",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"keyword": autoform.NewShortTextField().
					SetDisplayName("Keywords").
					SetDescription("keywords to search in tracking").
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

func (c *GetAllTrackingsOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing aftership api key")
	}
	input := sdk.InputToType[getAllTrackingsOperationProps](ctx)

	afterShipSdk, err := tracking.New(tracking.WithApiKey(ctx.Auth.Extra["api-key"]))
	if err != nil {
		return nil, err
	}

	result, err := afterShipSdk.Tracking.
		GetTrackings().
		BuildQuery(model.GetTrackingsQuery{Keyword: input.Keyword}).
		Execute()
	if err != nil {
		return nil, err
	}
	return result.Tracking, nil
}

func (c *GetAllTrackingsOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *GetAllTrackingsOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
