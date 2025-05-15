package actions

import (
	"errors"
	"net/url"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/gumroad/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type listSalesActionProps struct {
	After     string `json:"after"`
	Before    string `json:"before"`
	ProductID string `json:"product_id"`
	Email     string `json:"email"`
	OrderID   string `json:"order_id"`
	PageKey   string `json:"page_key"`
}

type ListSalesAction struct{}

// Metadata returns metadata about the action
func (a *ListSalesAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "list_sales",
		DisplayName:   "List Sales",
		Description:   "Retrieves a list of sales from your Gumroad account with optional filtering.",
		Type:          core.ActionTypeAction,
		Documentation: listSalesDocs,
		Icon:          "mdi:receipt-text",
		SampleOutput: map[string]any{
			"sales": []map[string]any{
				{
					"id":             "xyz789",
					"product_id":     "abc123",
					"product_name":   "Digital Product",
					"price":          19.99,
					"email":          "customer1@example.com",
					"full_name":      "John Doe",
					"currency":       "usd",
					"order_number":   "10051",
					"sale_timestamp": "2023-05-15T14:30:45Z",
					"refunded":       false,
				},
				{
					"id":             "uvw456",
					"product_id":     "def456",
					"product_name":   "Creative Template",
					"price":          "29.99",
					"email":          "customer2@example.com",
					"full_name":      "Jane Smith",
					"currency":       "usd",
					"order_number":   "10052",
					"sale_timestamp": "2023-05-16T09:15:30Z",
					"refunded":       false,
				},
			},
			"next_page_key": "next_page_123",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *ListSalesAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("list_sales", "List Sales")

	form.DateTimeField("after", "After").
		Required(false).
		HelpText("Only return sales after this date.")

	form.DateTimeField("before", "Before").
		Required(false).
		HelpText("Only return sales before this date.")

	form.TextField("product_id", "Product ID").
		Required(false).
		HelpText("Filter sales by this product.")

	form.TextField("email", "Email").
		Required(false).
		HelpText("Filter sales by this email.")

	form.TextField("order_id", "Order ID").
		Required(false).
		HelpText("Filter sales by this Order ID.")

	form.TextField("page_key", "Page Key").
		Required(false).
		HelpText("A key representing a page of results.")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *ListSalesAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *ListSalesAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[listSalesActionProps](ctx)
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

	// Build query parameters
	params := url.Values{}

	if input.Before != "" {
		formattedBefore, err := shared.FormatDateInput(input.Before)
		if err != nil {
			return nil, errors.New("invalid before date: " + err.Error())
		}
		params.Set("before", formattedBefore)
	}

	if input.After != "" {
		formattedAfter, err := shared.FormatDateInput(input.After)
		if err != nil {
			return nil, errors.New("invalid after date: " + err.Error())
		}
		params.Set("after", formattedAfter)
	}

	if input.Email != "" {
		params.Set("email", input.Email)
	}

	if input.OrderID != "" {
		params.Set("order_id", input.OrderID)
	}

	if input.PageKey != "" {
		params.Set("page_key", input.PageKey)
	}

	if input.ProductID != "" {
		params.Set("product_id", input.ProductID)
	}

	sales, err := shared.ListSales(accessToken, params)

	return sales, nil
}

func NewListSalesAction() sdk.Action {
	return &ListSalesAction{}
}
