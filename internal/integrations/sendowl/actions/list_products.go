package actions

import (
	"fmt"

	"github.com/wakflo/extensions/internal/integrations/sendowl/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type listProductsActionProps struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type ListProductsAction struct{}

func (a *ListProductsAction) Name() string {
	return "List Products"
}

func (a *ListProductsAction) Description() string {
	return "Retrieve a comprehensive list of all products from your SendOwl account, allowing you to access product information for reporting, analysis, or to trigger subsequent workflow steps."
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
	icon := "mdi:cart-outline"
	return &icon
}

func (a *ListProductsAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"page": autoform.NewNumberField().
			SetDisplayName("Page").
			SetDescription("Page number for pagination, defaults to 1").
			SetDefaultValue(1).
			Build(),
		"limit": autoform.NewNumberField().
			SetDisplayName("Limit").
			SetDescription("Number of products per page, defaults to 50").
			SetDefaultValue(50).
			Build(),
	}
}

func (a *ListProductsAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[listProductsActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	url := "/products"

	// Add pagination parameters
	if input.Page > 0 || input.Limit > 0 {
		page := 1
		limit := 50

		if input.Page > 0 {
			page = input.Page
		}

		if input.Limit > 0 {
			limit = input.Limit
		}

		url = fmt.Sprintf("%s?page=%d&per_page=%d", url, page, limit)
	}

	response, err := shared.GetSendOwlClient(ctx.Auth.Extra["api_key"], ctx.Auth.Extra["api_secret"], url)
	if err != nil {
		return nil, err
	}

	fmt.Println("Response from SendOwl-------------------------->:", response)

	return response, nil
}

func (a *ListProductsAction) Auth() *sdk.Auth {
	return nil
}

func (a *ListProductsAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"products": []map[string]any{
			{
				"id":           123456,
				"name":         "Digital Marketing Guide",
				"price":        29.99,
				"currency":     "USD",
				"product_type": "pdf",
				"sales_count":  156,
				"created_at":   "2023-07-15T10:30:00Z",
				"updated_at":   "2023-08-20T14:45:22Z",
			},
			{
				"id":           789012,
				"name":         "SEO Masterclass",
				"price":        99.00,
				"currency":     "USD",
				"product_type": "video",
				"sales_count":  87,
				"created_at":   "2023-06-10T09:15:30Z",
				"updated_at":   "2023-08-18T11:20:15Z",
			},
		},
	}
}

func (a *ListProductsAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewListProductsAction() sdk.Action {
	return &ListProductsAction{}
}
