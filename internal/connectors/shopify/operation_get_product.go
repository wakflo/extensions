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

	// "strings"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type getProductOperationProps struct {
	ProductID uint64  `json:"productId"`
}

type GetProductOperation struct {
	options *sdk.OperationInfo
}

func NewGetProductOperation() *GetProductOperation {
	return &GetProductOperation{
		options: &sdk.OperationInfo{
			Name:        "Get Product",
			Description: "Get an existing product by id.",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"productId": autoform.NewNumberField().
					SetDisplayName("Product").
					SetDescription("The ID of the product.").
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

func (c *GetProductOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {

	if ctx.Auth.Extra["token"] == "" {
		return nil, errors.New("missing shopify auth token")
	}

	input := sdk.InputToType[getProductOperationProps](ctx)

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

	productMap := map[string]interface{}{
		"ID":          product.Id,
		"Title":       product.Title,
		"Description": product.BodyHTML,
		"Price":       product.Variants[0].Price,
		"Variants":    product.Variants,
	}

		    return sdk.JSON(map[string]interface{}{
        "product details": productMap,
    }), nil
}

func (c *GetProductOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *GetProductOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
