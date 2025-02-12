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

type createCollectActionProps struct {
	ProductID    uint64 `json:"productId"`
	CollectionID uint64 `json:"collectionId"`
}

type CreateCollectAction struct{}

func (a *CreateCollectAction) Name() string {
	return "Create Collect"
}

func (a *CreateCollectAction) Description() string {
	return "Create Collect: Automatically generates and sends customizable collection requests to stakeholders, streamlining data gathering and improving collaboration across teams."
}

func (a *CreateCollectAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *CreateCollectAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createCollectDocs,
	}
}

func (a *CreateCollectAction) Icon() *string {
	return nil
}

func (a *CreateCollectAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
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
	}
}

func (a *CreateCollectAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createCollectActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	client, err := shared.CreateClient(ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	collect := shopify.Collect{
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

func (a *CreateCollectAction) Auth() *sdk.Auth {
	return nil
}

func (a *CreateCollectAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *CreateCollectAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreateCollectAction() sdk.Action {
	return &CreateCollectAction{}
}
