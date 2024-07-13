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

package zohoinventory

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/wakflo/go-sdk/autoform"
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

func (t *TriggerNewPayment) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	input := sdk.InputToType[getPaymentListOperationProps](ctx)

	// Get the last run time or use a default
	lastRunTime := ctx.Metadata.LastRun
	if lastRunTime == nil {
		defaultTime := time.Now().Add(-24 * time.Hour)
		lastRunTime = &defaultTime
	}

	fromDate := lastRunTime.Format("2006-01-02")

	baseURL := "https://www.zohoapis.com/inventory/v1/customerpayments"

	// Create URL with query parameters
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing URL: %v", err)
	}

	q := u.Query()
	q.Set("organization_id", input.OrganizationID)
	q.Set("date", fromDate)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Zoho-oauthtoken "+ctx.Auth.Token.AccessToken)
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return result, nil
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
