package actions

import (
	"context"
	"fmt"

	shopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/wakflo/extensions/internal/integrations/shopify/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type adjustInventoryLevelActionProps struct {
	InventoryItemID     uint64 `json:"inventoryItemId"`
	LocationID          uint64 `json:"locationId"`
	AvailableAdjustment int    `json:"available_adjustment"`
}

type AdjustInventoryLevelAction struct{}

func (a *AdjustInventoryLevelAction) Name() string {
	return "Adjust Inventory Level"
}

func (a *AdjustInventoryLevelAction) Description() string {
	return "Automatically updates the inventory level of a product in your system by adjusting the quantity available based on sales, returns, or other relevant factors."
}

func (a *AdjustInventoryLevelAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *AdjustInventoryLevelAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &adjustInventoryLevelDocs,
	}
}

func (a *AdjustInventoryLevelAction) Icon() *string {
	return nil
}

func (a *AdjustInventoryLevelAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
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
	}
}

func (a *AdjustInventoryLevelAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[adjustInventoryLevelActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	client, err := shared.CreateClient(ctx.BaseContext)
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

	return sdk.JSON(map[string]interface{}{
		"adjustedInventoryLevel": inventoryLevel,
	}), nil
}

func (a *AdjustInventoryLevelAction) Auth() *sdk.Auth {
	return nil
}

func (a *AdjustInventoryLevelAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *AdjustInventoryLevelAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewAdjustInventoryLevelAction() sdk.Action {
	return &AdjustInventoryLevelAction{}
}
