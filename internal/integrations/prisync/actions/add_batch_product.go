package actions

import (
	"github.com/wakflo/extensions/internal/integrations/prisync/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
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

func (a *AddBatchProductAction) Name() string {
	return "Add Batch Product"
}

func (a *AddBatchProductAction) Description() string {
	return "Add Batch Product: Automatically adds a new product to your batch, allowing you to manage and track multiple products within a single batch. This integration action enables seamless product management, streamlining your workflow and reducing manual errors."
}

func (a *AddBatchProductAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *AddBatchProductAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &addBatchProductDocs,
	}
}

func (a *AddBatchProductAction) Icon() *string {
	return nil
}

func (a *AddBatchProductAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"name": autoform.NewShortTextField().
			SetLabel("Name").
			SetRequired(true).
			SetPlaceholder("Your name").
			Build(),
	}
}

func (a *AddBatchProductAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[addBatchProductActionProps](ctx.BaseContext)
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
	resp, err := shared.PrisyncBatchRequest(ctx.Auth.Extra["api-key"], ctx.Auth.Extra["api-token"], endpoint, formData, false)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (a *AddBatchProductAction) Auth() *sdk.Auth {
	return nil
}

func (a *AddBatchProductAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *AddBatchProductAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewAddBatchProductAction() sdk.Action {
	return &AddBatchProductAction{}
}
