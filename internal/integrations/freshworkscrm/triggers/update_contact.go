package triggers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/wakflo/extensions/internal/integrations/freshworkscrm/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type ContactUpdatedTrigger struct{}

func (t *ContactUpdatedTrigger) Name() string {
	return "Contact Updated"
}

func (t *ContactUpdatedTrigger) Description() string {
	return "Triggers when a contact is updated in Freshworks CRM."
}

func (t *ContactUpdatedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *ContactUpdatedTrigger) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &updateContactDocs,
	}
}

func (t *ContactUpdatedTrigger) Icon() *string {
	icon := "mdi:account-edit-outline"
	return &icon
}

func (t *ContactUpdatedTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{}
}

// Start initializes the trigger, required for event and webhook triggers in a lifecycle context.
func (t *ContactUpdatedTrigger) Start(ctx sdk.LifecycleContext) error {
	return nil
}

// Stop shuts down the trigger, cleaning up resources and performing necessary teardown operations.
func (t *ContactUpdatedTrigger) Stop(ctx sdk.LifecycleContext) error {
	return nil
}

// Execute performs the main action logic by checking for updated contacts.
func (t *ContactUpdatedTrigger) Execute(ctx sdk.ExecuteContext) (sdkcore.JSON, error) {
	if ctx.Auth.Extra["api-key"] == "" || ctx.Auth.Extra["domain"] == "" {
		return nil, errors.New("missing freshworks auth parameters")
	}

	lastRunTime := ctx.Metadata().LastRun

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

	domain := ctx.Auth.Extra["domain"]
	freshworksDomain := "https://" + domain + ".myfreshworks.com"

	response, err := shared.ListContacts(freshworksDomain, ctx.Auth.Extra["api-key"], queryParams)
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

func (t *ContactUpdatedTrigger) Auth() *sdk.Auth {
	return nil
}

func (t *ContactUpdatedTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"id":            "12345",
		"first_name":    "John",
		"last_name":     "Doe",
		"email":         "john.doe@example.com",
		"mobile_number": "+1234567890",
		"job_title":     "Software Engineer",
		"company":       "Example Inc.",
		"created_at":    "2023-01-01T12:00:00Z",
		"updated_at":    "2023-01-02T14:30:00Z",
	}
}

func (t *ContactUpdatedTrigger) Settings() sdkcore.TriggerSettings {
	return sdkcore.TriggerSettings{}
}

func NewContactUpdatedTrigger() sdk.Trigger {
	return &ContactUpdatedTrigger{}
}
