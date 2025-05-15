package actions

import (
	"errors"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/gumroad/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type getProductActionProps struct {
	ProductID string `json:"product_id"`
}

type GetProductAction struct{}

// Metadata returns metadata about the action
func (a *GetProductAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_product",
		DisplayName:   "Get Product",
		Description:   "Retrieves a product from your Gumroad store.",
		Type:          core.ActionTypeAction,
		Documentation: getProductDocs,
		Icon:          "mdi:package-variant-closed-multiple",
		SampleOutput: map[string]any{
			"product": map[string]any{
				"id":               "abc123",
				"name":             "Digital Product",
				"preview_url":      "https://gumroad.com/l/abc123",
				"description":      "A great digital product for your needs",
				"price":            19.99,
				"custom_permalink": "my-digital-product",
				"currency":         "usd",
				"published":        true,
				"sales_count":      152,
				"custom_fields":    []map[string]any{},
				"custom_summary":   "Digital downloads for creators",
			},
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *GetProductAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_product", "Get Product")

	// Register products selection field
	shared.RegisterProductsProps(form)

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

	if authCtx.AccessToken == "" {
		return nil, errors.New("missing Gumroad auth token")
	}
	accessToken := authCtx.Token.AccessToken

	if input.ProductID == "" {
		return nil, errors.New("product ID is required")
	}

	product, err := shared.GetProduct(accessToken, input.ProductID)

	return product, nil
}

func NewGetProductAction() sdk.Action {
	return &GetProductAction{}
}
