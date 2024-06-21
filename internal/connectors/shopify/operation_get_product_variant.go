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

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type getProductVariantOperationProps struct {
	ProductID uint64 `json:"productId"`
	VariantID uint64 `json:"variantId"`
}

type GetProductVariantOperation struct {
	options *sdk.OperationInfo
}

func NewGetProductVariantOperation() *GetProductVariantOperation {
	return &GetProductVariantOperation{
		options: &sdk.OperationInfo{
			Name:        "Get Product Variant",
			Description: "Get a product variant",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"productId": autoform.NewNumberField().
					SetDisplayName("Product ID").
					SetDescription("ID of product").
					SetRequired(false).
					Build(),
				"variantId": autoform.NewNumberField().
					SetDisplayName("Variant ID").
					SetDescription("product variant ID").
					SetRequired(false).
					Build(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *GetProductVariantOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {

	if  ctx.Auth.Extra["token"] == "" {
		return nil, errors.New("missing shopify auth token")
	}

	input := sdk.InputToType[getProductVariantOperationProps](ctx)

	domain := ctx.Auth.Extra["domain"]
	shopName := domain + ".myshopify.com"

	client := getShopifyClient(shopName, ctx.Auth.Extra["token"])

	product, err := client.Product.Get(context.Background(), input.ProductID, nil)
	if err != nil {
		return nil, err
	}

		if product == nil {
		return nil, fmt.Errorf("no product found with ID '%d'", input.ProductID)
	}

	for _, variant := range product.Variants {
		if variant.Id == input.VariantID {
			return variant, nil
		}
	}
	return nil, nil
}

func (c *GetProductVariantOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *GetProductVariantOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
