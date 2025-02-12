package actions

import (
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/prisync/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type getProductsSmartpriceActionProps struct {
	Startfrom string `json:"start-from"`
}

type GetProductsSmartpriceAction struct{}

func (a *GetProductsSmartpriceAction) Name() string {
	return "Get Products Smartprice"
}

func (a *GetProductsSmartpriceAction) Description() string {
	return "Retrieves product prices from various e-commerce platforms and marketplaces, providing real-time smart pricing data to inform business decisions."
}

func (a *GetProductsSmartpriceAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *GetProductsSmartpriceAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getProductsSmartpriceDocs,
	}
}

func (a *GetProductsSmartpriceAction) Icon() *string {
	return nil
}

func (a *GetProductsSmartpriceAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"start-from": autoform.NewShortTextField().
			SetDisplayName("Start From (Optional)").
			SetDescription("Offset for pagination. It can take 0 and exact multiples of 100 as a value.").
			SetRequired(false).Build(),
	}
}

func (a *GetProductsSmartpriceAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	_, err := sdk.InputToTypeSafely[getProductsSmartpriceActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	endpoint := "/api/v2/list/smartprice/startFrom/0"
	resp, err := shared.PrisyncRequest(ctx.Auth.Extra["api-key"], ctx.Auth.Extra["api-token"], endpoint, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (a *GetProductsSmartpriceAction) Auth() *sdk.Auth {
	return nil
}

func (a *GetProductsSmartpriceAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *GetProductsSmartpriceAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewGetProductsSmartpriceAction() sdk.Action {
	return &GetProductsSmartpriceAction{}
}
