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

package stripe

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"

	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type TriggerNewCustomer struct {
	options *sdk.TriggerInfo
}

func NewTriggerCustomer() *TriggerNewCustomer {
	return &TriggerNewCustomer{
		options: &sdk.TriggerInfo{
			Name:        "New Customer",
			Description: "triggers workflow when a new customer is created",
			RequireAuth: true,
			Auth:        sharedAuth,
			Type:        sdkcore.TriggerTypePolling,
			Input:       map[string]*sdkcore.AutoFormSchema{},
			Settings:    &sdkcore.TriggerSettings{},
			ErrorSettings: &sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (t *TriggerNewCustomer) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing stripe secret api-key")
	}
	apiKey := ctx.Auth.Extra["api-key"]

	_ = sdk.InputToType[createCustomerOperationProps](ctx)
	var fromDate int64
	lastRunTime := ctx.Metadata.LastRun

	if lastRunTime != nil {
		fromDate = lastRunTime.UTC().Unix()
	}

	stringValue := strconv.FormatInt(fromDate, 10)

	params := url.Values{}
	params.Add("query", "created:'"+stringValue+"'")

	reqURL := "/v1/customers"

	resp, err := stripClient(apiKey, reqURL, http.MethodGet, nil, params)
	if err != nil {
		return nil, err
	}

	nodes, ok := resp["data"].([]interface{})
	if !ok {
		return nil, errors.New("failed to extract data from response")
	}

	return nodes, nil
}

func (t *TriggerNewCustomer) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return t.Run(ctx)
}

func (t *TriggerNewCustomer) OnEnabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *TriggerNewCustomer) OnDisabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *TriggerNewCustomer) GetInfo() *sdk.TriggerInfo {
	return t.options
}
