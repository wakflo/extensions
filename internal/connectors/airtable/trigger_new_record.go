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

package airtable

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type listNewRecordProps struct {
	Bases string `json:"bases,omitempty"`
}

type RecordTrigger struct {
	options *sdk.TriggerInfo
}

func NewRecordTrigger() *RecordTrigger {
	return &RecordTrigger{
		options: &sdk.TriggerInfo{
			Name:        "New Record Created",
			Description: "triggers workflow when a new record is created",
			RequireAuth: true,
			Auth:        sharedAuth,
			Type:        sdkcore.TriggerTypeCron,
			Input: map[string]*sdkcore.AutoFormSchema{
				"bases": getBasesInput(),
			},
			Settings: &sdkcore.TriggerSettings{},
			ErrorSettings: &sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (t *RecordTrigger) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing airtable api key")
	}
	apiKey := ctx.Auth.Extra["api-key"]
	input := sdk.InputToType[listNewRecordProps](ctx)
	lastRunTime := ctx.Metadata.LastRun

	var createdTime string
	if lastRunTime != nil {
		createdTime = lastRunTime.UTC().Format(time.RFC3339)
	}
	reqURL := fmt.Sprintf("%s/v0/meta/bases/%s/tables?updated_since=%s", baseAPI, input.Bases, createdTime)

	response, err := airtableRequest(apiKey, reqURL, http.MethodGet)
	if err != nil {
		return nil, errors.New("error fetching data")
	}

	return response, nil
}

func (t *RecordTrigger) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return t.Run(ctx)
}

func (t *RecordTrigger) OnEnabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *RecordTrigger) OnDisabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *RecordTrigger) GetInfo() *sdk.TriggerInfo {
	return t.options
}
