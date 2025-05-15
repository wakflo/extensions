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

type createDraftOrderActionProps struct {
	ProductID  uint64               `json:"productId"`
	LineItems  []goshopify.LineItem `json:"line_items"`
	VariantID  uint64               `json:"variantId"`
	CustomerID uint64               `json:"customerId"`
	Note       string               `json:"note"`
	Title      string               `json:"title"`
	Quantity   int                  `json:"quantity"`
	Price      *decimal.Decimal     `json:"price"`
	Tags       string               `json:"tags"`
}

type CreateDraftOrderAction struct{}

func (a *CreateDraftOrderAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_draft_order",
		DisplayName:   "Create Draft Order",
		Description:   "Create Draft Order: Automatically generates a draft order in your e-commerce platform, allowing you to review and customize order details before submitting it for fulfillment.",
		Type:          core.ActionTypeAction,
		Documentation: createDraftOrderDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *CreateDraftOrderAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_draft_order", "Create Draft Order")

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
		Placeholder("Note about the order").
		HelpText("Note about the order")

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

func (a *CreateDraftOrderAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[createDraftOrderActionProps](ctx)
	if err != nil {
		return nil, err
	}

	client, err := shared.CreateClient(ctx)
	if err != nil {
		return nil, err
	}
	newDraftOrder := goshopify.DraftOrder{
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
	}
	order, err := client.DraftOrder.Create(context.Background(), newDraftOrder)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"new_draft_order": order,
	}, nil
}

func (a *CreateDraftOrderAction) Auth() *core.AuthMetadata {
	return nil
}

func NewCreateDraftOrderAction() sdk.Action {
	return &CreateDraftOrderAction{}
}
