package actions

import (
	"context"

	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/shopspring/decimal"
	"github.com/wakflo/extensions/internal/integrations/shopify/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type createOrderActionProps struct {
	ProductID  uint64               `json:"productId"`
	LineItems  []goshopify.LineItem `json:"line_items"`
	VariantID  uint64               `json:"variantId"`
	CustomerID uint64               `json:"customerId"`
	Title      string               `json:"title"`
	Tags       string               `json:"tags"`
	Note       string               `json:"note"`
	Quantity   int                  `json:"quantity"`
	Price      *decimal.Decimal     `json:"price"`
}

type CreateOrderAction struct{}

func (a *CreateOrderAction) Name() string {
	return "Create Order"
}

func (a *CreateOrderAction) Description() string {
	return "Create Order: Automatically generates and submits a new order to your e-commerce platform or inventory management system, streamlining the ordering process and reducing manual errors."
}

func (a *CreateOrderAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *CreateOrderAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createOrderDocs,
	}
}

func (a *CreateOrderAction) Icon() *string {
	return nil
}

func (a *CreateOrderAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"productId": autoform.NewNumberField().
			SetDisplayName("Product ID").
			SetDescription("The ID of the product to create the order with.").
			SetRequired(false).
			Build(),
		"variantId": autoform.NewNumberField().
			SetDisplayName("Product Variant").
			SetDescription("The ID of the variant to create the order with.").
			SetRequired(false).
			Build(),
		"customerId": autoform.NewNumberField().
			SetDisplayName("Customer ID").
			SetDescription("The ID of the customer to use.").
			SetRequired(false).
			Build(),
		"title": autoform.NewShortTextField().
			SetDisplayName("Title").
			SetRequired(false).
			Build(),
		"note": autoform.NewLongTextField().
			SetDisplayName("Note about the order.").
			SetRequired(false).
			Build(),
		"tags": autoform.NewLongTextField().
			SetDisplayName("Tags").
			SetDescription("A string of comma-separated tags for filtering and search").
			SetRequired(false).
			Build(),
		"quantity": autoform.NewNumberField().
			SetDisplayName("Quantity").
			SetDescription("The ID of the variant to create the order with.").
			SetRequired(false).
			SetDefaultValue(1).
			Build(),
		"price": autoform.NewShortTextField().
			SetDisplayName("Price").
			SetRequired(false).
			Build(),
	}
}

func (a *CreateOrderAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createOrderActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	client, err := shared.CreateClient(ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	newOrder := goshopify.Order{
		LineItems: []goshopify.LineItem{
			{
				ProductId: input.ProductID,
				VariantId: input.VariantID,
				Quantity:  input.Quantity,
				Price:     input.Price,
				Title:     input.Title,
			},
		},
		Note: input.Note,
		Tags: input.Tags,
		Customer: &goshopify.Customer{
			Id: input.CustomerID,
		},
	}
	order, err := client.Order.Create(context.Background(), newOrder)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"new_order": order,
	}, nil
}

func (a *CreateOrderAction) Auth() *sdk.Auth {
	return nil
}

func (a *CreateOrderAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *CreateOrderAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreateOrderAction() sdk.Action {
	return &CreateOrderAction{}
}
