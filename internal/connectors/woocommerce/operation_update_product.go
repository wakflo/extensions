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

	"github.com/hiscaler/woocommerce-go"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type UpdateProductOperation struct {
	options *sdk.OperationInfo
}

type updateProductOperationProps struct {
	ProductID        int     `json:"productId"`
	Name             string  `json:"name"`
	RegularPrice     float64 `json:"regular_price"`
	SalePrice        float64 `json:"sale_price"`
	Description      string  `json:"description"`
	ShortDescription string  `json:"short_description"`
	Weight           string  `json:"weight"`
	Length           float64 `json:"length"`
	Width            float64 `json:"width"`
	Height           float64 `json:"height"`
}

func NewUpdateProductOperation() *UpdateProductOperation {
	return &UpdateProductOperation{
		options: &sdk.OperationInfo{
			Name:        "Update a product",
			Description: "update a product",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"productId": autoform.NewNumberField().
					SetDisplayName("Product ID").
					SetDescription("Enter the product ID").
					SetRequired(true).
					Build(),
				"name": autoform.NewShortTextField().
					SetDisplayName("Product Name").
					SetDescription("Enter product Name").
					Build(),
				"description": autoform.NewLongTextField().
					SetDisplayName(" Description").
					SetDescription("Enter product description").
					Build(),
				"short_description": autoform.NewLongTextField().
					SetDisplayName("Short Description").
					SetDescription("Enter the short description").
					Build(),
				"length": autoform.NewNumberField().
					SetDisplayName("Length").
					SetDescription("Enter Product Length").
					Build(),
				"regular_price": autoform.NewNumberField().
					SetDisplayName("Regular Price").
					SetDescription("Enter Regular Price").
					SetRequired(true).
					Build(),
				"sale_price": autoform.NewNumberField().
					SetDisplayName("Discounted Price").
					SetDescription("Enter Discounted Price").
					SetRequired(true).
					Build(),
				"height": autoform.NewNumberField().
					SetDisplayName("Height").
					SetDescription("Enter Product Height").
					Build(),
				"width": autoform.NewNumberField().
					SetDisplayName("Width").
					SetDescription("Enter Product Width").
					Build(),
				"weight": autoform.NewShortTextField().
					SetDisplayName("Weight").
					SetDescription("weight").
					Build(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *UpdateProductOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	baseURL := ctx.Auth.Extra["shop-url"]
	consumerKey := ctx.Auth.Extra["consumer-key"]
	consumerSecret := ctx.Auth.Extra["consumer-secret"]

	if baseURL == "" || consumerKey == "" || consumerSecret == "" {
		return nil, errors.New("missing WooCommerce authentication credentials")
	}

	input := sdk.InputToType[updateProductOperationProps](ctx)

	wooClient := initializeWooCommerceClient(baseURL, consumerKey, consumerSecret)

	// Create a query parameters struct
	params := woocommerce.UpdateProductRequest{}

	if input.Name != "" {
		params.Name = input.Name
	}

	if input.Description != "" {
		params.Description = input.Description
	}

	if input.RegularPrice != 0 {
		params.RegularPrice = input.RegularPrice
	}

	if input.SalePrice != 0 {
		params.SalePrice = input.SalePrice
	}

	if input.ShortDescription != "" {
		params.ShortDescription = input.ShortDescription
	}

	if input.Weight != "" {
		params.Weight = input.Weight
	}

	if input.Length != 0 {
		params.Dimensions.Length = input.Length
	}

	if input.Width != 0 {
		params.Dimensions.Width = input.Width
	}

	if input.Height != 0 {
		params.Dimensions.Length = input.Height
	}

	updatedProduct, err := wooClient.Services.Product.Update(input.ProductID, params)
	if err != nil {
		return nil, err
	}

	return updatedProduct, nil
}

func (c *UpdateProductOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *UpdateProductOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
