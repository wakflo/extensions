package actions

import (
	"net/http"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/prisync/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type getProductActionProps struct {
	ProductID string `json:"id"`
}

type GetProductAction struct{}

// Metadata returns metadata about the action
func (a *GetProductAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_product",
		DisplayName:   "Get Product",
		Description:   "Retrieves product information from the specified source, such as an e-commerce platform or inventory management system.",
		Type:          core.ActionTypeAction,
		Documentation: getProductDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *GetProductAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_product", "Get Product")

	form.TextField("id", "Product ID").
		Required(true)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *GetProductAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *GetProductAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[getProductActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	endpoint := "/api/v2/get/product/" + input.ProductID
	resp, err := shared.PrisyncRequest(authCtx.Extra["api-key"], authCtx.Extra["api-token"], endpoint, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NewGetProductAction() sdk.Action {
	return &GetProductAction{}
}
