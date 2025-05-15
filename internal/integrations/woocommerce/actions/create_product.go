package actions

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hiscaler/woocommerce-go"
	"github.com/hiscaler/woocommerce-go/entity"
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/woocommerce/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type createProductActionProps struct {
	Name             string  `json:"name"`
	Type             string  `json:"type"`
	RegularPrice     float64 `json:"regular_price"`
	Description      string  `json:"description"`
	ShortDescription string  `json:"short_description"`
	Categories       string  `json:"categories"`
	ImageURL         string  `json:"image_url"`
}

type CreateProductAction struct{}

func (a *CreateProductAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_product",
		DisplayName:   "Create Product",
		Description:   "Create Product: Automatically generates and creates new products in your woocommerce store, including product details such as name, description, price, and inventory levels.",
		Type:          core.ActionTypeAction,
		Documentation: createProductDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *CreateProductAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_product", "Create Product")

	form.TextField("name", "Product Name").
		Placeholder("Enter product Name").
		Required(true).
		HelpText("Enter product Name")

	form.SelectField("type", "Type").
		AddOption("simple", "Simple").
		AddOption("grouped", "Grouped").
		AddOption("external", "External").
		AddOption("variable", "Variable").
		Required(true).
		HelpText("Select the type")

	form.TextareaField("description", "Description").
		Placeholder("Enter product description").
		Required(true).
		HelpText("Enter product description")

	form.TextareaField("short_description", "Short Description").
		Placeholder("Enter the short description").
		Required(true).
		HelpText("Enter the short description")

	form.NumberField("regular_price", "Regular Price").
		Placeholder("Enter Regular Price").
		Required(true).
		HelpText("Enter Regular Price")

	form.TextField("categories", "Category").
		Placeholder("Enter the category IDs (comma separated)").
		HelpText("Enter the category IDs (comma separated)")

	form.TextField("image_url", "Image URL").
		Placeholder("Enter image URL. Must end with a valid image extension (.jpg, .jpeg, .png, .gif)").
		HelpText("Enter image URL. Must end with a valid image extension (.jpg, .jpeg, .png, .gif)")

	schema := form.Build()

	return schema
}

func (a *CreateProductAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[createProductActionProps](ctx)
	if err != nil {
		return nil, err
	}

	wooClient, err := shared.InitClient(ctx)
	if err != nil {
		return nil, err
	}

	var categories []entity.ProductCategory
	if input.Categories != "" {
		categoryIDs := strings.Split(input.Categories, ",")
		for _, idStr := range categoryIDs {
			id, err := strconv.Atoi(strings.TrimSpace(idStr))
			if err != nil {
				return nil, fmt.Errorf("invalid category ID: %s", idStr)
			}
			categories = append(categories, entity.ProductCategory{ID: id})
		}
	}

	params := woocommerce.CreateProductRequest{
		Name:             input.Name,
		Description:      input.Description,
		Type:             input.Type,
		RegularPrice:     input.RegularPrice,
		ShortDescription: input.ShortDescription,
		Categories:       categories,
	}

	if input.ImageURL != "" {
		params.Images = []entity.ProductImage{
			{
				Src: input.ImageURL,
			},
		}
	}

	product, err := wooClient.Services.Product.Create(params)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (a *CreateProductAction) Auth() *core.AuthMetadata {
	return nil
}

func NewCreateProductAction() sdk.Action {
	return &CreateProductAction{}
}
