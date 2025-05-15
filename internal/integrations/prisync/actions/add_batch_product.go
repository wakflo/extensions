package actions

import (
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/prisync/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type addBatchProductActionProps struct {
	Products []ProductInput `json:"products"`
}

type ProductInput struct {
	Name           string `json:"name"`
	Brand          string `json:"brand"`
	Category       string `json:"category"`
	ProductCode    string `json:"product_code"`
	Barcode        string `json:"barcode"`
	Cost           string `json:"cost"`
	AdditionalCost string `json:"additional_cost"`
}

type AddBatchProductAction struct{}

// Metadata returns metadata about the action
func (a *AddBatchProductAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "add_batch_product",
		DisplayName:   "Add Batch Product",
		Description:   "Add Batch Product: Automatically adds a new product to your batch, allowing you to manage and track multiple products within a single batch. This integration action enables seamless product management, streamlining your workflow and reducing manual errors.",
		Type:          core.ActionTypeAction,
		Documentation: addBatchProductDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *AddBatchProductAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("add_batch_product", "Add Batch Product")

	// Create a products array field
	productsArray := form.ArrayField("products", "Products")

	// Configure the product item template
	productGroup := productsArray.ObjectTemplate("product", "")

	// Add all the required fields to the product template
	productGroup.TextField("name", "Product Name").
		Placeholder("Enter product name").
		Required(true).
		HelpText("The name of the product")

	productGroup.TextField("brand", "Brand").
		Placeholder("Enter brand name").
		Required(true).
		HelpText("The brand of the product")

	productGroup.TextField("category", "Category").
		Placeholder("Enter product category").
		Required(true).
		HelpText("The category of the product")

	productGroup.TextField("product_code", "Product Code").
		Placeholder("Enter product code").
		Required(false).
		HelpText("The unique code for the product")

	productGroup.TextField("barcode", "Barcode").
		Placeholder("Enter barcode").
		Required(false).
		HelpText("The barcode of the product")

	productGroup.TextField("cost", "Cost").
		Placeholder("Enter cost").
		Required(false).
		HelpText("The cost of the product")

	productGroup.TextField("additional_cost", "Additional Cost").
		Placeholder("Enter additional cost").
		Required(false).
		HelpText("Any additional costs associated with the product")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *AddBatchProductAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *AddBatchProductAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[addBatchProductActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	formData := make([]map[string]string, 0, len(input.Products))

	for _, product := range input.Products {
		formData = append(formData, map[string]string{
			"name":            product.Name,
			"brand":           product.Brand,
			"category":        product.Category,
			"product_code":    product.ProductCode,
			"barcode":         product.Barcode,
			"cost":            product.Cost,
			"additional_cost": product.AdditionalCost,
		})
	}

	endpoint := "/api/v2/add/batch/"
	resp, err := shared.PrisyncBatchRequest(authCtx.Extra["api-key"], authCtx.Extra["api-token"], endpoint, formData, false)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NewAddBatchProductAction() sdk.Action {
	return &AddBatchProductAction{}
}
