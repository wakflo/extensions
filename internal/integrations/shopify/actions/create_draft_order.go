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

func (a *CreateDraftOrderAction) Name() string {
	return "Create Draft Order"
}

func (a *CreateDraftOrderAction) Description() string {
	return "Create Draft Order: Automatically generates a draft order in your e-commerce platform, allowing you to review and customize order details before submitting it for fulfillment."
}

func (a *CreateDraftOrderAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *CreateDraftOrderAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createDraftOrderDocs,
	}
}

func (a *CreateDraftOrderAction) Icon() *string {
	return nil
}

func (a *CreateDraftOrderAction) Properties() map[string]*sdkcore.AutoFormSchema {
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
			SetDisplayName("Note about the order").
			SetRequired(false).
			Build(),
		"tags": autoform.NewLongTextField().
			SetDisplayName("A string of comma-separated tags for filtering and search").
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

func (a *CreateDraftOrderAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createDraftOrderActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	client, err := shared.CreateClient(ctx.BaseContext)
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

func (a *CreateDraftOrderAction) Auth() *sdk.Auth {
	return nil
}

func (a *CreateDraftOrderAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *CreateDraftOrderAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreateDraftOrderAction() sdk.Action {
	return &CreateDraftOrderAction{}
}
