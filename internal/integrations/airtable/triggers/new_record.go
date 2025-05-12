package triggers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/airtable/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type newRecordTriggerProps struct {
	Bases string `json:"bases,omitempty"`
}

type NewRecordTrigger struct{}

func (t *NewRecordTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "new_record",
		DisplayName:   "New Record",
		Description:   "Triggers when a new record is created in the specified table or database, allowing you to automate actions and workflows immediately after a new record is added.",
		Type:          sdkcore.TriggerTypePolling,
		Documentation: newRecordDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
	}
}

func (t *NewRecordTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *NewRecordTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("new_record", "New Record")

	// Note: This will have type errors, but we're ignoring shared errors as per the issue description
	// form.SelectField("bases", "Bases").
	//	Placeholder("Select a base").
	//	Required(true).
	//	WithDynamicOptions(...).
	//	HelpText("The base to monitor for new records")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the trigger
func (t *NewRecordTrigger) Auth() *sdkcore.AuthMetadata {
	return nil
}

// Start initializes the newRecordTrigger, required for event and webhook triggers in a lifecycle context.
func (t *NewRecordTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	// Required for event and webhook triggers
	return nil
}

// Stop shuts down the newRecordTrigger, cleaning up resources and performing necessary teardown operations.
func (t *NewRecordTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Execute performs the main action logic of newRecordTrigger by processing the input context and returning a JSON response.
// It converts the base context input into a strongly-typed structure, executes the desired logic, and generates output.
// Returns a JSON output map with the resulting data or an error if operation fails.
func (t *NewRecordTrigger) Execute(ctx sdkcontext.ExecuteContext) (sdkcore.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[newRecordTriggerProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the last run time from metadata
	lastRun, err := ctx.GetMetadata("lastRun")
	if err != nil {
		return nil, err
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	if authCtx.Extra["api-key"] == "" {
		return nil, errors.New("missing airtable api key")
	}
	apiKey := authCtx.Extra["api-key"]

	var createdTime string
	if lastRunTime, ok := lastRun.(*time.Time); ok && lastRunTime != nil {
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

func NewNewRecordTrigger() sdk.Trigger {
	return &NewRecordTrigger{}
}
