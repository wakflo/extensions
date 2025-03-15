package triggers

import (
	"context"
	"errors"
	"time"

	"github.com/wakflo/extensions/internal/integrations/activecampaign/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type contactCreatedTriggerProps struct {
	ListID string `json:"list-id"`
}

type ContactCreatedTrigger struct{}

func (t *ContactCreatedTrigger) Name() string {
	return "Contact Created"
}

func (t *ContactCreatedTrigger) Description() string {
	return "Triggers when a new contact is created in ActiveCampaign."
}

func (t *ContactCreatedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *ContactCreatedTrigger) Icon() *string {
	icon := "mdi:account-plus-outline"
	return &icon
}

func (t *ContactCreatedTrigger) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &contactCreatedDocs,
	}
}

func (t *ContactCreatedTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"list-id": shared.GetActiveCampaignListsInput(),
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

	lastRunTime := ctx.Metadata().LastRun
	var createdSince string

	if lastRunTime != nil {
		createdSince = lastRunTime.UTC().Format(time.RFC3339)
	} else {
		createdSince = ""
	}

	endpoint := "contacts?filters[created_after]=" + createdSince

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

func (t *ContactCreatedTrigger) Auth() *sdk.Auth {
	return nil
}

func (t *ContactCreatedTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *ContactCreatedTrigger) SampleData() sdkcore.JSON {
	return []map[string]any{
		{
			"id":        "123",
			"email":     "john.doe@example.com",
			"firstName": "John",
			"lastName":  "Doe",
			"phone":     "+1234567890",
			"cdate":     "2023-10-15T15:30:00-05:00",
		},
	}
}

func NewContactCreatedTrigger() sdk.Trigger {
	return &ContactCreatedTrigger{}
}
