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

type ContactCreatedTrigger struct{}

func (t *ContactCreatedTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "contact_created",
		DisplayName:   "Contact Created",
		Description:   "Triggers when a new contact is created in Keap",
		Type:          sdkcore.TriggerTypePolling,
		Documentation: contactCreatedDocs,
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
					"date_created": "2023-06-15T14:30:25Z",
					"last_updated": "2023-06-15T14:30:25Z",
				},
			},
		},
	}
}

func (t *ContactCreatedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *ContactCreatedTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("contact_created", "Contact Created")

	schema := form.Build()

	return schema
}

func (t *ContactCreatedTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	return nil
}

func (t *ContactCreatedTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

func (t *ContactCreatedTrigger) Execute(ctx sdkcontext.ExecuteContext) (sdkcore.JSON, error) {
	lastRunTime, err := ctx.GetMetadata("lastRun")
	if err != nil {
		return nil, err
	}

	tokenSource := ctx.Auth().Token
	if tokenSource == nil {
		return nil, errors.New("missing authentication token")
	}
	token := tokenSource.AccessToken

	var createdSince string
	if lastRunTime != nil {
		createdSince = lastRunTime.(*time.Time).UTC().Format(time.RFC3339Nano)
	} else {
		createdSince = ""
	}

	queryParams := url.Values{}
	queryParams.Add("limit", "100")
	queryParams.Add("order", "date_created")
	queryParams.Add("order_direction", "ascending")
	queryParams.Add("since", createdSince)

	queryParams.Add("optional_properties", "email")
	queryParams.Add("optional_properties", "last_updated")

	endpoint := "/contacts?" + queryParams.Encode()

	response, err := shared.MakeKeapRequest(token, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}

	return response, nil
}

func (t *ContactCreatedTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *ContactCreatedTrigger) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewContactCreatedTrigger() sdk.Trigger {
	return &ContactCreatedTrigger{}
}
