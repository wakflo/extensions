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
	"strconv"
	"strings"

	"github.com/hiscaler/woocommerce-go"
	"github.com/hiscaler/woocommerce-go/entity"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type CreateProductOperation struct {
	options *sdk.OperationInfo
}

type createProductOperationProps struct {
	Name             string  `json:"name"`
	Type             string  `json:"type"`
	RegularPrice     float64 `json:"regular_price"`
	Description      string  `json:"description"`
	ShortDescription string  `json:"short_description"`
	Categories       string  `json:"categories"`
}

func NewCreateProductOperation() *CreateProductOperation {
	return &CreateProductOperation{
		options: &sdk.OperationInfo{
			Name:        "Create a product",
			Description: "create a product",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"name": autoform.NewShortTextField().
					SetDisplayName("Product Name").
					SetDescription("Enter product Name").
					SetRequired(true).
					Build(),
				"type": autoform.NewSelectField().
					SetDisplayName("Type").
					SetDescription("Select the type").
					SetOptions(productType).
					SetRequired(true).
					Build(),
				"description": autoform.NewLongTextField().
					SetDisplayName(" Description").
					SetDescription("Enter product description").
					SetRequired(true).
					Build(),
				"short_description": autoform.NewLongTextField().
					SetDisplayName("Short Description").
					SetDescription("Enter the short description").
					SetRequired(true).
					Build(),
				"regular_price": autoform.NewNumberField().
					SetDisplayName("Regular Price").
					SetDescription("Enter Regular Price").
					SetRequired(true).
					Build(),
				"categories": autoform.NewShortTextField().
					SetDisplayName("Category").
					SetDescription("Enter the category IDs (comma separated)").
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
	baseURL := ctx.Auth.Extra["shop-url"]
	consumerKey := ctx.Auth.Extra["consumer-key"]
	consumerSecret := ctx.Auth.Extra["consumer-secret"]

	if baseURL == "" || consumerKey == "" || consumerSecret == "" {
		return nil, errors.New("missing WooCommerce authentication credentials")
	}

	input := sdk.InputToType[createProductOperationProps](ctx)

	wooClient := initializeWooCommerceClient(baseURL, consumerKey, consumerSecret)

	// Parse categories
	var categories []entity.ProductCategory
	if input.Categories != "" {
		categoryIDs := strings.Split(input.Categories, ",")
		for _, idStr := range categoryIDs {
			id, err := strconv.Atoi(strings.TrimSpace(idStr))
			if err != nil {
				return nil, fmt.Errorf("invalid category ID: %s", idStr)
			}
			categories = append(categories, entity.ProductCategory{ID: id})
		}
	}

	// Create a query parameters struct
	params := woocommerce.CreateProductRequest{
		Name:             input.Name,
		Description:      input.Description,
		Type:             input.Type,
		RegularPrice:     input.RegularPrice,
		ShortDescription: input.ShortDescription,
		Categories:       categories,
	}

	product, err := wooClient.Services.Product.Create(params)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (c *CreateProductOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CreateProductOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
