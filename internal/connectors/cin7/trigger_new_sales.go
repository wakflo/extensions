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

package cin7

import (
	"log"
	"time"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type TriggerNewSales struct {
	options *sdk.TriggerInfo
}

func NewTriggerNewSales() *TriggerNewSales {
	return &TriggerNewSales{
		options: &sdk.TriggerInfo{
			Name:        "New Sales",
			Description: "triggers workflow when a new sales is initiated",
			RequireAuth: true,
			Auth:        sharedAuth,
			Type:        sdkcore.TriggerTypePolling,
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

func (t *TriggerNewSales) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	endpoint := "/ExternalApi/SaleList"
	accountID := ctx.Auth.Extra["account_id"]
	applicationKey := ctx.Auth.Extra["key"]

	lastRunTime := ctx.Metadata.LastRun

	var fromDate string
	if lastRunTime != nil {
		fromDate = lastRunTime.UTC().Format(time.RFC3339)
	}

	queryParams := map[string]interface{}{
		"Page":         1,
		"CreatedSince": fromDate,
	}

	response, err := fetchData(endpoint, accountID, applicationKey, queryParams)
	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}

	return response, nil
}

func (t *TriggerNewSales) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return t.Run(ctx)
}

func (t *TriggerNewSales) OnEnabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *TriggerNewSales) OnDisabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *TriggerNewSales) GetInfo() *sdk.TriggerInfo {
	return t.options
}
