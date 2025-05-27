package actions

import (
	"net/http"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/prisync/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type editProductActionProps struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Brand       string `json:"brand"`
	Category    string `json:"category"`
	BarCode     string `json:"barcode"`
	Cost        string `json:"cost"`
	ProductCode string `json:"product_code"`
	Tags        string `json:"tags"`
}

type EditProductAction struct{}

// Metadata returns metadata about the action
func (a *EditProductAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "edit_product",
		DisplayName:   "Edit Product",
		Description:   "Edit Product: Update product information by modifying existing product details, such as name, description, price, and inventory levels.",
		Type:          core.ActionTypeAction,
		Documentation: editProductDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *EditProductAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("edit_product", "Edit Product")

	shared.GetProductProp("id", "Product ID", "ID of product", true, form)

	form.TextField("name", "Product Name").
		Required(false).
		HelpText("name of product")

	form.TextField("brand", "Brand").
		Required(false).
		HelpText("Brand name")

	form.TextField("category", "Category").
		Required(false).
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
func (a *EditProductAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *EditProductAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[editProductActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	formData := map[string]string{}

	if input.Name != "" {
		formData["name"] = input.Name
	}

	if input.Brand != "" {
		formData["brand"] = input.Brand
	}

	if input.Category != "" {
		formData["category"] = input.Category
	}

	if input.ProductCode != "" {
		formData["product_code"] = input.ProductCode
	}

	if input.BarCode != "" {
		formData["barcode"] = input.BarCode
	}

	if input.Cost != "" {
		formData["cost"] = input.Cost
	}

	if input.Tags != "" {
		formData["tags"] = input.Tags
	}

	endpoint := "/api/v2/edit/product/id/" + input.ID
	resp, err := shared.PrisyncRequest(authCtx.Extra["api-key"], authCtx.Extra["api-token"], endpoint, http.MethodPost, formData)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NewEditProductAction() sdk.Action {
	return &EditProductAction{}
}
