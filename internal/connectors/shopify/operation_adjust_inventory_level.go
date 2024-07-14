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

	shopify "github.com/bold-commerce/go-shopify/v4"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type adjustInventoryLevelOperationProps struct {
	InventoryItemID     uint64 `json:"inventoryItemId"`
	LocationID          uint64 `json:"locationId"`
	AvailableAdjustment int    `json:"available_adjustment"`
}
type AdjustInventoryLevelOperation struct {
	options *sdk.OperationInfo
}

func NewAdjustInventoryLevelOperation() *AdjustInventoryLevelOperation {
	return &AdjustInventoryLevelOperation{
		options: &sdk.OperationInfo{
			Name:        "Adjust Inventory Level",
			Description: "Adjust inventory level of an item at a location.",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"inventoryItemId": autoform.NewNumberField().
					SetDisplayName("Inventory Item").
					SetDescription("The ID of the inventory item.").
					SetRequired(true).
					Build(),
				"locationId": autoform.NewNumberField().
					SetDisplayName("Location").
					SetDescription("The ID of the location.").
					SetRequired(true).
					Build(),
				"available_adjustment": autoform.NewNumberField().
					SetDisplayName("Adjustment Quantity").
					SetDescription("Positive values increase inventory, negative values decrease it.").
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

func (c *AdjustInventoryLevelOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["token"] == "" {
		return nil, errors.New("missing shopify auth token")
	}
	input := sdk.InputToType[adjustInventoryLevelOperationProps](ctx)
	domain := ctx.Auth.Extra["domain"]
	shopName := domain + ".myshopify.com"
	client := getShopifyClient(shopName, ctx.Auth.Extra["token"])

	options := shopify.InventoryLevelAdjustOptions{
		InventoryItemId: input.InventoryItemID,
		LocationId:      input.LocationID,
		Adjust:          input.AvailableAdjustment,
	}

	inventoryLevel, err := client.InventoryLevel.Adjust(context.Background(), options)
	if err != nil {
		return nil, fmt.Errorf("failed to adjust inventory: %v", err)
	}

	return sdk.JSON(map[string]interface{}{
		"adjustedInventoryLevel": inventoryLevel,
	}), nil
}

func (c *AdjustInventoryLevelOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *AdjustInventoryLevelOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
