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

	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)


type ListProductsOperation struct {
	options *sdk.OperationInfo
}

func NewListProductsOperation() *ListProductsOperation {
	return &ListProductsOperation{
		options: &sdk.OperationInfo{
			Name:        "List Products",
			Description: "Fetch total number of products in store",
			RequireAuth: true,
			Auth:        sharedAuth,
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *ListProductsOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {

	if ctx.Auth.Extra["token"] == "" {
		return nil, errors.New("missing shopify auth token")
	}

	domain := ctx.Auth.Extra["domain"]
	shopName := domain + ".myshopify.com"
	client := getShopifyClient(shopName, ctx.Auth.Extra["token"])

	count, err := client.Product.Count(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	
	    return sdk.JSON(map[string]interface{}{
        "total products": count,
    }), nil
}

func (c *ListProductsOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *ListProductsOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
