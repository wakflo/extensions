package actions

import (
	"github.com/hiscaler/woocommerce-go"
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/woocommerce/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type ListProductsAction struct{}

func (a *ListProductsAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "list_products",
		DisplayName:   "List Products",
		Description:   "Retrieves a list of products from your WooCommerce store.",
		Type:          core.ActionTypeAction,
		Documentation: listProductsDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *ListProductsAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("list_products", "List Products")

	form.TextareaField("projectId", "Limit").
		Placeholder("").
		HelpText("")

	schema := form.Build()

	return schema
}

func (a *ListProductsAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	wooClient, err := shared.InitClient(ctx)
	if err != nil {
		return nil, err
	}

	params := woocommerce.ProductsQueryParams{}
	products, _, _, _, err := wooClient.Services.Product.All(params)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (a *ListProductsAction) Auth() *core.AuthMetadata {
	return nil
}

func NewListProductsAction() sdk.Action {
	return &ListProductsAction{}
}
