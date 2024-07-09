// Copyright 2022-present Wakflo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
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
	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type updateProductOperationProps struct {
	ProductID   uint64 `json:"productID"`
	Title       string `json:"title"`
	BodyHTML    string `json:"bodyHTML"`
	Vendor      string `json:"vendor"`
	ProductType string `json:"productType"`
	Status      string `json:"status"`
	Tags        string `json:"tags"`
}
type UpdateProductOperation struct {
	options *sdk.OperationInfo
}

func NewUpdateProductOperation() *UpdateProductOperation {
	return &UpdateProductOperation{
		options: &sdk.OperationInfo{
			Name:        "Update Product",
			Description: "update an existing product.",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"productID": autoform.NewNumberField().
					SetDisplayName("Product ID").
					SetDescription("The id of the product.").
					SetRequired(true).
					Build(),
				"title": autoform.NewShortTextField().
					SetDisplayName("Product title").
					SetDescription("The title of the product.").
					SetRequired(false).
					Build(),
				"bodyHTML": autoform.NewLongTextField().
					SetDisplayName("Product description").
					SetDescription("The description of the product.").
					SetRequired(false).
					Build(),
				"vendor": autoform.NewShortTextField().
					SetDisplayName("Vendor").
					SetDescription("Vendor.").
					SetRequired(false).
					Build(),
				"productType": autoform.NewShortTextField().
					SetDisplayName("Product type").
					SetDescription("A categorization for the product used for filtering and searching products.").
					SetRequired(false).
					Build(),
				"tags": autoform.NewLongTextField().
					SetDisplayName("Tags").
					SetDescription("A string of comma-separated tags for filtering and search.").
					SetRequired(false).
					Build(),
				"status": autoform.NewSelectField().
					SetDisplayName("Status").
					SetDescription("The status of the product: active or draft").
					SetOptions(statusFormat).
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
func (c *UpdateProductOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["token"] == "" {
		return nil, errors.New("missing shopify auth token")
	}
	input := sdk.InputToType[updateProductOperationProps](ctx)
	domain := ctx.Auth.Extra["domain"]
	shopName := domain + ".myshopify.com"
	client := getShopifyClient(shopName, ctx.Auth.Extra["token"])

	existingProduct, err := client.Product.Get(context.Background(), input.ProductID, nil)

	if err != nil {
		return nil, err
	}
	if input.Title != "" {
		existingProduct.Title = input.Title
	}
	if input.BodyHTML != "" {
		existingProduct.BodyHTML = input.BodyHTML
	}
	if input.Vendor != "" {
		existingProduct.Vendor = input.Vendor
	}
	if input.ProductType != "" {
		existingProduct.ProductType = input.ProductType
	}
	if input.Tags != "" {
		existingProduct.Tags = input.Tags
	}
	if input.Status != "" {
		existingProduct.Status = goshopify.ProductStatus(input.Status)
	}

	updatedProduct, err := client.Product.Update(context.Background(), *existingProduct)

	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}
	return map[string]interface{}{
		"updated_product": updatedProduct,
	}, nil
}
func (c *UpdateProductOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}
func (c *UpdateProductOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
