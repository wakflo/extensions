package triggers

import (
	"context"
	"errors"
	"time"

	"github.com/wakflo/extensions/internal/integrations/activecampaign/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type contactUpdatedTriggerProps struct {
	ListID string `json:"list-id"`
}

type ContactUpdatedTrigger struct{}

func (t *ContactUpdatedTrigger) Name() string {
	return "Contact Updated"
}

func (t *ContactUpdatedTrigger) Description() string {
	return "Automatically trigger workflows when a contact is updated in ActiveCampaign. This trigger polls ActiveCampaign at regular intervals to detect recently updated contacts."
}

func (t *ContactUpdatedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *ContactUpdatedTrigger) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &contactUpdatedDocs,
	}
}

func (t *ContactUpdatedTrigger) Icon() *string {
	icon := "mdi:account-edit"
	return &icon
}

func (t *ContactUpdatedTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"list_id": shared.GetActiveCampaignListsInput(),
	}
}

func (t *ContactUpdatedTrigger) Start(ctx sdk.LifecycleContext) error {
	return nil
}

func (t *ContactUpdatedTrigger) Stop(ctx sdk.LifecycleContext) error {
	return nil
}

func (t *ContactUpdatedTrigger) Execute(ctx sdk.ExecuteContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[contactUpdatedTriggerProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	lastRunTime := ctx.Metadata().LastRun
	var updatedSince string

	if lastRunTime != nil {
		updatedSince = lastRunTime.UTC().Format(time.RFC3339)
	} else {
		updatedSince = ""
	}

	endpoint := "contacts?filters[updated_after]=" + updatedSince

	if input.ListID != "" {
		endpoint += "&filters[listid]=" + input.ListID
	}

	response, err := shared.GetActiveCampaignClient(
		ctx.Auth.Extra["api_url"],
		ctx.Auth.Extra["api_key"],
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

func (t *ContactUpdatedTrigger) Auth() *sdk.Auth {
	return nil
}

func (t *ContactUpdatedTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
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
	}
}

func NewContactUpdatedTrigger() sdk.Trigger {
	return &ContactUpdatedTrigger{}
}
