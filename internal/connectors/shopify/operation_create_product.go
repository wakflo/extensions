// Copyright 2022-present Wakflo
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
	// "fmt"
	// "strings"
	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type createProductOperationProps struct {
	Title       string `json:"title"`
	BodyHTML    string `json:"bodyHTML"`
	Vendor      string `json:"vendor"`
	ProductType string `json:"productType"`
	Status      string `json:"status"`
	Tags        string `json:"tags"`
}
type CreateProductOperation struct {
	options *sdk.OperationInfo
}

func NewCreateProductOperation() *CreateProductOperation {
	return &CreateProductOperation{
		options: &sdk.OperationInfo{
			Name:        "Create Product",
			Description: "Create a new product.",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"title": autoform.NewShortTextField().
					SetDisplayName("Product title").
					SetDescription("The title of the product.").
					SetRequired(true).
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
					SetDisplayName("Product status").
					SetDescription("The status of the product: active or draft").
					SetOptions(statusFormat).
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
func (c *CreateProductOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["token"] == "" {
		return nil, errors.New("missing shopify auth token")
	}
	input := sdk.InputToType[createProductOperationProps](ctx)

	domain := ctx.Auth.Extra["domain"]
	shopName := domain + ".myshopify.com"
	client := getShopifyClient(shopName, ctx.Auth.Extra["token"])

	newProduct := goshopify.Product{
		Title:       input.Title,
		BodyHTML:    input.BodyHTML,
		Vendor:      input.Vendor,
		ProductType: input.ProductType,
		Status:      goshopify.ProductStatus(input.Status),
		Tags:        input.Tags,
	}
	product, err := client.Product.Create(context.Background(), newProduct)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("product not created ")
	}
	return map[string]interface{}{
		"new product": product,
	}, nil
}
func (c *CreateProductOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}
func (c *CreateProductOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
