package actions

import (
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/prisync/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
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

func (a *AddProductAction) Name() string {
	return "Add Product"
}

func (a *AddProductAction) Description() string {
	return "Adds a new product to your inventory or database, allowing you to track and manage products efficiently. This integration action can be used to populate product information from various sources, such as e-commerce platforms, marketplaces, or product information management systems."
}

func (a *AddProductAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *AddProductAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &addProductDocs,
	}
}

func (a *AddProductAction) Icon() *string {
	return nil
}

func (a *AddProductAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"name": autoform.NewShortTextField().
			SetDisplayName("Product Name").
			SetDescription("name of product").
			SetRequired(true).Build(),
		"brand": autoform.NewShortTextField().
			SetDisplayName("Brand").
			SetDescription("Brand name").
			SetRequired(true).Build(),
		"category": autoform.NewShortTextField().
			SetDisplayName("Category").
			SetDescription("Category name").
			SetRequired(true).Build(),
		"product_code": autoform.NewShortTextField().
			SetDisplayName("Product Code").
			SetDescription("Product code").
			SetRequired(false).Build(),
		"barcode": autoform.NewShortTextField().
			SetDisplayName("Bar Code").
			SetDescription("Bar code").
			SetRequired(false).Build(),
		"cost": autoform.NewShortTextField().
			SetDisplayName("Product Cost").
			SetDescription("Product cost").
			SetRequired(false).Build(),
		"tags": autoform.NewLongTextField().
			SetDisplayName("Product Tags").
			SetDescription("Product tags").
			SetRequired(false).Build(),
	}
}

func (a *AddProductAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[addProductActionProps](ctx.BaseContext)
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
	resp, err := shared.PrisyncRequest(ctx.Auth.Extra["api-key"], ctx.Auth.Extra["api-token"], endpoint, http.MethodPost, formData)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (a *AddProductAction) Auth() *sdk.Auth {
	return nil
}

func (a *AddProductAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *AddProductAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewAddProductAction() sdk.Action {
	return &AddProductAction{}
}
