package actions

import (
	"errors"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/harvest/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type listInvoicesActionProps struct {
	Name string `json:"name"`
}

type ListInvoicesAction struct{}

// Metadata returns metadata about the action
func (a *ListInvoicesAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "list_invoices",
		DisplayName:   "List Invoices",
		Description:   "Retrieve and list all invoices associated with a specific account or organization, allowing you to easily track and manage your financial transactions.",
		Type:          core.ActionTypeAction,
		Documentation: listInvoicesDocs,
		SampleOutput: []map[string]any{
			{
				"id":         13150453,
				"client_id":  5735776,
				"number":     "1001",
				"amount":     288.23,
				"due_amount": 288.23,
				"subject":    "Web Design",
				"state":      "open",
				"issue_date": "2023-05-01",
				"due_date":   "2023-05-31",
				"sent_at":    "2023-05-02T14:30:00Z",
			},
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *ListInvoicesAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("list_invoices", "List Invoices")

	// The original doesn't have any properties

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *ListInvoicesAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *ListInvoicesAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	url := "/v2/invoices"

	invoices, err := shared.GetHarvestClient(authCtx.Token.AccessToken, url)
	if err != nil {
		return nil, err
	}
	invoiceArray, ok := invoices["invoices"].(interface{})
	if !ok {
		return nil, errors.New("failed to extract issues from response")
	}
	return invoiceArray, nil
}

func NewListInvoicesAction() sdk.Action {
	return &ListInvoicesAction{}
}
