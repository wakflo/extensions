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

package square

import (
	"time"

	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type TriggerNewPayment struct {
	options *sdk.TriggerInfo
}

func NewTriggerNewPayment() *TriggerNewPayment {
	return &TriggerNewPayment{
		options: &sdk.TriggerInfo{
			Name:        "New Payment",
			Description: "triggers workflow when a new payment is created",
			RequireAuth: true,
			Auth:        sharedAuth,
			Strategy:    sdkcore.TriggerStrategyPolling,
			Input:       map[string]*sdkcore.AutoFormSchema{},
			Settings:    &sdkcore.TriggerSettings{},
			ErrorSettings: &sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (t *TriggerNewPayment) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	_ = sdk.InputToType[getPaymentsOperationProps](ctx)
	var fromDate string
	lastRunTime := ctx.Metadata.LastRun

	if lastRunTime != nil {
		fromDate = lastRunTime.UTC().Format(time.RFC3339)
	} else {
		defaultStartDate := time.Date(2006, 1, 1, 0, 0, 0, 0, time.UTC)
		fromDate = defaultStartDate.UTC().Format(time.RFC3339)
	}

	request := "/v2/payments?begin_time=" + fromDate

	payments, err := getSquareClient(ctx.Auth.AccessToken, request)
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func (t *TriggerNewPayment) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return t.Run(ctx)
}

func (t *TriggerNewPayment) OnEnabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *TriggerNewPayment) OnDisabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *TriggerNewPayment) GetInfo() *sdk.TriggerInfo {
	return t.options
}
