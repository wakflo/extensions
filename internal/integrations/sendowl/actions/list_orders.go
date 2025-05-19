// actions/list_orders.go
package actions

import (
	"fmt"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/sendowl/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
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

func (a *ListOrdersAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "list_orders",
		DisplayName:   "List Orders",
		Description:   "Retrieve a list of orders from your SendOwl account with optional filtering by date range, status, and other criteria.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: listOrdersDocs,
		SampleOutput: map[string]any{
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
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *ListOrdersAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("list_order", "List Orders")

	form.TextField("page", "Page").
		Required(false).
		HelpText("Page number for pagination, defaults to 1").
		Placeholder("1").
		DefaultValue("1")

	form.TextField("limit", "Limit").
		Required(false).
		HelpText("Number of orders per page, defaults to 50").
		Placeholder("50").
		DefaultValue("50")

	form.TextField("status", "Order Status").
		Required(false).
		HelpText("Filter orders by status").
		Placeholder("select status").
		AddOptions([]*smartform.Option{
			{
				Value: "",
				Label: "All",
			},
			{
				Value: "completed",
				Label: "Completed",
			},
			{
				Value: "pending",
				Label: "Pending",
			},
			{
				Value: "refunded",
				Label: "Refunded",
			},
			{
				Value: "cancelled",
				Label: "Cancelled",
			},
		}...)

	form.DateField("start_date", "Start Date").
		Required(false).
		HelpText("Filter orders from this date (format: YYYY-MM-DD)").
		Placeholder("YYYY-MM-DD")

	form.DateField("end_date", "End Date").
		Required(false).
		HelpText("Filter orders until this date (format: YYYY-MM-DD)").
		Placeholder("YYYY-MM-DD")

	shared.GetProductProp("product_id", "Product ID", "Filter orders by specific product ID", true, form)

	form.TextField("buyer_email", "Buyer Email").
		Required(false).
		HelpText("Filter orders by buyer's email address").
		Placeholder("buyer@example.com")

	schema := form.Build()

	return schema
}

func (a *ListOrdersAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[listOrdersActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
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

	response, err := shared.GetSendOwlClient(shared.AltBaseURL, authCtx.Extra["api_key"], authCtx.Extra["api_secret"], url)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (a *ListOrdersAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewListOrdersAction() sdk.Action {
	return &ListOrdersAction{}
}
