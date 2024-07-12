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
	"fmt"
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
			Description: "Triggers workflow when a new customer is created",
			RequireAuth: true,
			Auth:        sharedAuth,
			Type:        sdkcore.TriggerTypeCron,
			Input: map[string]*sdkcore.AutoFormSchema{
				"tag": autoform.NewShortTextField().
					SetDisplayName("Tag").
					SetDescription("Only trigger for customers with this tag").
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

	// Get the last run time from metadata, or use a default if it's the first run
	lastRunTime := ctx.Metadata.LastRun
	if lastRunTime == nil {
		defaultTime := time.Now().Add(-24 * time.Hour)
		lastRunTime = &defaultTime
	}

	query := fmt.Sprintf("created_at:>='%s'", lastRunTime.Format("2006-01-02T15:04:05-07:00"))

	// Set up the search options
	options := &goshopify.CustomerSearchOptions{
		Query: query,
	}

	customers, err := client.Customer.Search(context.Background(), options)
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
