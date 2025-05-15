package actions

import (
	"net/http"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/prisync/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type addProductActionProps struct {
	Name        string `json:"name"`
	Brand       string `json:"brand"`
	Category    string `json:"category"`
	Cost        string `json:"cost"`
	ProductCode string `json:"product_code"`
	BarCode     string `json:"barcode"`
	Tags        string `json:"tags"`
}

type AddProductAction struct{}

// Metadata returns metadata about the action
func (a *AddProductAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "add_product",
		DisplayName:   "Add Product",
		Description:   "Adds a new product to your inventory or database, allowing you to track and manage products efficiently. This integration action can be used to populate product information from various sources, such as e-commerce platforms, marketplaces, or product information management systems.",
		Type:          core.ActionTypeAction,
		Documentation: addProductDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *AddProductAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("add_product", "Add Product")

	form.TextField("name", "Product Name").
		Required(true).
		HelpText("name of product")

	form.TextField("brand", "Brand").
		Required(true).
		HelpText("Brand name")

	form.TextField("category", "Category").
		Required(true).
		HelpText("Category name")

	form.TextField("product_code", "Product Code").
		Required(false).
		HelpText("Product code")

	form.TextField("barcode", "Bar Code").
		Required(false).
		HelpText("Bar code")

	form.TextField("cost", "Product Cost").
		Required(false).
		HelpText("Product cost")

	form.TextField("tags", "Product Tags").
		Required(false).
		HelpText("Product tags")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *AddProductAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *AddProductAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[addProductActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	formData := map[string]string{
		"name":     input.Name,
		"brand":    input.Brand,
		"category": input.Category,
	}

	if input.ProductCode != "" {
		formData["product_code"] = input.ProductCode
	}

	if input.Cost != "" {
		formData["cost"] = input.Cost
	}

	if input.BarCode != "" {
		formData["barcode"] = input.BarCode
	}

	if input.Tags != "" {
		formData["tags"] = input.Tags
	}

	endpoint := "/api/v2/add/product/"
	resp, err := shared.PrisyncRequest(authCtx.Extra["api-key"], authCtx.Extra["api-token"], endpoint, http.MethodPost, formData)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NewAddProductAction() sdk.Action {
	return &AddProductAction{}
}
