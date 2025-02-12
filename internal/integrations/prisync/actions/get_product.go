package actions

import (
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/prisync/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type getProductActionProps struct {
	ProductID string `json:"id"`
}

type GetProductAction struct{}

func (a *GetProductAction) Name() string {
	return "Get Product"
}

func (a *GetProductAction) Description() string {
	return "Retrieves product information from the specified source, such as an e-commerce platform or inventory management system."
}

func (a *GetProductAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *GetProductAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getProductDocs,
	}
}

func (a *GetProductAction) Icon() *string {
	return nil
}

func (a *GetProductAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"id": autoform.NewShortTextField().
			SetDisplayName("Product ID").
			SetRequired(true).Build(),
	}
}

func (a *GetProductAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getProductActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	endpoint := "/api/v2/get/product/" + input.ProductID
	resp, err := shared.PrisyncRequest(ctx.Auth.Extra["api-key"], ctx.Auth.Extra["api-token"], endpoint, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (a *GetProductAction) Auth() *sdk.Auth {
	return nil
}

func (a *GetProductAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *GetProductAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewGetProductAction() sdk.Action {
	return &GetProductAction{}
}
