package actions

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/gumroad/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type listProductsActionProps struct {
	Published *bool `json:"published"`
	Limit     *int  `json:"limit"`
	Page      *int  `json:"page"`
}

type ListProductsAction struct{}

// Metadata returns metadata about the action
func (a *ListProductsAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "list_products",
		DisplayName:   "List Products",
		Description:   "Retrieves a list of all products in your Gumroad store.",
		Type:          core.ActionTypeAction,
		Documentation: listProductsDocs,
		Icon:          "mdi:package-variant-closed-multiple",
		SampleOutput: map[string]any{
			"products": []map[string]any{
				{
					"id":               "abc123",
					"name":             "Digital Product",
					"preview_url":      "https://gumroad.com/l/abc123",
					"description":      "A great digital product for your needs",
					"price":            19.99,
					"custom_permalink": "my-digital-product",
					"currency":         "usd",
					"published":        true,
					"sales_count":      152,
				},
			},
			"total_products": 2,
			"has_more":       false,
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *ListProductsAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("list_products", "List Products")

	form.CheckboxField("published", "Published Only").
		Required(false).
		HelpText("If true, only returns published products.")

	form.NumberField("limit", "Limit").
		Required(false).
		HelpText("Maximum number of products to return (default: 10).")

	form.NumberField("page", "Page").
		Required(false).
		HelpText("Page number for pagination (default: 1).")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *ListProductsAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *ListProductsAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[listProductsActionProps](ctx)
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

	// Build query parameters
	params := url.Values{}

	if input.Published != nil {
		if *input.Published {
			params.Set("published", "true")
		} else {
			params.Set("published", "false")
		}
	}

	if input.Limit != nil {
		params.Set("limit", fmt.Sprintf("%d", *input.Limit))
	}

	if input.Page != nil {
		params.Set("page", fmt.Sprintf("%d", *input.Page))
	}

	products, err := shared.ListProducts(accessToken, params)

	return products, nil
}

func NewListProductsAction() sdk.Action {
	return &ListProductsAction{}
}
