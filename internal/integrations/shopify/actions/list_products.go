package actions

import (
	"context"
	"errors"

	"github.com/wakflo/extensions/internal/integrations/shopify/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type listProductsActionProps struct {
	Name string `json:"name"`
}

type ListProductsAction struct{}

func (a *ListProductsAction) Name() string {
	return "List Products"
}

func (a *ListProductsAction) Description() string {
	return "Retrieves a list of products from a specified data source or API, allowing you to automate tasks that require product information, such as updating inventory levels or sending notifications."
}

func (a *ListProductsAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *ListProductsAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &listProductsDocs,
	}
}

func (a *ListProductsAction) Icon() *string {
	return nil
}

func (a *ListProductsAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"name": autoform.NewShortTextField().
			SetLabel("Name").
			SetRequired(true).
			SetPlaceholder("Your name").
			Build(),
	}
}

func (a *ListProductsAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	client, err := shared.CreateClient(ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	products, err := client.Product.List(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	if products == nil {
		return nil, errors.New("no products found")
	}

	return sdk.JSON(map[string]interface{}{
		"Total count of products": products,
	}), err
}

func (a *ListProductsAction) Auth() *sdk.Auth {
	return nil
}

func (a *ListProductsAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *ListProductsAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewListProductsAction() sdk.Action {
	return &ListProductsAction{}
}
