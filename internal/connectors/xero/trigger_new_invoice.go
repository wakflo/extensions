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

package xero

import (
	"errors"
	"fmt"
	"strings"

	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type TriggerNewInvoice struct {
	options *sdk.TriggerInfo
}

func NewTriggerNewInvoice() *TriggerNewInvoice {
	return &TriggerNewInvoice{
		options: &sdk.TriggerInfo{
			Name:        "New Invoice",
			Description: "triggers workflow when a new invoice is created",
			RequireAuth: true,
			Auth:        sharedAuth,
			Type:        sdkcore.TriggerTypeCron,
			Input:       map[string]*sdkcore.AutoFormSchema{},
			Settings:    &sdkcore.TriggerSettings{},
			ErrorSettings: &sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (t *TriggerNewInvoice) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing Xero access token")
	}
	var endpoint string

	lastRunTime := ctx.Metadata.LastRun

	if lastRunTime != nil {
		fromDate := lastRunTime.Format("2006-01-02")
		fromDate = strings.ReplaceAll(fromDate, "-", ",")
		endpoint = fmt.Sprintf("/Invoices?where=Date>=DateTime(%s)", fromDate)
	} else {
		endpoint = "/Invoices"
	}

	invoices, err := getXeroNewClient(ctx.Auth.AccessToken, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch invoices: %v", err)
	}

	return invoices, nil
}

func (t *TriggerNewInvoice) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return t.Run(ctx)
}

func (t *TriggerNewInvoice) OnEnabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *TriggerNewInvoice) OnDisabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *TriggerNewInvoice) GetInfo() *sdk.TriggerInfo {
	return t.options
}
