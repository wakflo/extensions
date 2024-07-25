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

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type ListDraftOrdersOperation struct {
	options *sdk.OperationInfo
}

func NewListDraftOrdersOperation() *ListDraftOrdersOperation {
	return &ListDraftOrdersOperation{
		options: &sdk.OperationInfo{
			Name:        "List Draft Orders",
			Description: "List Draft Orders",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"projectId": autoform.NewShortTextField().
					SetDisplayName("").
					SetDescription("").
					Build(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *ListDraftOrdersOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["token"] == "" {
		return nil, errors.New("missing shopify auth token")
	}

	domain := ctx.Auth.Extra["domain"]
	shopName := domain + ".myshopify.com"
	client := getShopifyClient(shopName, ctx.Auth.Extra["token"])

	darftOrders, err := client.DraftOrder.List(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	if darftOrders == nil {
		return nil, errors.New("no draft orders found")
	}

	return sdk.JSON(map[string]interface{}{
		"Draft orders": darftOrders,
	}), err
}

func (c *ListDraftOrdersOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *ListDraftOrdersOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
