package actions

import (
	"errors"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/gumroad/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type enableProductActionProps struct {
	ProductID string `json:"product_id"`
}

type EnableProductAction struct{}

// Metadata returns metadata about the action
func (a *EnableProductAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "enable_product",
		DisplayName:   "Enable Product",
		Description:   "Enable a product in your Gumroad store.",
		Type:          core.ActionTypeAction,
		Documentation: getProductDocs,
		Icon:          "mdi:package-variant-closed-multiple",
		SampleOutput: map[string]any{
			"success": true,
			"message": "Product successfully enabled.",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *EnableProductAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("enable_product", "Enable Product")

	shared.RegisterProductsProps(form)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *EnableProductAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *EnableProductAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[enableProductActionProps](ctx)
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
	accessToken := authCtx.AccessToken

	if input.ProductID == "" {
		return nil, errors.New("product ID is required")
	}

	product, err := shared.EnableProduct(accessToken, input.ProductID)

	return product, nil
}

func NewEnableProductAction() sdk.Action {
	return &EnableProductAction{}
}
