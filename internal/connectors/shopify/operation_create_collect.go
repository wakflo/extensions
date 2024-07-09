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
	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type createCollectOperationProps struct {
	ProductID    uint64 `json:"productId"`
	CollectionID uint64 `json:"collectionId"`
}
type CreateCollectOperation struct {
	options *sdk.OperationInfo
}

func NewCreateCollectOperation() *CreateCollectOperation {
	return &CreateCollectOperation{
		options: &sdk.OperationInfo{
			Name:        "Create Collect",
			Description: "Add a product to a collection.",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"productId": autoform.NewNumberField().
					SetDisplayName("Product ID").
					SetDescription("The ID of the product.").
					SetRequired(true).
					Build(),
				"collectionID": autoform.NewNumberField().
					SetDisplayName("Collection ID").
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

func (c *CreateCollectOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["token"] == "" {
		return nil, errors.New("missing shopify auth token")
	}
	input := sdk.InputToType[createCollectOperationProps](ctx)
	domain := ctx.Auth.Extra["domain"]
	shopName := domain + ".myshopify.com"

	client := getShopifyClient(shopName, ctx.Auth.Extra["token"])

	collect := goshopify.Collect{
		CollectionId: input.CollectionID,
		ProductId:    input.ProductID,
	}

	newCollect, err := client.Collect.Create(context.Background(), collect)
	if err != nil {
		return nil, err
	}
	if newCollect == nil {
		return nil, fmt.Errorf("no collection found with ID '%d'", input.CollectionID)
	}
	return sdk.JSON(map[string]interface{}{
		"collection details": newCollect,
	}), nil
}

func (c *CreateCollectOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CreateCollectOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
