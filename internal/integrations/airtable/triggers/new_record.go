package triggers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/wakflo/extensions/internal/integrations/airtable/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type newRecordTriggerProps struct {
	Bases string `json:"bases,omitempty"`
}

type NewRecordTrigger struct{}

func (t *NewRecordTrigger) Name() string {
	return "New Record"
}

func (t *NewRecordTrigger) Description() string {
	return "Triggers when a new record is created in the specified table or database, allowing you to automate actions and workflows immediately after a new record is added."
}

func (t *NewRecordTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *NewRecordTrigger) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &newRecordDocs,
	}
}

func (t *NewRecordTrigger) Icon() *string {
	return nil
}

func (t *NewRecordTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"bases": shared.GetBasesInput(),
	}
}

// Start initializes the newRecordTrigger, required for event and webhook triggers in a lifecycle context.
func (t *NewRecordTrigger) Start(ctx sdk.LifecycleContext) error {
	// Required for event and webhook triggers
	return nil
}

// Stop shuts down the newRecordTrigger, cleaning up resources and performing necessary teardown operations.
func (t *NewRecordTrigger) Stop(ctx sdk.LifecycleContext) error {
	return nil
}

// Execute performs the main action logic of newRecordTrigger by processing the input context and returning a JSON response.
// It converts the base context input into a strongly-typed structure, executes the desired logic, and generates output.
// Returns a JSON output map with the resulting data or an error if operation fails. required for Pooling triggers
func (t *NewRecordTrigger) Execute(ctx sdk.ExecuteContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[newRecordTriggerProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	lastRunTime := ctx.Metadata().LastRun

	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing airtable api key")
	}
	apiKey := ctx.Auth.Extra["api-key"]

	var createdTime string
	if lastRunTime != nil {
		createdTime = lastRunTime.UTC().Format(time.RFC3339)
	}
	reqURL := fmt.Sprintf("%s/v0/meta/bases/%s/tables?updated_since=%s", shared.BaseAPI, input.Bases, createdTime)

	response, err := shared.AirtableRequest(apiKey, reqURL, http.MethodGet)
	if err != nil {
		return nil, errors.New("error fetching data")
	}

	return response, nil
}

func (t *NewRecordTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *NewRecordTrigger) Auth() *sdk.Auth {
	return nil
}

func (t *NewRecordTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func NewNewRecordTrigger() sdk.Trigger {
	return &NewRecordTrigger{}
}
