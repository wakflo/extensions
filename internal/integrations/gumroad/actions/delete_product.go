package actions

import (
	"errors"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/gumroad/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type deleteProductActionProps struct {
	ProductID string `json:"product_id"`
}

type DeleteProductAction struct{}

// Metadata returns metadata about the action
func (a *DeleteProductAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "delete_product",
		DisplayName:   "Delete Product",
		Description:   "Permanently a product from your Gumroad store.",
		Type:          core.ActionTypeAction,
		Documentation: deleteProductDocs,
		Icon:          "mdi:package-variant-closed-multiple",
		SampleOutput: map[string]any{
			"success": true,
			"message": "Product successfully deleted.",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *DeleteProductAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("delete_product", "Delete Product")

	shared.RegisterProductsProps(form)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *DeleteProductAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *DeleteProductAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[deleteProductActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	if authCtx.AccessToken == "" {
		return nil, errors.New("missing Gumroad auth token")
	}
	accessToken := authCtx.Token.AccessToken

	if input.ProductID == "" {
		return nil, errors.New("product ID is required")
	}

	product, err := shared.DeleteProduct(accessToken, input.ProductID)

	return product, nil
}

func NewDeleteProductAction() sdk.Action {
	return &DeleteProductAction{}
}
