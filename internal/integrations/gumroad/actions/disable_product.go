package actions

import (
	"errors"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/gumroad/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type disableProductActionProps struct {
	ProductID string `json:"product_id"`
}

type DisableProductAction struct{}

// Metadata returns metadata about the action
func (a *DisableProductAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "disable_product",
		DisplayName:   "Disable Product",
		Description:   "Disable a product in your Gumroad store.",
		Type:          core.ActionTypeAction,
		Documentation: getProductDocs,
		Icon:          "mdi:package-variant-closed-multiple",
		SampleOutput: map[string]any{
			"success": true,
			"message": "Product successfully disabled.",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *DisableProductAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("disable_product", "Disable Product")

	// Register products selection field
	shared.RegisterProductsProps(form)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *DisableProductAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *DisableProductAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[disableProductActionProps](ctx)
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

	product, err := shared.DisableProduct(accessToken, input.ProductID)

	return product, nil
}

func NewDisableProductAction() sdk.Action {
	return &DisableProductAction{}
}
