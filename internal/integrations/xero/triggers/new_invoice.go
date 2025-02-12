package triggers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/wakflo/extensions/internal/integrations/xero/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type newInvoiceTriggerProps struct {
	TenantID string `json:"tenant_id"`
}

type NewInvoiceTrigger struct{}

func (t *NewInvoiceTrigger) Name() string {
	return "New Invoice"
}

func (t *NewInvoiceTrigger) Description() string {
	return "Triggered when a new invoice is created in your accounting system, this integration allows you to automate workflows and processes immediately after an invoice is generated, streamlining your financial operations and reducing manual errors."
}

func (t *NewInvoiceTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *NewInvoiceTrigger) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &newInvoiceDocs,
	}
}

func (t *NewInvoiceTrigger) Icon() *string {
	return nil
}

func (t *NewInvoiceTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"tenant_id": shared.GetTenantInput("Organization", "select organization", true),
	}
}

// Start initializes the newInvoiceTrigger, required for event and webhook triggers in a lifecycle context.
func (t *NewInvoiceTrigger) Start(ctx sdk.LifecycleContext) error {
	// Required for event and webhook triggers
	return nil
}

// Stop shuts down the newInvoiceTrigger, cleaning up resources and performing necessary teardown operations.
func (t *NewInvoiceTrigger) Stop(ctx sdk.LifecycleContext) error {
	return nil
}

// Execute performs the main action logic of newInvoiceTrigger by processing the input context and returning a JSON response.
// It converts the base context input into a strongly-typed structure, executes the desired logic, and generates output.
// Returns a JSON output map with the resulting data or an error if operation fails. required for Pooling triggers
func (t *NewInvoiceTrigger) Execute(ctx sdk.ExecuteContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[newInvoiceTriggerProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	var endpoint string

	lastRunTime := ctx.Metadata().LastRun

	if lastRunTime != nil {
		fromDate := lastRunTime.UTC().Format(time.RFC3339)
		fromDate = strings.ReplaceAll(fromDate, "-", ",")
		endpoint = fmt.Sprintf("/Invoices?where=Date>=DateTime(%s)", fromDate)
	} else {
		endpoint = "/Invoices"
	}

	invoices, err := shared.GetXeroNewClient(ctx.Auth.AccessToken, endpoint, input.TenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch invoices: %v", err)
	}

	return invoices, nil
}

func (t *NewInvoiceTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *NewInvoiceTrigger) Auth() *sdk.Auth {
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
