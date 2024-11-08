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
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type addBatchProductOperationProps struct {
	Products []ProductInput `json:"products"`
}

type ProductInput struct {
	Name           string `json:"name"`
	Brand          string `json:"brand"`
	Category       string `json:"category"`
	ProductCode    string `json:"product_code"`
	Barcode        string `json:"barcode"`
	Cost           string `json:"cost"`
	AdditionalCost string `json:"additional_cost"`
}

type AddBatchProductOperation struct {
	options *sdk.OperationInfo
}

func NewAddBatchProductOperation() *AddBatchProductOperation {
	return &AddBatchProductOperation{
		options: &sdk.OperationInfo{
			Name:        "Add Batch Products",
			Description: "Add Batch Products",
			Auth:        sharedAuth,
			Input:       map[string]*sdkcore.AutoFormSchema{
				// "products": autoform.NewArrayField().
				//	SetDisplayName("Products").
				//	SetDescription("The products to add in batch").
				//	SetItems(
				//		autoform.NewObjectField().
				//			SetFields(map[string]*sdkcore.AutoFormSchema{
				//				"name":            autoform.NewShortTextField().SetDisplayName("Product Name").SetRequired(true).Build(),
				//				"brand":           autoform.NewShortTextField().SetDisplayName("Brand").SetRequired(true).Build(),
				//				"category":        autoform.NewShortTextField().SetDisplayName("Category").SetRequired(true).Build(),
				//				"product_code":    autoform.NewShortTextField().SetDisplayName("Product Code").SetRequired(false).Build(),
				//				"barcode":         autoform.NewShortTextField().SetDisplayName("Barcode").SetRequired(false).Build(),
				//				"cost":            autoform.NewShortTextField().SetDisplayName("Cost").SetRequired(false).Build(),
				//				"additional_cost": autoform.NewShortTextField().SetDisplayName("Additional Cost").SetRequired(false).Build(),
				//			}).Build(),
				//	).
				//	SetRequired(true).Build(),
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

func (c *AddBatchProductOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	input := sdk.InputToType[addBatchProductOperationProps](ctx)

	formData := make([]map[string]string, 0, len(input.Products))

	for _, product := range input.Products {
		formData = append(formData, map[string]string{
			"name":            product.Name,
			"brand":           product.Brand,
			"category":        product.Category,
			"product_code":    product.ProductCode,
			"barcode":         product.Barcode,
			"cost":            product.Cost,
			"additional_cost": product.AdditionalCost,
		})
	}

	endpoint := "/api/v2/add/batch/"
	resp, err := prisyncBatchRequest(ctx.Auth.Extra["api-key"], ctx.Auth.Extra["api-token"], endpoint, formData, false)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *AddBatchProductOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *AddBatchProductOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
