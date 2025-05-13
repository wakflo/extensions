package triggers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/freshworkscrm/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type contactCreatedTriggerProps struct {
	Page string `json:"page"`
}

type ContactCreatedTrigger struct{}

// Metadata returns metadata about the trigger
func (t *ContactCreatedTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "contact_created",
		DisplayName:   "Contact Created",
		Description:   "Triggers when a new contact is created in Freshworks CRM.",
		Type:          core.TriggerTypePolling,
		Documentation: newContactDocs,
		SampleOutput: map[string]any{
			"id":            "12345",
			"first_name":    "John",
			"last_name":     "Doe",
			"email":         "john.doe@example.com",
			"mobile_number": "+1234567890",
			"job_title":     "Software Engineer",
			"company":       "Example Inc.",
			"created_at":    "2023-01-01T12:00:00Z",
			"updated_at":    "2023-01-01T12:00:00Z",
		},
	}
}

// Props returns the schema for the trigger's input configuration
func (t *ContactCreatedTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("contact_created", "Contact Created")

	// Add page field
	form.TextField("page", "Page Limit").
		Placeholder("Enter page limit").
		Required(false).
		HelpText("Maximum number of contacts to retrieve per page")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the trigger
func (t *ContactCreatedTrigger) Auth() *core.AuthMetadata {
	return nil
}

// Start initializes the trigger
func (t *ContactCreatedTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Stop shuts down the trigger
func (t *ContactCreatedTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Execute performs the main trigger logic
func (t *ContactCreatedTrigger) Execute(ctx sdkcontext.ExecuteContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[contactCreatedTriggerProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	if authCtx.Extra["api-key"] == "" || authCtx.Extra["domain"] == "" {
		return nil, errors.New("missing freshworks auth parameters")
	}

	// Get the last run time
	lastRunTime, err := ctx.GetMetadata("lastRun")
	if err != nil {
		return nil, err
	}

	queryParams := map[string]string{
		"page":      "1",
		"per_page":  input.Page,
		"sort":      "created_at",
		"sort_type": "desc",
	}

	if lastRunTime != nil {
		lastRunTimeValue := lastRunTime.(*time.Time)
		createdSince := lastRunTimeValue.UTC().Format(time.RFC3339)
		filterJSON := fmt.Sprintf(`{"created_at":{"gt":"%s"}}`, createdSince)
		queryParams["filter"] = filterJSON
	}

	domain := authCtx.Extra["domain"]
	freshworksDomain := "https://" + domain + ".myfreshworks.com"

	response, err := shared.ListContacts(freshworksDomain, authCtx.Extra["api-key"], queryParams)
	if err != nil {
		return nil, fmt.Errorf("error fetching contacts: %v", err)
	}

	responseMap, ok := response.(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid response format")
	}

	contacts, ok := responseMap["contacts"].([]interface{})
	if !ok {
		return nil, errors.New("invalid response format: contacts field is not an array")
	}

	if len(contacts) == 0 {
		return []interface{}{}, nil
	}

	return contacts, nil
}

// Criteria returns the criteria for the trigger
func (t *ContactCreatedTrigger) Criteria(ctx context.Context) core.TriggerCriteria {
	return core.TriggerCriteria{}
}

func NewContactCreatedTrigger() sdk.Trigger {
	return &ContactCreatedTrigger{}
}
