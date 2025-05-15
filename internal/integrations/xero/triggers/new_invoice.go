package triggers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/xero/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type newInvoiceTriggerProps struct {
	TenantID string `json:"tenant_id"`
}

type NewInvoiceTrigger struct{}

func (t *NewInvoiceTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "new_invoice",
		DisplayName:   "New Invoice",
		Description:   "Triggers when a new invoice is created in your accounting system, this integration allows you to automate workflows and processes immediately after an invoice is generated, streamlining your financial operations and reducing manual errors.",
		Type:          sdkcore.TriggerTypePolling,
		Documentation: newInvoiceDocs,
		SampleOutput: map[string]interface{}{
			"message": "hello world",
		},
	}
}

func (t *NewInvoiceTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("new_invoice", "New Invoice")

	shared.GetTenantProps("tenant_id", "Organization", "select organization", true, form)

	schema := form.Build()

	return schema
}

// Start initializes the newInvoiceTrigger, required for event and webhook triggers in a lifecycle context.
func (t *NewInvoiceTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	// Required for event and webhook triggers
	return nil
}

// Stop shuts down the newInvoiceTrigger, cleaning up resources and performing necessary teardown operations.
func (t *NewInvoiceTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Execute performs the main action logic of newInvoiceTrigger by processing the input context and returning a JSON response.
// It converts the base context input into a strongly-typed structure, executes the desired logic, and generates output.
// Returns a JSON output map with the resulting data or an error if operation fails. required for Pooling triggers
func (t *NewInvoiceTrigger) Execute(ctx sdkcontext.ExecuteContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[newInvoiceTriggerProps](ctx)
	if err != nil {
		return nil, err
	}

	tokenSource := ctx.Auth().Token
	if tokenSource == nil {
		return nil, errors.New("missing authentication token")
	}
	token := tokenSource.AccessToken

	var endpoint string

	lastRunTime, err := ctx.GetMetadata("lastrun")

	if lastRunTime != nil {
		fromDate := lastRunTime.(*time.Time).UTC().Format(time.RFC3339)
		fromDate = strings.ReplaceAll(fromDate, "-", ",")
		endpoint = fmt.Sprintf("/Invoices?where=Date>=DateTime(%s)", fromDate)
	} else {
		endpoint = "/Invoices"
	}

	invoices, err := shared.GetXeroNewClient(token, endpoint, input.TenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch invoices: %v", err)
	}

	return invoices, nil
}

func (t *NewInvoiceTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *NewInvoiceTrigger) Auth() *sdkcore.AuthMetadata {
	return nil
}

func (t *NewInvoiceTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func NewNewInvoiceTrigger() sdk.Trigger {
	return &NewInvoiceTrigger{}
}
