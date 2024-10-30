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

package prisync

import (
	"net/http"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type editProductOperationProps struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Brand       string `json:"brand"`
	Category    string `json:"category"`
	BarCode     string `json:"barcode"`
	Cost        string `json:"cost"`
	ProductCode string `json:"product_code"`
	Tags        string `json:"tags"`
}

type EditProductOperation struct {
	options *sdk.OperationInfo
}

func NewEditProductOperation() *EditProductOperation {
	return &EditProductOperation{
		options: &sdk.OperationInfo{
			Name:        "Edit Product",
			Description: "Edit a Product",
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"id": autoform.NewShortTextField().
					SetDisplayName("Product ID").
					SetDescription("ID of product").
					SetRequired(true).Build(),
				"name": autoform.NewShortTextField().
					SetDisplayName("Product Name").
					SetDescription("name of product").
					SetRequired(false).Build(),
				"brand": autoform.NewShortTextField().
					SetDisplayName("Brand").
					SetDescription("Brand name").
					SetRequired(false).Build(),
				"category": autoform.NewShortTextField().
					SetDisplayName("Category").
					SetDescription("Category name").
					SetRequired(false).Build(),
				"product_code": autoform.NewShortTextField().
					SetDisplayName("Product Code").
					SetDescription("Product code").
					SetRequired(false).Build(),
				"barcode": autoform.NewShortTextField().
					SetDisplayName("Bar Code").
					SetDescription("Bar code").
					SetRequired(false).Build(),
				"cost": autoform.NewShortTextField().
					SetDisplayName("Product Cost").
					SetDescription("Product cost").
					SetRequired(false).Build(),
				"tags": autoform.NewLongTextField().
					SetDisplayName("Product Tags").
					SetDescription("Product tags").
					SetRequired(false).Build(),
			},
			SampleOutput: map[string]interface{}{},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
			RequireAuth: true,
		},
	}
}

func (c *EditProductOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	input := sdk.InputToType[editProductOperationProps](ctx)
	formData := map[string]string{}

	if input.Name != "" {
		formData["name"] = input.Name
	}

	if input.Brand != "" {
		formData["brand"] = input.Brand
	}

	if input.Category != "" {
		formData["category"] = input.Category
	}

	if input.ProductCode != "" {
		formData["product_code"] = input.ProductCode
	}

	if input.BarCode != "" {
		formData["barcode"] = input.BarCode
	}

	if input.Cost != "" {
		formData["cost"] = input.Cost
	}

	if input.Tags != "" {
		formData["tags"] = input.Tags
	}

	reqURL := "https://prisync.com/api/v2/edit/product/id/" + input.ID
	resp, err := prisyncRequest(ctx.Auth.Extra["api-key"], ctx.Auth.Extra["api-token"], reqURL, http.MethodPost, formData)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *EditProductOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *EditProductOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
