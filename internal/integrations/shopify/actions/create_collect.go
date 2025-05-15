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

type createCollectActionProps struct {
	ProductID    uint64 `json:"productId"`
	CollectionID uint64 `json:"collectionId"`
}

type CreateCollectAction struct{}

func (a *CreateCollectAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_collect",
		DisplayName:   "Create Collect",
		Description:   "Create Collect: Automatically generates and sends customizable collection requests to stakeholders, streamlining data gathering and improving collaboration across teams.",
		Type:          core.ActionTypeAction,
		Documentation: createCollectDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *CreateCollectAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_collect", "Create Collect")

	form.NumberField("productId", "Product ID").
		Placeholder("The ID of the product.").
		Required(true).
		HelpText("The ID of the product.")

	form.NumberField("collectionId", "Collection ID").
		Placeholder("The ID of the collection.").
		Required(true).
		HelpText("The ID of the collection.")

	schema := form.Build()
	return schema
}

func (a *CreateCollectAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[createCollectActionProps](ctx)
	if err != nil {
		return nil, err
	}

	client, err := shared.CreateClient(ctx)
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
	return core.JSON(map[string]interface{}{
		"collection details": newCollect,
	}), nil
}

func (a *CreateCollectAction) Auth() *core.AuthMetadata {
	return nil
}

func NewCreateCollectAction() sdk.Action {
	return &CreateCollectAction{}
}
