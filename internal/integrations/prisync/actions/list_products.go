package actions

import (
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/prisync/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type listProductsActionProps struct {
	Startfrom string `json:"start-from"`
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
		"start-from": autoform.NewShortTextField().
			SetDisplayName("Start From (Optional)").
			SetDescription("Offset for pagination. It can take 0 and exact multiples of 100 as a value.").
			SetRequired(false).Build(),
	}
}

func (a *ListProductsAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	_, err := sdk.InputToTypeSafely[listProductsActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	endpoint := "/api/v2/list/product/startFrom/0"
	resp, err := shared.PrisyncRequest(ctx.Auth.Extra["api-key"], ctx.Auth.Extra["api-token"], endpoint, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
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
