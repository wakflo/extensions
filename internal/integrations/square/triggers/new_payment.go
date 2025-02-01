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

package triggers

import (
	"context"
	"time"

	"github.com/wakflo/extensions/internal/integrations/square/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/integration"
)

type NewPaymentTrigger struct {
	getParentFolders func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error)
}

func (e *NewPaymentTrigger) Name() string {
	return "New Payment"
}

func (e *NewPaymentTrigger) Description() string {
	return "triggers workflow when a new payment is created"
}

func (e *NewPaymentTrigger) Documentation() *integration.OperationDocumentation {
	return &integration.OperationDocumentation{
		Documentation: &newPaymentDocs,
	}
}

func (e *NewPaymentTrigger) Icon() *string {
	return nil
}

func (e *NewPaymentTrigger) SampleData() (sdkcore.JSON, error) {
	return map[string]interface{}{
		"kind":     "drive#file",
		"mimeType": "image/jpeg",
		"id":       "1dpv4-sKJfKRwI9qx1vWqQhEGEn3EpbI5",
		"name":     "example.jpg",
	}, nil
}

func (e *NewPaymentTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{}
}

func (e *NewPaymentTrigger) Auth() *integration.Auth {
	return &integration.Auth{
		Inherit: true,
	}
}

func (e *NewPaymentTrigger) Start(ctx integration.LifecycleContext) error {
	return nil
}

func (e *NewPaymentTrigger) Stop(ctx integration.LifecycleContext) error {
	return nil
}

func (e *NewPaymentTrigger) Execute(ctx integration.ExecuteContext) (sdkcore.JSON, error) {
	var fromDate string
	lastRunTime := ctx.Metadata().LastRun

	if lastRunTime != nil {
		fromDate = lastRunTime.UTC().Format(time.RFC3339)
	} else {
		defaultStartDate := time.Date(2006, 1, 1, 0, 0, 0, 0, time.UTC)
		fromDate = defaultStartDate.UTC().Format(time.RFC3339)
	}

	request := "/v2/payments?begin_time=" + fromDate

	payments, err := shared.GetSquareClient(ctx.Auth.AccessToken, request)
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func (e *NewPaymentTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypeScheduled
}

func (e *NewPaymentTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{
		Schedule: &sdkcore.ScheduleTriggerCriteria{
			CronExpression: "",
			StartTime:      nil,
			EndTime:        nil,
			TimeZone:       "",
			Enabled:        true,
		},
	}
}

func NewNewPaymentTrigger() integration.Trigger {
	return &NewPaymentTrigger{}
}
