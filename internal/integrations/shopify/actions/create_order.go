package actions

import (
	"context"

	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/juicycleff/smartform/v1"
	"github.com/shopspring/decimal"
	"github.com/wakflo/extensions/internal/integrations/shopify/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
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

func (a *CreateOrderAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_order",
		DisplayName:   "Create Order",
		Description:   "Create Order: Automatically generates and submits a new order to your e-commerce platform or inventory management system, streamlining the ordering process and reducing manual errors.",
		Type:          core.ActionTypeAction,
		Documentation: createOrderDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *CreateOrderAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_order", "Create Order")

	form.NumberField("productId", "Product ID").
		Placeholder("The ID of the product to create the order with.").
		HelpText("The ID of the product to create the order with.")

	form.NumberField("variantId", "Product Variant").
		Placeholder("The ID of the variant to create the order with.").
		HelpText("The ID of the variant to create the order with.")

	form.NumberField("customerId", "Customer ID").
		Placeholder("The ID of the customer to use.").
		HelpText("The ID of the customer to use.")

	form.TextField("title", "Title").
		Placeholder("Title").
		HelpText("Title")

	form.TextareaField("note", "Note about the order").
		Placeholder("Note about the order.").
		HelpText("Note about the order.")

	form.TextareaField("tags", "Tags").
		Placeholder("A string of comma-separated tags for filtering and search").
		HelpText("A string of comma-separated tags for filtering and search")

	form.NumberField("quantity", "Quantity").
		Placeholder("The ID of the variant to create the order with.").
		DefaultValue(1).
		HelpText("The ID of the variant to create the order with.")

	form.TextField("price", "Price").
		Placeholder("Price").
		HelpText("Price")

	schema := form.Build()
	return schema
}

func (a *CreateOrderAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[createOrderActionProps](ctx)
	if err != nil {
		return nil, err
	}

	client, err := shared.CreateClient(ctx)
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

func (a *CreateOrderAction) Auth() *core.AuthMetadata {
	return nil
}

func NewCreateOrderAction() sdk.Action {
	return &CreateOrderAction{}
}
