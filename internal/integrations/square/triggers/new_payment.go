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

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/square/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type NewPaymentTrigger struct{}

func (e *NewPaymentTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "new_payment",
		DisplayName:   "New Payment",
		Description:   "triggers workflow when a new payment is created",
		Type:          core.TriggerTypeScheduled,
		Documentation: newPaymentDocs,
		SampleOutput: map[string]interface{}{
			"kind":     "drive#file",
			"mimeType": "image/jpeg",
			"id":       "1dpv4-sKJfKRwI9qx1vWqQhEGEn3EpbI5",
			"name":     "example.jpg",
		},
	}
}

func (e *NewPaymentTrigger) Auth() *core.AuthMetadata {
	return nil
}

func (e *NewPaymentTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("new-payment", "New Payment")

	schema := form.Build()
	return schema
}

func (e *NewPaymentTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	return nil
}

func (e *NewPaymentTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

func (e *NewPaymentTrigger) Execute(ctx sdkcontext.ExecuteContext) (core.JSON, error) {
	var fromDate string

	// Get the last run time
	lastRun, err := ctx.GetMetadata("lastRun")
	if err != nil {
		return nil, err
	}

	if lastRun != nil {
		lastRunTime, ok := lastRun.(*time.Time)
		if ok && lastRunTime != nil {
			fromDate = lastRunTime.UTC().Format(time.RFC3339)
		} else {
			defaultStartDate := time.Date(2006, 1, 1, 0, 0, 0, 0, time.UTC)
			fromDate = defaultStartDate.UTC().Format(time.RFC3339)
		}
	} else {
		defaultStartDate := time.Date(2006, 1, 1, 0, 0, 0, 0, time.UTC)
		fromDate = defaultStartDate.UTC().Format(time.RFC3339)
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	request := "/v2/payments?begin_time=" + fromDate

	payments, err := shared.GetSquareClient(authCtx.Token.AccessToken, request)
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func (e *NewPaymentTrigger) Criteria(ctx context.Context) core.TriggerCriteria {
	return core.TriggerCriteria{
		Schedule: &core.ScheduleTriggerCriteria{
			CronExpression: "",
			StartTime:      nil,
			EndTime:        nil,
			TimeZone:       "",
			Enabled:        true,
		},
	}
}

func NewNewPaymentTrigger() sdk.Trigger {
	return &NewPaymentTrigger{}
}
