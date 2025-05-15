package actions

import (
	"context"
	"fmt"

	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/shopify/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type updateProductActionProps struct {
	ProductID   uint64 `json:"productID"`
	Title       string `json:"title"`
	BodyHTML    string `json:"bodyHTML"`
	Vendor      string `json:"vendor"`
	ProductType string `json:"productType"`
	Status      string `json:"status"`
	Tags        string `json:"tags"`
}

type UpdateProductAction struct{}

// Metadata returns metadata about the action
func (a *UpdateProductAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "update_product",
		DisplayName:   "Update Product",
		Description:   "Updates product information in your e-commerce platform or CRM system by mapping to specific fields such as product name, description, price, and inventory levels.",
		Type:          core.ActionTypeAction,
		Documentation: updateProductDocs,
		SampleOutput: map[string]any{
			"updated_product": map[string]any{},
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *UpdateProductAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("update_product", "Update Product")

	form.NumberField("productID", "Product ID").
		Required(true).
		HelpText("The id of the product.")

	form.TextField("title", "Product title").
		Required(false).
		HelpText("The title of the product.")

	form.TextField("bodyHTML", "Product description").
		Required(false).
		HelpText("The description of the product.")

	form.TextField("vendor", "Vendor").
		Required(false).
		HelpText("Vendor.")

	form.TextField("productType", "Product type").
		Required(false).
		HelpText("A categorization for the product used for filtering and searching products.")

	form.TextField("tags", "Tags").
		Required(false).
		HelpText("A string of comma-separated tags for filtering and search.")

	form.SelectField("status", "Status").
		Required(true).
		AddOption("active", "Active").
		AddOption("draft", "Draft").
		HelpText("The status of the product: active or draft")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *UpdateProductAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *UpdateProductAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[updateProductActionProps](ctx)
	if err != nil {
		return nil, err
	}

	client, err := shared.CreateClient(ctx)
	if err != nil {
		return nil, err
	}

	existingProduct, err := client.Product.Get(context.Background(), input.ProductID, nil)
	if err != nil {
		return nil, err
	}
	if input.Title != "" {
		existingProduct.Title = input.Title
	}
	if input.BodyHTML != "" {
		existingProduct.BodyHTML = input.BodyHTML
	}
	if input.Vendor != "" {
		existingProduct.Vendor = input.Vendor
	}
	if input.ProductType != "" {
		existingProduct.ProductType = input.ProductType
	}
	if input.Tags != "" {
		existingProduct.Tags = input.Tags
	}
	if input.Status != "" {
		existingProduct.Status = goshopify.ProductStatus(input.Status)
	}

	updatedProduct, err := client.Product.Update(context.Background(), *existingProduct)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}
	return map[string]interface{}{
		"updated_product": updatedProduct,
	}, nil
}

func NewUpdateProductAction() sdk.Action {
	return &UpdateProductAction{}
}
