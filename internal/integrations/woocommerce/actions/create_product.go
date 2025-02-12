package actions

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hiscaler/woocommerce-go"
	"github.com/hiscaler/woocommerce-go/entity"
	"github.com/wakflo/extensions/internal/integrations/woocommerce/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type createProductActionProps struct {
	Name             string  `json:"name"`
	Type             string  `json:"type"`
	RegularPrice     float64 `json:"regular_price"`
	Description      string  `json:"description"`
	ShortDescription string  `json:"short_description"`
	Categories       string  `json:"categories"`
}

type CreateProductAction struct{}

func (a *CreateProductAction) Name() string {
	return "Create Product"
}

func (a *CreateProductAction) Description() string {
	return "Create Product: Automatically generates and creates new products in your system, including product details such as name, description, price, and inventory levels."
}

func (a *CreateProductAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *CreateProductAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createProductDocs,
	}
}

func (a *CreateProductAction) Icon() *string {
	return nil
}

func (a *CreateProductAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"name": autoform.NewShortTextField().
			SetDisplayName("Product Name").
			SetDescription("Enter product Name").
			SetRequired(true).
			Build(),
		"type": autoform.NewSelectField().
			SetDisplayName("Type").
			SetDescription("Select the type").
			SetOptions(shared.ProductType).
			SetRequired(true).
			Build(),
		"description": autoform.NewLongTextField().
			SetDisplayName(" Description").
			SetDescription("Enter product description").
			SetRequired(true).
			Build(),
		"short_description": autoform.NewLongTextField().
			SetDisplayName("Short Description").
			SetDescription("Enter the short description").
			SetRequired(true).
			Build(),
		"regular_price": autoform.NewNumberField().
			SetDisplayName("Regular Price").
			SetDescription("Enter Regular Price").
			SetRequired(true).
			Build(),
		"categories": autoform.NewShortTextField().
			SetDisplayName("Category").
			SetDescription("Enter the category IDs (comma separated)").
			Build(),
	}
}

func (a *CreateProductAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createProductActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	wooClient, err := shared.InitClient(ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	// Parse categories
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

	// Create a query parameters struct
	params := woocommerce.CreateProductRequest{
		Name:             input.Name,
		Description:      input.Description,
		Type:             input.Type,
		RegularPrice:     input.RegularPrice,
		ShortDescription: input.ShortDescription,
		Categories:       categories,
	}

	product, err := wooClient.Services.Product.Create(params)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (a *CreateProductAction) Auth() *sdk.Auth {
	return nil
}

func (a *CreateProductAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *CreateProductAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreateProductAction() sdk.Action {
	return &CreateProductAction{}
}
