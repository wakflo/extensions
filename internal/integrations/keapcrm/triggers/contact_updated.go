package triggers

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/juicycleff/smartform/v1"

	"github.com/wakflo/extensions/internal/integrations/keapcrm/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type ContactUpdatedTrigger struct{}

func (t *ContactUpdatedTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "contact_updated",
		DisplayName:   "Contact Updated",
		Description:   "Triggers when a contact is updated in Keap",
		Type:          sdkcore.TriggerTypePolling,
		Documentation: contactUpdatedDocs,
		SampleOutput: map[string]any{
			"contacts": []map[string]any{
				{
					"id":          "123",
					"given_name":  "John",
					"family_name": "Doe",
					"email_addresses": []map[string]string{
						{
							"email": "john.doe@example.com",
							"type":  "PRIMARY",
						},
					},
					"date_created": "2023-06-15T14:30:25.039Z",
					"last_updated": "2023-06-16T14:30:25.039Z",
				},
			},
		},
	}
}

func (t *ContactUpdatedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *ContactUpdatedTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("contact_created", "Contact Created")

	schema := form.Build()

	return schema
}

func (t *ContactUpdatedTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	return nil
}

func (t *ContactUpdatedTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

func (t *ContactUpdatedTrigger) Execute(ctx sdkcontext.ExecuteContext) (sdkcore.JSON, error) {
	// Determine the time range to check for updated contacts
	lastRunTime, err := ctx.GetMetadata("lastRun")
	if err != nil {
		return nil, err
	}

	tokenSource := ctx.Auth().Token
	if tokenSource == nil {
		return nil, errors.New("missing authentication token")
	}
	token := tokenSource.AccessToken

	var updatedSince string
	if lastRunTime != nil {
		updatedSince = lastRunTime.(*time.Time).UTC().Format(time.RFC3339Nano)
	} else {
		updatedSince = ""
	}

	queryParams := url.Values{}
	queryParams.Add("limit", "100")
	queryParams.Add("order", "last_updated")
	queryParams.Add("order_direction", "ascending")
	queryParams.Add("since", updatedSince)

	endpoint := "/contacts?" + queryParams.Encode()

	response, err := shared.MakeKeapRequest(token, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}

	return response, nil
}

func (t *ContactUpdatedTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *ContactUpdatedTrigger) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewContactUpdatedTrigger() sdk.Trigger {
	return &ContactUpdatedTrigger{}
}
