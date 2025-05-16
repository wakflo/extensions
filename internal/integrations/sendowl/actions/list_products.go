package actions

import (
	"fmt"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/sendowl/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type listProductsActionProps struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type ListProductsAction struct{}

func (a *ListProductsAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		DisplayName:   "List Products",
		Description:   "Retrieve a comprehensive list of all products from your SendOwl account, allowing you to access product information for reporting, analysis, or to trigger subsequent workflow steps.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: listProductsDocs,
		SampleOutput: map[string]any{
			"Count": 1,
			"Data": []map[string]any{
				{
					"ID":          "123456",
					"Name":        "Product Name",
					"Description": "Product Description",
					"Price":       10.0,
					"CreatedAt":   "2023-01-01T00:00:00Z",
				},
			},
			"Total": 1,
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *ListProductsAction) Properties() *smartform.FormSchema {

	form := smartform.NewForm("list_products", "List Products")

	form.NumberField("page", "Page").
		Placeholder("1").
		Required(false).
		HelpText("Page number for pagination, defaults to 1")

	form.NumberField("limit", "Limit").
		Placeholder("50").
		Required(false).
		HelpText("Number of products per page, defaults to 50")

	schema := form.Build()

	return schema
}

func (a *ListProductsAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[listProductsActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
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

	response, err := shared.GetSendOwlClient(shared.BaseURL, authCtx.Extra["api_key"], authCtx.Extra["api_secret"], url)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (a *ListProductsAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewListProductsAction() sdk.Action {
	return &ListProductsAction{}
}
