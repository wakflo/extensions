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
	"github.com/aftership/tracking-sdk-go/v5/model"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type getAllTrackingsActionProps struct {
	Keyword string `json:"keyword"`
}

type GetAllTrackingsAction struct{}

func (c *GetAllTrackingsAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func (c GetAllTrackingsAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (c GetAllTrackingsAction) Name() string {
	return "Get All Trackings"
}

func (c GetAllTrackingsAction) Description() string {
	return "get all available trackings"
}

func (c GetAllTrackingsAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getAllTrackingsDocs,
	}
}

func (c GetAllTrackingsAction) Icon() *string {
	return nil
}

func (c GetAllTrackingsAction) SampleData() sdkcore.JSON {
	return nil
}

func (c GetAllTrackingsAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"keyword": autoform.NewShortTextField().
			SetDisplayName("Keywords").
			SetDescription("keywords to search in tracking").
			SetRequired(false).
			Build(),
	}
}

func (c GetAllTrackingsAction) Auth() *sdk.Auth {
	return &sdk.Auth{
		Inherit: true,
	}
}

func (c GetAllTrackingsAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing aftership api key")
	}
	input, err := sdk.InputToTypeSafely[getAllTrackingsActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

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

func NewGetAllTrackingsAction() sdk.Action {
	return &GetAllTrackingsAction{}
}
