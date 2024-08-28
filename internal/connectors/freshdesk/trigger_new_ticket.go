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

package freshdesk

import (
	"errors"
	"fmt"
	"time"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type TriggerNewTicket struct {
	options *sdk.TriggerInfo
}

func NewTriggerNewTicket() *TriggerNewTicket {
	return &TriggerNewTicket{
		options: &sdk.TriggerInfo{
			Name:        "New Ticket",
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

func (t *TriggerNewTicket) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["api-key"] == "" || ctx.Auth.Extra["domain"] == "" {
		return nil, errors.New("missing freshdesk auth values")
	}

	domain := ctx.Auth.Extra["domain"]
	freshdeskDomain := "https://" + domain + ".freshdesk.com"
	fmt.Println(freshdeskDomain)

	lastRunTime := ctx.Metadata.LastRun

	var fromDate string
	if lastRunTime != nil {
		fromDate = lastRunTime.UTC().Format(time.RFC3339)
	} else {
		fromDate = ""
	}

	tickets, err := GetTicketQuery(freshdeskDomain, ctx.Auth.Extra["api-key"], fromDate)
	if err != nil {
		return nil, err
	}

	return tickets, nil
}

func (t *TriggerNewTicket) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return t.Run(ctx)
}

func (t *TriggerNewTicket) OnEnabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *TriggerNewTicket) OnDisabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *TriggerNewTicket) GetInfo() *sdk.TriggerInfo {
	return t.options
}
