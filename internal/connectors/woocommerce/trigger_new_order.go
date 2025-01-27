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

package woocommerce

import (
	"errors"
	"fmt"
	"time"

	"github.com/hiscaler/woocommerce-go"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type TriggerNewOrder struct {
	options *sdk.TriggerInfo
}

func NewTriggerNewOrder() *TriggerNewOrder {
	return &TriggerNewOrder{
		options: &sdk.TriggerInfo{
			Name:        "New Order Update",
			Description: "Triggers when a new item is added or updated in an order.",
			RequireAuth: true,
			Auth:        sharedAuth,
			Type:        sdkcore.TriggerTypePolling,
			Input: map[string]*sdkcore.AutoFormSchema{
				"tag": autoform.NewShortTextField().
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

func (t *TriggerNewOrder) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	baseURL := ctx.Auth.Extra["shop-url"]
	consumerKey := ctx.Auth.Extra["consumer-key"]
	consumerSecret := ctx.Auth.Extra["consumer-secret"]

	if baseURL == "" || consumerKey == "" || consumerSecret == "" {
		return nil, errors.New("missing WooCommerce authentication credentials")
	}

	wooClient := initializeWooCommerceClient(baseURL, consumerKey, consumerSecret)
	lastRunTime := ctx.Metadata.LastRun

	var formattedTime string
	if lastRunTime != nil {
		utcTime := lastRunTime.UTC()
		formattedTime = utcTime.Format(time.RFC3339)
	}

	params := woocommerce.OrdersQueryParams{
		After: formattedTime,
	}

	newOrder, total, totalPages, isLastPage, err := wooClient.Services.Order.All(params)
	if err != nil {
		return nil, err
	}
	fmt.Println(totalPages, total, isLastPage)

	return newOrder, nil
}

func (t *TriggerNewOrder) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return t.Run(ctx)
}

func (t *TriggerNewOrder) OnEnabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *TriggerNewOrder) OnDisabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *TriggerNewOrder) GetInfo() *sdk.TriggerInfo {
	return t.options
}
