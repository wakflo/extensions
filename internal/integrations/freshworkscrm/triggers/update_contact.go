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
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type ContactUpdatedTrigger struct{}

func (t *ContactUpdatedTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "contact_updated",
		DisplayName:   "Contact Updated",
		Description:   "Triggers when a contact is updated in Freshworks CRM.",
		Type:          sdkcore.TriggerTypePolling,
		Documentation: updateContactDocs,
		Icon:          "mdi:account-edit-outline",
		SampleOutput: map[string]any{
			"id":            "12345",
			"first_name":    "John",
			"last_name":     "Doe",
			"email":         "john.doe@example.com",
			"mobile_number": "+1234567890",
			"job_title":     "Software Engineer",
			"company":       "Example Inc.",
			"created_at":    "2023-01-01T12:00:00Z",
			"updated_at":    "2023-01-02T14:30:00Z",
		},
	}
}

func (t *ContactUpdatedTrigger) Auth() *sdkcore.AuthMetadata {
	return nil
}

func (t *ContactUpdatedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *ContactUpdatedTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("freshworks-contact-updated", "Contact Updated")

	// No properties needed for this trigger

	schema := form.Build()

	return schema
}

// Start initializes the trigger, required for event and webhook triggers in a lifecycle context.
func (t *ContactUpdatedTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Stop shuts down the trigger, cleaning up resources and performing necessary teardown operations.
func (t *ContactUpdatedTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Execute performs the main action logic by checking for updated contacts.
func (t *ContactUpdatedTrigger) Execute(ctx sdkcontext.ExecuteContext) (sdkcore.JSON, error) {
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	if authCtx.Extra["api-key"] == "" || authCtx.Extra["domain"] == "" {
		return nil, errors.New("missing freshworks auth parameters")
	}

	var lastRunTime *time.Time
	lastRun, err := ctx.GetMetadata("lastRun")
	if err == nil && lastRun != nil {
		lastRunTime = lastRun.(*time.Time)
	}

	queryParams := map[string]string{
		"page":      "1",
		"per_page":  "50",
		"sort":      "updated_at",
		"sort_type": "desc",
	}

	if lastRunTime != nil {
		// Format as ISO8601 for the API
		updatedSince := lastRunTime.UTC().Format(time.RFC3339)
		filterJSON := fmt.Sprintf(`{"updated_at":{"gt":"%s"}}`, updatedSince)
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

	var updatedContacts []interface{}
	for _, contact := range contacts {
		contactMap, ok := contact.(map[string]interface{})
		if !ok {
			continue
		}

		createdAt, createdOk := contactMap["created_at"].(string)
		updatedAt, updatedOk := contactMap["updated_at"].(string)
		if !createdOk || !updatedOk {
			continue
		}

		// Skip if created_at is the same as updated_at (newly created)
		if createdAt != updatedAt {
			updatedContacts = append(updatedContacts, contact)
		}
	}

	// If no updated contacts, return empty array
	if len(updatedContacts) == 0 {
		return []interface{}{}, nil
	}

	return updatedContacts, nil
}

func (t *ContactUpdatedTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func NewContactUpdatedTrigger() sdk.Trigger {
	return &ContactUpdatedTrigger{}
}
