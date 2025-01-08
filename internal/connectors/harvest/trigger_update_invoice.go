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

package harvest

import (
	"errors"
	"log"
	"time"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type InvoiceUpdate struct {
	options *sdk.TriggerInfo
}

func NewInvoiceUpdate() *InvoiceUpdate {
	return &InvoiceUpdate{
		options: &sdk.TriggerInfo{
			Name:        "Invoice Update",
			Description: "Triggers workflow when an invoice has been updated",
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

func (t *InvoiceUpdate) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing Harvest auth token")
	}

	lastRunTime := ctx.Metadata.LastRun

	var updatedTime string
	if lastRunTime != nil {
		updatedTime = lastRunTime.UTC().Format(time.RFC3339)
	}

	url := "/v2/invoices?updated_since=" + updatedTime

	response, err := getHarvestClient(ctx.Auth.AccessToken, url)
	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}

	data, ok := response["invoices"].([]interface{})
	if !ok {
		return nil, errors.New("invalid response format: data field is not an array")
	}

	return data, nil
}

func (t *InvoiceUpdate) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return t.Run(ctx)
}

func (t *InvoiceUpdate) OnEnabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *InvoiceUpdate) OnDisabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *InvoiceUpdate) GetInfo() *sdk.TriggerInfo {
	return t.options
}
