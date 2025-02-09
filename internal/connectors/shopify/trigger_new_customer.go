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

package shopify

import (
	"context"
	"errors"
	"time"

	goshopify "github.com/bold-commerce/go-shopify/v4"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type TriggerNewCustomer struct {
	options *sdk.TriggerInfo
}

func NewTriggerNewCustomer() *TriggerNewCustomer {
	return &TriggerNewCustomer{
		options: &sdk.TriggerInfo{
			Name:        "New Customer",
			Description: "Triggers when a new customer is created",
			RequireAuth: true,
			Auth:        sharedAuth,
			Type:        sdkcore.TriggerTypePolling,
			Input: map[string]*sdkcore.AutoFormSchema{
				"tag": autoform.NewShortTextField().
					SetDisplayName("").
					SetDescription("").
					SetRequired(false).
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

func (t *TriggerNewCustomer) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["token"] == "" {
		return nil, errors.New("missing shopify auth token")
	}

	domain := ctx.Auth.Extra["domain"]
	shopName := domain + ".myshopify.com"
	client := getShopifyClient(shopName, ctx.Auth.Extra["token"])

	var lastRunTime time.Time
	if ctx.Metadata.LastRun != nil {
		lastRunTime = *ctx.Metadata.LastRun
	}

	options := &goshopify.ListOptions{
		CreatedAtMin: lastRunTime.UTC(),
	}

	customers, err := client.Customer.List(context.Background(), options)
	if err != nil {
		return nil, err
	}

	return customers, nil
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
