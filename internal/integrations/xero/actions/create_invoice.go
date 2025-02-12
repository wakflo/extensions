package actions

import (
	"fmt"
	"time"

	"github.com/wakflo/extensions/internal/integrations/xero/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type createInvoiceActionProps struct {
	TenantID  string                   `json:"tenant_id"`
	ContactID string                   `json:"contact-id"`
	Contact   string                   `json:"contact"`
	LineItems []map[string]interface{} `json:"line_items"`
	DueDate   string                   `json:"due_date"`
	Status    string                   `json:"status"`
	Date      string                   `json:"date"`
	Email     string                   `json:"email"`
	Reference string                   `json:"reference"`
}

type CreateInvoiceAction struct{}

func (a *CreateInvoiceAction) Name() string {
	return "Create Invoice"
}

func (a *CreateInvoiceAction) Description() string {
	return "Create Invoice: Automatically generates and sends professional-looking invoices to customers based on predefined templates and payment terms, streamlining your accounting process and ensuring timely payments."
}

func (a *CreateInvoiceAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *CreateInvoiceAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createInvoiceDocs,
	}
}

func (a *CreateInvoiceAction) Icon() *string {
	return nil
}

func (a *CreateInvoiceAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"tenant_id": shared.GetTenantInput("Organization", "select organization", true),
		"contact-id": autoform.NewShortTextField().
			SetDisplayName("Contact ID").
			SetDescription("Contact ID").
			SetRequired(false).
			Build(),
		"email": autoform.NewShortTextField().
			SetDisplayName("Contact Email").
			SetDescription("Contact Email").
			SetRequired(false).
			Build(),
		"due_date": autoform.NewShortTextField().
			SetDisplayName("Due Date").
			SetDescription("Due date of the invoice. Format example: 2019-03-11").
			SetDefaultValue(time.Now().Format("2006-01-02")).
			SetRequired(true).
			Build(),
		"date": autoform.NewShortTextField().
			SetDisplayName("Date").
			SetDescription("Date the invoice was created. Format example: 2019-03-11").
			SetRequired(false).
			Build(),
		"reference": autoform.NewShortTextField().
			SetDisplayName("Invoice Reference").
			SetDescription("Reference number of the Invoice").
			SetRequired(false).
			Build(),
		"contact": autoform.NewShortTextField().
			SetDisplayName("Contact Full Name").
			SetDescription("Contact Name").
			SetRequired(true).
			Build(),
		"line_items": autoform.NewArrayField().
			SetDisplayName("Line Items").
			SetDescription("List of line items for the invoice.").
			SetRequired(true).
			SetItems(
				autoform.NewShortTextField().
					SetDisplayName("Label").
					SetDescription("Label").
					SetRequired(true).
					Build(),
			).
			SetDefaultValue([]map[string]interface{}{
				{
					"Description": "Default item description",
					"Quantity":    0,
					"UnitAmount":  0,
					"AccountCode": "200",
					"TaxType":     "NONE",
					"LineAmount":  0,
				},
			}).
			Build(),
		"status": autoform.NewSelectField().
			SetDisplayName("Status").
			SetOptions(shared.XeroInvoiceStatus).
			SetRequired(true).
			Build(),
	}
}

func (a *CreateInvoiceAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createInvoiceActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	body := map[string]interface{}{
		"Invoices": []map[string]interface{}{
			{
				"Type": "ACCREC",
				"Contact": map[string]interface{}{
					"Name": input.Contact,
				},
				"LineItems": func() []map[string]interface{} {
					if len(input.LineItems) > 0 {
						return input.LineItems
					}
					return []map[string]interface{}{}
				}(),
				"Date":      input.Date,
				"DueDate":   input.DueDate,
				"Reference": input.Reference,
				"Status":    input.Status,
			},
		},
	}

	_, err = shared.CreateDraftInvoice(ctx.Auth.AccessToken, input.TenantID, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create invoice: %v", err)
	}

	return map[string]interface{}{
		"Report": "Invoice created Successfully",
	}, nil
}

func (a *CreateInvoiceAction) Auth() *sdk.Auth {
	return nil
}

func (a *CreateInvoiceAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *CreateInvoiceAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreateInvoiceAction() sdk.Action {
	return &CreateInvoiceAction{}
}
