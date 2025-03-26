package actions

import (
	"errors"
	"net/url"

	"github.com/wakflo/extensions/internal/integrations/gumroad/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
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

func (a *ListSalesAction) Name() string {
	return "List Sales"
}

func (a *ListSalesAction) Description() string {
	return "Retrieves a list of sales from your Gumroad account with optional filtering."
}

func (a *ListSalesAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *ListSalesAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &listSalesDocs,
	}
}

func (a *ListSalesAction) Icon() *string {
	icon := "mdi:receipt-text"
	return &icon
}

func (a *ListSalesAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"after": autoform.NewDateTimeField().
			SetDisplayName("After").
			SetDescription("Only return sales after this date.").
			SetRequired(false).
			Build(),
		"before": autoform.NewDateTimeField().
			SetDisplayName("Before").
			SetDescription("Only return sales before this date.").
			SetRequired(false).
			Build(),
		"product_id": autoform.NewShortTextField().
			SetDisplayName("Product ID").
			SetDescription("Filter sales by this product.").
			SetRequired(false).
			Build(),
		"email": autoform.NewShortTextField().
			SetDisplayName("Email").
			SetDescription("Filter sales by this email.").
			SetRequired(false).
			Build(),
		"order_id": autoform.NewShortTextField().
			SetDisplayName("Order ID").
			SetDescription("Filter sales by this Order ID.").
			SetRequired(false).
			Build(),
		"page_key": autoform.NewShortTextField().
			SetDisplayName("Page Key").
			SetDescription("A key representing a page of results.").
			SetRequired(false).
			Build(),
	}
}

func (a *ListSalesAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[listSalesActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing calendly auth token")
	}
	accessToken := ctx.Auth.AccessToken

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

func (a *ListSalesAction) Auth() *sdk.Auth {
	return nil
}

func (a *ListSalesAction) SampleData() sdkcore.JSON {
	return nil
}

func (a *ListSalesAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewListSalesAction() sdk.Action {
	return &ListSalesAction{}
}
