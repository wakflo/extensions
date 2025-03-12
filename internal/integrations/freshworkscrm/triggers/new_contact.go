package triggers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/wakflo/extensions/internal/integrations/freshworkscrm/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type contactCreatedTriggerProps struct {
	Page string `json:"page"`
}

type ContactCreatedTrigger struct{}

func (t *ContactCreatedTrigger) Name() string {
	return "Contact Created"
}

func (t *ContactCreatedTrigger) Description() string {
	return "Triggers when a new contact is created in Freshworks CRM."
}

func (t *ContactCreatedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *ContactCreatedTrigger) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &newContactDocs,
	}
}

func (t *ContactCreatedTrigger) Icon() *string {
	icon := "mdi:account-plus-outline"
	return &icon
}

func (t *ContactCreatedTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"page": autoform.NewShortTextField().
			SetDisplayName("Page Limit").
			SetDescription("Maximum number of contacts to retrieve per page").
			SetRequired(false).
			Build(),
	}
}

func (t *ContactCreatedTrigger) Start(ctx sdk.LifecycleContext) error {
	return nil
}

func (t *ContactCreatedTrigger) Stop(ctx sdk.LifecycleContext) error {
	return nil
}

func (t *ContactCreatedTrigger) Execute(ctx sdk.ExecuteContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[contactCreatedTriggerProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}
	if ctx.Auth.Extra["api-key"] == "" || ctx.Auth.Extra["domain"] == "" {
		return nil, errors.New("missing freshworks auth parameters")
	}

	lastRunTime := ctx.Metadata().LastRun

	queryParams := map[string]string{
		"page":      "1",
		"per_page":  input.Page,
		"sort":      "created_at",
		"sort_type": "desc",
	}

	if lastRunTime != nil {
		createdSince := lastRunTime.UTC().Format(time.RFC3339)
		filterJSON := fmt.Sprintf(`{"created_at":{"gt":"%s"}}`, createdSince)
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

	if len(contacts) == 0 {
		return []interface{}{}, nil
	}

	return contacts, nil
}

func (t *ContactCreatedTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *ContactCreatedTrigger) Auth() *sdk.Auth {
	return nil
}

func (t *ContactCreatedTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"id":            "12345",
		"first_name":    "John",
		"last_name":     "Doe",
		"email":         "john.doe@example.com",
		"mobile_number": "+1234567890",
		"job_title":     "Software Engineer",
		"company":       "Example Inc.",
		"created_at":    "2023-01-01T12:00:00Z",
		"updated_at":    "2023-01-01T12:00:00Z",
	}
}

func (t *ContactCreatedTrigger) Settings() sdkcore.TriggerSettings {
	return sdkcore.TriggerSettings{}
}

func NewContactCreatedTrigger() sdk.Trigger {
	return &ContactCreatedTrigger{}
}
