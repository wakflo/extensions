package actions

import (
	"github.com/hiscaler/woocommerce-go"
	"github.com/wakflo/extensions/internal/integrations/woocommerce/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)


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
		"projectId": autoform.NewLongTextField().
			SetDisplayName("Limit").
			SetDescription("").
			Build(),
	}
}

func (a *ListProductsAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	wooClient, err := shared.InitClient(ctx.BaseContext)
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
