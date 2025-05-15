package actions

import (
	"errors"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/gumroad/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type getSaleActionProps struct {
	SaleID string `json:"sale_id"`
}

type GetSaleAction struct{}

// Metadata returns metadata about the action
func (a *GetSaleAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_sale",
		DisplayName:   "Get Sale",
		Description:   "Retrieves a sale from your Gumroad store.",
		Type:          core.ActionTypeAction,
		Documentation: getSaleDocs,
		Icon:          "mdi:package-variant-closed-multiple",
		SampleOutput: map[string]any{
			"sale": map[string]any{
				"id":             "xyz789",
				"product_id":     "abc123",
				"product_name":   "Digital Product",
				"price":          "19.99",
				"email":          "customer@example.com",
				"full_name":      "John Doe",
				"currency":       "usd",
				"order_number":   "10051",
				"sale_timestamp": "2023-05-15T14:30:45Z",
				"refunded":       false,
				"variants": map[string]string{
					"Package": "Standard",
				},
				"license_key": "GUMROAD-XYZ123-ABC789",
			},
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *GetSaleAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_sale", "Get Sale")

	// Register sales selection field
	shared.RegisterSalesProps(form)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *GetSaleAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *GetSaleAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[getSaleActionProps](ctx)
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

	if input.SaleID == "" {
		return nil, errors.New("sale ID is required")
	}

	sale, err := shared.GetSale(accessToken, input.SaleID)

	return sale, nil
}

func NewGetSaleAction() sdk.Action {
	return &GetSaleAction{}
}
