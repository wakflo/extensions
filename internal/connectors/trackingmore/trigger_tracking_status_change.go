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
	"log"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type TrackingStatusChange struct {
	options *sdk.TriggerInfo
}

func NewTrackingStatusChange() *TrackingStatusChange {
	return &TrackingStatusChange{
		options: &sdk.TriggerInfo{
			Name:        "Tracking Status Change",
			Description: "triggers workflow when a new sales is initiated",
			RequireAuth: true,
			Auth:        sharedAuth,
			Type:        sdkcore.TriggerTypeCron,
			Input: map[string]*sdkcore.AutoFormSchema{
				"id": autoform.NewShortTextField().
					SetDisplayName("").
					SetDescription("").
					Build(),
			},
			Settings: &sdkcore.TriggerSettings{},
			ErrorSettings: &sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (t *TrackingStatusChange) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	endpoint := "/v4/trackings/get"
	applicationKey := ctx.Auth.Extra["key"]

	lastRunTime := ctx.Metadata.LastRun

	var fromDate string
	if lastRunTime != nil {
		fromDate = lastRunTime.UTC().Format("2006-01-02T15:04:05-07:00")
	}

	response, err := listTracking(endpoint, applicationKey, fromDate)
	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}

	data, ok := response["data"].([]interface{})
	if !ok {
		return nil, errors.New("invalid response format: data field is not an array")
	}

	return data, nil
}

func (t *TrackingStatusChange) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return t.Run(ctx)
}

func (t *TrackingStatusChange) OnEnabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *TrackingStatusChange) OnDisabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *TrackingStatusChange) GetInfo() *sdk.TriggerInfo {
	return t.options
}
