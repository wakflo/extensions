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

	"github.com/hiscaler/woocommerce-go"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type TriggerNewProduct struct {
	options *sdk.TriggerInfo
}

func NewTriggerNewProduct() *TriggerNewProduct {
	return &TriggerNewProduct{
		options: &sdk.TriggerInfo{
			Name:        "New Product",
			Description: "Triggers when a new product is created",
			RequireAuth: true,
			Auth:        sharedAuth,
			Type:        sdkcore.TriggerTypeCron,
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

func (t *TriggerNewProduct) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
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
		formattedTime = utcTime.Format("2006-01-02T15:04:05")
	}

	params := woocommerce.ProductsQueryParams{
		After: formattedTime,
	}

	newProduct, total, totalPages, isLastPage, err := wooClient.Services.Product.All(params)
	if err != nil {
		return nil, err
	}
	fmt.Println(totalPages, total, isLastPage)

	return newProduct, nil
}

func (t *TriggerNewProduct) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return t.Run(ctx)
}

func (t *TriggerNewProduct) OnEnabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *TriggerNewProduct) OnDisabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *TriggerNewProduct) GetInfo() *sdk.TriggerInfo {
	return t.options
}
