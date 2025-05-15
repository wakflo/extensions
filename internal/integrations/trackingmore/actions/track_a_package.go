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

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/trackingmore/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type trackAPackageActionProps struct {
	TrackingNumber string `json:"tracking_number"`
	CourierCode    string `json:"courier_code"`
	Title          string `json:"title"`
	Note           string `json:"note"`
}

type TrackAPackageAction struct{}

func (c *TrackAPackageAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_a_tracking",
		DisplayName:   "Create A tracking",
		Description:   "create a tracking for a package",
		Type:          core.ActionTypeAction,
		Documentation: createTrackingDocs,
		SampleOutput:  nil,
		Settings:      core.ActionSettings{},
	}
}

func (c *TrackAPackageAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_a_tracking", "Create A tracking")

	form.TextField("tracking_number", "Tracking Number").
		Placeholder("Tracking number of a package.").
		Required(true).
		HelpText("Tracking number of a package.")

	form.SelectField("courier_code", "Courier Code").
		Required(true).
		AddOptions(shared.CourierCodes...).
		Placeholder("Courier code").
		HelpText("Courier code")

	form.TextField("title", "Title").
		Placeholder("Title of the package.").
		Required(true).
		HelpText("Title of the package.")

	form.TextareaField("note", "Note").
		Placeholder("Note about the package").
		HelpText("Note about the package")

	schema := form.Build()
	return schema
}

func (c *TrackAPackageAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[trackAPackageActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	endpoint := "/v4/trackings/create"
	applicationKey := authCtx.Extra["key"]

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

	return core.JSON(data), nil
}

func (c *TrackAPackageAction) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Inherit: true,
	}
}

func NewTrackAPackageAction() sdk.Action {
	return &TrackAPackageAction{}
}
