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

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type FindCouponOperation struct {
	options *sdk.OperationInfo
}

type findCouponOperationProps struct {
	CouponID int `json:"couponId"`
}

func NewFindCouponOperation() *FindCouponOperation {
	return &FindCouponOperation{
		options: &sdk.OperationInfo{
			Name:        "Find Coupon",
			Description: "Find a coupon",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"couponId": autoform.NewNumberField().
					SetDisplayName("Product ID").
					SetDescription("Enter the product ID").
					SetRequired(true).
					Build(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *FindCouponOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	baseURL := ctx.Auth.Extra["shop-url"]
	consumerKey := ctx.Auth.Extra["consumer-key"]
	consumerSecret := ctx.Auth.Extra["consumer-secret"]

	if baseURL == "" || consumerKey == "" || consumerSecret == "" {
		return nil, errors.New("missing WooCommerce authentication credentials")
	}

	input := sdk.InputToType[findCouponOperationProps](ctx)

	wooClient := initializeWooCommerceClient(baseURL, consumerKey, consumerSecret)

	product, err := wooClient.Services.Coupon.One(input.CouponID)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (c *FindCouponOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *FindCouponOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
