package actions

import (
	"net/http"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/prisync/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type listProductsActionProps struct{}

type ListProductsAction struct{}

// Metadata returns metadata about the action
func (a *ListProductsAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "list_products",
		DisplayName:   "List Products",
		Description:   "Retrieves a list of products from a specified data source or API, allowing you to automate tasks that require product information, such as updating inventory levels or sending notifications.",
		Type:          core.ActionTypeAction,
		Documentation: listProductsDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *ListProductsAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("list_products", "List Products")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *ListProductsAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *ListProductsAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	_, err := sdk.InputToTypeSafely[listProductsActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	endpoint := "/api/v2/list/product/startFrom/0"
	resp, err := shared.PrisyncRequest(authCtx.Extra["api-key"], authCtx.Extra["api-token"], endpoint, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NewListProductsAction() sdk.Action {
	return &ListProductsAction{}
}
