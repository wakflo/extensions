package actions

import (
	"context"
	"fmt"

	shopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/shopify/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type adjustInventoryLevelActionProps struct {
	InventoryItemID     uint64 `json:"inventoryItemId"`
	LocationID          uint64 `json:"locationId"`
	AvailableAdjustment int    `json:"available_adjustment"`
}

type AdjustInventoryLevelAction struct{}

func (a *AdjustInventoryLevelAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "adjust_inventory_level",
		DisplayName:   "Adjust Inventory Level",
		Description:   "Automatically updates the inventory level of a product in your system by adjusting the quantity available based on sales, returns, or other relevant factors.",
		Type:          core.ActionTypeAction,
		Documentation: adjustInventoryLevelDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *AdjustInventoryLevelAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("adjust_inventory_level", "Adjust Inventory Level")

	form.NumberField("inventoryItemId", "Inventory Item").
		Placeholder("The ID of the inventory item.").
		Required(true).
		HelpText("The ID of the inventory item.")

	form.NumberField("locationId", "Location").
		Placeholder("The ID of the location.").
		Required(true).
		HelpText("The ID of the location.")

	form.NumberField("available_adjustment", "Adjustment Quantity").
		Placeholder("Positive values increase inventory, negative values decrease it.").
		Required(true).
		HelpText("Positive values increase inventory, negative values decrease it.")

	schema := form.Build()
	return schema
}

func (a *AdjustInventoryLevelAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[adjustInventoryLevelActionProps](ctx)
	if err != nil {
		return nil, err
	}

	client, err := shared.CreateClient(ctx)
	if err != nil {
		return nil, err
	}

	options := shopify.InventoryLevelAdjustOptions{
		InventoryItemId: input.InventoryItemID,
		LocationId:      input.LocationID,
		Adjust:          input.AvailableAdjustment,
	}

	inventoryLevel, err := client.InventoryLevel.Adjust(context.Background(), options)
	if err != nil {
		return nil, fmt.Errorf("failed to adjust inventory: %v", err)
	}

	return core.JSON(map[string]interface{}{
		"adjustedInventoryLevel": inventoryLevel,
	}), nil
}

func (a *AdjustInventoryLevelAction) Auth() *core.AuthMetadata {
	return nil
}

func NewAdjustInventoryLevelAction() sdk.Action {
	return &AdjustInventoryLevelAction{}
}
