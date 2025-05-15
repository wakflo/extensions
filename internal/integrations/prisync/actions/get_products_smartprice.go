package actions

import (
	"net/http"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/prisync/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type getProductsSmartpriceActionProps struct{}

type GetProductsSmartpriceAction struct{}

// Metadata returns metadata about the action
func (a *GetProductsSmartpriceAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_products_smartprice",
		DisplayName:   "Get Products Smartprice",
		Description:   "Retrieves product prices from various e-commerce platforms and marketplaces, providing real-time smart pricing data to inform business decisions.",
		Type:          core.ActionTypeAction,
		Documentation: getProductsSmartpriceDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *GetProductsSmartpriceAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_products_smartprice", "Get Products Smartprice")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *GetProductsSmartpriceAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *GetProductsSmartpriceAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	_, err := sdk.InputToTypeSafely[getProductsSmartpriceActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	endpoint := "/api/v2/list/smartprice/startFrom/0"
	resp, err := shared.PrisyncRequest(authCtx.Extra["api-key"], authCtx.Extra["api-token"], endpoint, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NewGetProductsSmartpriceAction() sdk.Action {
	return &GetProductsSmartpriceAction{}
}
