// actions/list_orders.go
package actions

import (
	"fmt"

	"github.com/wakflo/extensions/internal/integrations/sendowl/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type listOrdersActionProps struct {
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	Status     string `json:"status"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	ProductID  string `json:"product_id"`
	BuyerEmail string `json:"buyer_email"`
}

type ListOrdersAction struct{}

func (a *ListOrdersAction) Name() string {
	return "List Orders"
}

func (a *ListOrdersAction) Description() string {
	return "Retrieve a list of orders from your SendOwl account with optional filtering by date range, status, and other criteria."
}

func (a *ListOrdersAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *ListOrdersAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &listOrdersDocs,
	}
}

func (a *ListOrdersAction) Icon() *string {
	icon := "mdi:receipt"
	return &icon
}

func (a *ListOrdersAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"page": autoform.NewNumberField().
			SetDisplayName("Page").
			SetDescription("Page number for pagination, defaults to 1").
			SetDefaultValue(1).
			Build(),
		"limit": autoform.NewNumberField().
			SetDisplayName("Limit").
			SetDescription("Number of orders per page, defaults to 50").
			SetDefaultValue(50).
			Build(),
		"status": autoform.NewSelectField().
			SetDisplayName("Order Status").
			SetDescription("Filter orders by status").
			SetOptions([]*sdkcore.AutoFormSchema{
				autoform.NewShortTextField().
					SetDisplayName("All").
					SetDefaultValue("").
					Build(),
				autoform.NewShortTextField().
					SetDisplayName("Completed").
					SetDefaultValue("completed").
					Build(),
				autoform.NewShortTextField().
					SetDisplayName("Pending").
					SetDefaultValue("pending").
					Build(),
				autoform.NewShortTextField().
					SetDisplayName("Refunded").
					SetDefaultValue("refunded").
					Build(),
				autoform.NewShortTextField().
					SetDisplayName("Cancelled").
					SetDefaultValue("cancelled").
					Build(),
			}).
			Build(),
		"start_date": autoform.NewShortTextField().
			SetDisplayName("Start Date").
			SetDescription("Filter orders from this date (format: YYYY-MM-DD)").
			Build(),
		"end_date": autoform.NewShortTextField().
			SetDisplayName("End Date").
			SetDescription("Filter orders until this date (format: YYYY-MM-DD)").
			Build(),
		"product_id": shared.GetProductInput("Product ID", "Filter orders by specific product ID", true),
		"buyer_email": autoform.NewShortTextField().
			SetDisplayName("Buyer Email").
			SetDescription("Filter orders by buyer's email address").
			Build(),
	}
}

func (a *ListOrdersAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[listOrdersActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	url := "/orders"

	// Construct query parameters
	queryParams := make([]string, 0)

	if input.Page > 0 {
		queryParams = append(queryParams, fmt.Sprintf("page=%d", input.Page))
	} else {
		queryParams = append(queryParams, "page=1")
	}

	if input.Limit > 0 {
		queryParams = append(queryParams, fmt.Sprintf("per_page=%d", input.Limit))
	} else {
		queryParams = append(queryParams, "per_page=50")
	}

	if input.Status != "" {
		queryParams = append(queryParams, fmt.Sprintf("status=%s", input.Status))
	}

	if input.StartDate != "" {
		queryParams = append(queryParams, fmt.Sprintf("from_date=%s", input.StartDate))
	}

	if input.EndDate != "" {
		queryParams = append(queryParams, fmt.Sprintf("to_date=%s", input.EndDate))
	}

	if input.ProductID != "" {
		queryParams = append(queryParams, fmt.Sprintf("product_id=%s", input.ProductID))
	}

	if input.BuyerEmail != "" {
		queryParams = append(queryParams, fmt.Sprintf("buyer_email=%s", input.BuyerEmail))
	}

	// Add query parameters to URL
	if len(queryParams) > 0 {
		url = url + "?" + queryParams[0]
		for i := 1; i < len(queryParams); i++ {
			url = url + "&" + queryParams[i]
		}
	}

	response, err := shared.GetSendOwlClient(shared.AltBaseURL, ctx.Auth.Extra["api_key"], ctx.Auth.Extra["api_secret"], url)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (a *ListOrdersAction) Auth() *sdk.Auth {
	return nil
}

func (a *ListOrdersAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"orders": []map[string]any{
			{
				"id":             123456,
				"buyer_name":     "John Doe",
				"buyer_email":    "john.doe@example.com",
				"total":          29.99,
				"currency":       "USD",
				"status":         "completed",
				"payment_method": "credit_card",
				"created_at":     "2023-08-15T14:30:22Z",
				"completed_at":   "2023-08-15T14:35:45Z",
				"products": []map[string]any{
					{
						"id":    789012,
						"name":  "Digital Marketing Guide",
						"price": 29.99,
					},
				},
			},
			{
				"id":             789012,
				"buyer_name":     "Jane Smith",
				"buyer_email":    "jane.smith@example.com",
				"total":          99.00,
				"currency":       "USD",
				"status":         "completed",
				"payment_method": "paypal",
				"created_at":     "2023-08-16T09:12:35Z",
				"completed_at":   "2023-08-16T09:15:10Z",
				"products": []map[string]any{
					{
						"id":    345678,
						"name":  "SEO Masterclass",
						"price": 99.00,
					},
				},
			},
		},
	}
}

func (a *ListOrdersAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewListOrdersAction() sdk.Action {
	return &ListOrdersAction{}
}
