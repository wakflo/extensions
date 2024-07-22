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

	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type TriggerNewOrder struct {
	options *sdk.TriggerInfo
}

func NewTriggerNewOrder() *TriggerNewOrder {
	return &TriggerNewOrder{
		options: &sdk.TriggerInfo{
			Name:        "New Order",
			Description: "Triggers when a new order is created",
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

func (t *TriggerNewOrder) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
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
		CreatedAtMin: lastRunTime,
	}

	orders, err := client.Order.List(context.Background(), options)
	if err != nil {
		return nil, err
	}

	return sdk.JSON(map[string]interface{}{
		"order details": orders,
	}), err
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
