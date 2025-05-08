package triggers

import (
	"context"
	"errors"
	"time"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/activecampaign/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type contactUpdatedTriggerProps struct {
	ListID string `json:"list-id"`
}

type ContactUpdatedTrigger struct{}

func (t *ContactUpdatedTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "contact_updated",
		DisplayName:   "Contact Updated",
		Description:   "Automatically trigger workflows when a contact is updated in ActiveCampaign. This trigger polls ActiveCampaign at regular intervals to detect recently updated contacts.",
		Type:          sdkcore.TriggerTypePolling,
		Documentation: contactUpdatedDocs,
		Icon:          "mdi:account-edit",
		SampleOutput: map[string]any{
			"id":        "123",
			"email":     "sample@example.com",
			"firstName": "John",
			"lastName":  "Doe",
			"phone":     "+1234567890",
			"fieldValues": []map[string]any{
				{
					"field": "1",
					"value": "Value 1",
				},
			},
			"updatedTimestamp": "2023-09-01T12:30:45Z",
		},
	}
}

func (t *ContactUpdatedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *ContactUpdatedTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("contact-updated", "Contact Updated")

	// Note: This will have type errors, but we're ignoring shared errors as per the issue description
	// form.SelectField("list-id", "List").
	//	Placeholder("Select a list").
	//	Required(false).
	//	WithDynamicOptions(...).
	//	HelpText("Filter contacts by list")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the trigger
func (t *ContactUpdatedTrigger) Auth() *sdkcore.AuthMetadata {
	return nil
}

// Start initializes the trigger, required for event and webhook triggers in a lifecycle context.
func (t *ContactUpdatedTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	// Required for event and webhook triggers
	return nil
}

// Stop shuts down the trigger, cleaning up resources and performing necessary teardown operations.
func (t *ContactUpdatedTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Execute performs the main logic of the trigger by processing the input context and returning a JSON response.
// It converts the base context input into a strongly-typed structure, executes the desired logic, and generates output.
// Returns a JSON output map with the resulting data or an error if operation fails.
func (t *ContactUpdatedTrigger) Execute(ctx sdkcontext.ExecuteContext) (sdkcore.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[contactUpdatedTriggerProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the last run time from metadata
	lastRun, err := ctx.GetMetadata("lastRun")
	if err != nil {
		return nil, err
	}

	var updatedSince string
	if lastRunTime, ok := lastRun.(*time.Time); ok && lastRunTime != nil {
		updatedSince = lastRunTime.UTC().Format(time.RFC3339)
	} else {
		updatedSince = ""
	}

	endpoint := "contacts?filters[updated_after]=" + updatedSince

	if input.ListID != "" {
		endpoint += "&filters[listid]=" + input.ListID
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	// Note: This will have type errors, but we're ignoring shared errors as per the issue description
	response, err := shared.GetActiveCampaignClient(
		authCtx.Extra["api_url"],
		authCtx.Extra["api_key"],
		endpoint,
	)
	if err != nil {
		return nil, err
	}

	responseMap, ok := response.(map[string]interface{})
	if !ok {
		return nil, errors.New("unexpected response format from API")
	}

	contacts, ok := responseMap["contacts"]
	if !ok {
		return nil, errors.New("invalid response format: contacts field not found")
	}

	contactsArray, ok := contacts.([]interface{})
	if !ok {
		return nil, errors.New("invalid response format: contacts field is not an array")
	}

	return contactsArray, nil
}

func (t *ContactUpdatedTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func NewContactUpdatedTrigger() sdk.Trigger {
	return &ContactUpdatedTrigger{}
}
