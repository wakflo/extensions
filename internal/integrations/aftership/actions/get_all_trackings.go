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
	"github.com/juicycleff/smartform/v1"

	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type getAllTrackingsActionProps struct {
	Keyword string `json:"keyword"`
}

type GetAllTrackingsAction struct{}

// Metadata returns metadata about the action
func (c GetAllTrackingsAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_all_trackings",
		DisplayName:   "Get All Trackings",
		Description:   "get all available trackings",
		Type:          core.ActionTypeAction,
		Documentation: getAllTrackingsDocs,
		SampleOutput:  nil,
		Settings:      core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (c GetAllTrackingsAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_all_trackings", "Get All Trackings")

	form.TextField("keyword", "Keywords").
		Placeholder("Enter keywords").
		Required(false).
		HelpText("keywords to search in tracking")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (c GetAllTrackingsAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (c GetAllTrackingsAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	if authCtx.Extra["api-key"] == "" {
		return nil, errors.New("missing aftership api key")
	}

	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[getAllTrackingsActionProps](ctx)
	if err != nil {
		return nil, err
	}

	afterShipSdk, err := tracking.New(tracking.WithApiKey(authCtx.Extra["api-key"]))
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
