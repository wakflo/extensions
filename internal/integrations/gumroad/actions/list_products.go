package actions

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/wakflo/extensions/internal/integrations/gumroad/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type listProductsActionProps struct {
	Published *bool `json:"published"`
	Limit     *int  `json:"limit"`
	Page      *int  `json:"page"`
	// Category  *string `json:"category"`
}

type ListProductsAction struct{}

func (a *ListProductsAction) Name() string {
	return "List Products"
}

func (a *ListProductsAction) Description() string {
	return "Retrieves a list of all products in your Gumroad store."
}

func (a *ListProductsAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *ListProductsAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &listProductsDocs,
	}
}

func (a *ListProductsAction) Icon() *string {
	icon := "mdi:package-variant-closed-multiple"
	return &icon
}

func (a *ListProductsAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"published": autoform.NewBooleanField().
			SetDisplayName("Published Only").
			SetDescription("If true, only returns published products.").
			SetRequired(false).Build(),
		"limit": autoform.NewNumberField().
			SetDisplayName("Limit").
			SetDescription("Maximum number of products to return (default: 10).").
			SetRequired(false).Build(),
		"page": autoform.NewNumberField().
			SetDisplayName("Page").
			SetDescription("Page number for pagination (default: 1).").
			SetRequired(false).Build(),
		// "category": autoform.NewShortTextField().
		// 	SetDisplayName("Category").
		// 	SetDescription("Filter products by category.").
		// 	SetRequired(false).Build(),
	}
}

func (a *ListProductsAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[listProductsActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing calendly auth token")
	}
	accessToken := ctx.Auth.AccessToken

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

	// if input.Category != nil && *input.Category != "" {
	// 	params.Set("category", *input.Category)
	// }

	products, err := shared.ListProducts(accessToken, params)

	return products, nil
}

func (a *ListProductsAction) Auth() *sdk.Auth {
	return nil
}

func (a *ListProductsAction) SampleData() sdkcore.JSON {
	return nil
}

func (a *ListProductsAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewListProductsAction() sdk.Action {
	return &ListProductsAction{}
}
