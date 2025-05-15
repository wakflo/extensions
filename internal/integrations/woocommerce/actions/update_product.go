package actions

import (
	"github.com/hiscaler/woocommerce-go"
	"github.com/hiscaler/woocommerce-go/entity"
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/woocommerce/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type updateProductActionProps struct {
	ProductID        int     `json:"productId"`
	Name             string  `json:"name"`
	RegularPrice     float64 `json:"regular_price"`
	SalePrice        float64 `json:"sale_price"`
	Description      string  `json:"description"`
	ShortDescription string  `json:"short_description"`
	Weight           string  `json:"weight"`
	Length           float64 `json:"length"`
	Width            float64 `json:"width"`
	Height           float64 `json:"height"`
	ImageURL         string  `json:"image_url"`
}

type UpdateProductAction struct{}

func (a *UpdateProductAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "update_product",
		DisplayName:   "Update Product",
		Description:   "Updates product information in your WooCommerce store.",
		Type:          core.ActionTypeAction,
		Documentation: updateProductDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *UpdateProductAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("update_product", "Update Product")

	form.NumberField("productId", "Product ID").
		Placeholder("Enter the product ID").
		Required(true).
		HelpText("Enter the product ID")

	form.TextField("name", "Product Name").
		Placeholder("Enter product Name").
		HelpText("Enter product Name")

	form.TextareaField("description", "Description").
		Placeholder("Enter product description").
		HelpText("Enter product description")

	form.TextareaField("short_description", "Short Description").
		Placeholder("Enter the short description").
		HelpText("Enter the short description")

	form.NumberField("length", "Length").
		Placeholder("Enter Product Length").
		HelpText("Enter Product Length")

	form.NumberField("regular_price", "Regular Price").
		Placeholder("Enter Regular Price").
		Required(true).
		HelpText("Enter Regular Price")

	form.NumberField("sale_price", "Discounted Price").
		Placeholder("Enter Discounted Price").
		Required(true).
		HelpText("Enter Discounted Price")

	form.NumberField("height", "Height").
		Placeholder("Enter Product Height").
		HelpText("Enter Product Height")

	form.NumberField("width", "Width").
		Placeholder("Enter Product Width").
		HelpText("Enter Product Width")

	form.TextField("weight", "Weight").
		Placeholder("Enter weight").
		HelpText("weight")

	form.TextField("image_url", "Image URL").
		Placeholder("Enter image URL. Must end with a valid image extension (.jpg, .jpeg, .png, .gif)").
		HelpText("Enter image URL. Must end with a valid image extension (.jpg, .jpeg, .png, .gif)")

	schema := form.Build()

	return schema
}

func (a *UpdateProductAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[updateProductActionProps](ctx)
	if err != nil {
		return nil, err
	}

	wooClient, err := shared.InitClient(ctx)
	if err != nil {
		return nil, err
	}

	params := woocommerce.UpdateProductRequest{}
	if input.Name != "" {
		params.Name = input.Name
	}

	if input.Description != "" {
		params.Description = input.Description
	}

	if input.RegularPrice != 0 {
		params.RegularPrice = input.RegularPrice
	}

	if input.SalePrice != 0 {
		params.SalePrice = input.SalePrice
	}

	if input.ShortDescription != "" {
		params.ShortDescription = input.ShortDescription
	}

	if input.Weight != "" {
		params.Weight = input.Weight
	}

	if input.Length != 0 {
		params.Dimensions.Length = input.Length
	}

	if input.Width != 0 {
		params.Dimensions.Width = input.Width
	}

	if input.Height != 0 {
		params.Dimensions.Length = input.Height
	}

	if input.ImageURL != "" {
		params.Images = []entity.ProductImage{
			{
				Src: input.ImageURL,
			},
		}
	}

	updatedProduct, err := wooClient.Services.Product.Update(input.ProductID, params)
	if err != nil {
		return nil, err
	}

	return updatedProduct, nil
}

func (a *UpdateProductAction) Auth() *core.AuthMetadata {
	return nil
}

func NewUpdateProductAction() sdk.Action {
	return &UpdateProductAction{}
}
