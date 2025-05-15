package actions

import (
	"errors"
	"fmt"
	"time"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/xero/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
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

func (a *CreateInvoiceAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_invoice",
		DisplayName:   "Create Invoice",
		Description:   "Create Invoice: Automatically generates and sends professional-looking invoices to customers based on predefined templates and payment terms, streamlining your accounting process and ensuring timely payments.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: createInvoiceDocs,
		SampleOutput: map[string]any{
			"message": "hello world",
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *CreateInvoiceAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_invoice", "Create Invoice")

	shared.GetInvoiceProp("tenant_id", "Organization", "select organization", true, form)

	form.TextField("contact-id", "Contact ID").
		Placeholder("Contact ID").
		Required(false).
		HelpText("Contact ID")

	form.TextField("email", "Contact Email").
		Placeholder("Contact Email").
		Required(false).
		HelpText("Contact Email")

	form.DateField("due_date", "Due Date").
		Placeholder("Due date of the invoice. Format example: 2019-03-11").
		Required(true).
		DefaultValue(time.Now().Format("2006-01-02")).
		HelpText("Due date of the invoice. Format example: 2019-03-11")

	form.TextField("date", "Date").
		Placeholder("Date the invoice was created. Format example: 2019-03-11").
		Required(false).
		HelpText("Date the invoice was created. Format example: 2019-03-11")

	form.TextField("reference", "Invoice Reference").
		Placeholder("Invoice Reference").
		Required(false).
		HelpText("Reference number of the Invoice")

	form.TextField("contact", "Contact Full Name").
		Placeholder("Contact Name").
		Required(true).
		HelpText("Contact Name")

	lineItemsArray := form.ArrayField("line_items", "Line Items")
	lineItemsArray.Placeholder("List of line items for the invoice.")
	lineItemsArray.Required(true)

	lineItemTemplate := lineItemsArray.ObjectTemplate("lineItem", "")
	lineItemTemplate.TextField("label", "Label").
		Placeholder("Label").
		Required(true)

	lineItemsArray.DefaultValue([]map[string]interface{}{
		{
			"Description": "Default item description",
			"Quantity":    0,
			"UnitAmount":  0,
			"AccountCode": "200",
			"TaxType":     "NONE",
			"LineAmount":  0,
		},
	}).
		HelpText("List of line items for the invoice.")

	form.SelectField("status", "Status").
		Placeholder("Status").
		AddOption("DRAFT", "Draft").
		AddOption("SUBMITTED", "Submitted").
		AddOption("DELETED", "Deleted").
		AddOption("AUTHORISED", "Authorised").
		AddOption("VOIDED", "voided").
		Required(true).
		HelpText("Status")

	schema := form.Build()

	return schema
}

func (a *CreateInvoiceAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createInvoiceActionProps](ctx)
	if err != nil {
		return nil, err
	}

	tokenSource := ctx.Auth().Token
	if tokenSource == nil {
		return nil, errors.New("missing authentication token")
	}
	token := tokenSource.AccessToken

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

	_, err = shared.CreateDraftInvoice(token, input.TenantID, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create invoice: %v", err)
	}

	return map[string]interface{}{
		"Report": "Invoice created Successfully",
	}, nil
}

func (a *CreateInvoiceAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewCreateInvoiceAction() sdk.Action {
	return &CreateInvoiceAction{}
}
