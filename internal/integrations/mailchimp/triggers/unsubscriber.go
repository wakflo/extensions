package triggers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/mailchimp/shared"

	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type unsubscriberTriggerProps struct {
	ListID string `json:"list_id"`
}

type UnsubscriberTrigger struct{}

func (t *UnsubscriberTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "unsubscriber",
		DisplayName:   "Unsubscriber",
		Description:   "Triggers when a subscriber unsubscribes from your application or service, allowing you to automate tasks and workflows immediately after unsubscription.",
		Type:          sdkcore.TriggerTypePolling,
		Documentation: unsubscriberDocs,
		SampleOutput: map[string]any{
			"email": "john.doe@example.com",
		},
	}
}

func (t *UnsubscriberTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *UnsubscriberTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("unsubscriber", "Unsubscriber")

	form.TextField("list_id", "List ID").
		Placeholder("Enter a value for List ID.").
		Required(true).
		HelpText("The ID of the list to check for unsubscribers.")

	schema := form.Build()

	return schema
}

// Start initializes the unsubscriberTrigger, required for event and webhook triggers in a lifecycle context.
func (t *UnsubscriberTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	// Required for event and webhook triggers
	return nil
}

// Stop shuts down the unsubscriberTrigger, cleaning up resources and performing necessary teardown operations.
func (t *UnsubscriberTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Execute performs the main action logic of unsubscriberTrigger by processing the input context and returning a JSON response.
// It converts the base context input into a strongly-typed structure, executes the desired logic, and generates output.
// Returns a JSON output map with the resulting data or an error if operation fails. required for Pooling triggers
func (t *UnsubscriberTrigger) Execute(ctx sdkcontext.ExecuteContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[unsubscriberTriggerProps](ctx)
	if err != nil {
		return nil, err
	}

	tokenSource := ctx.Auth().Token
	if tokenSource == nil {
		return nil, errors.New("missing authentication token")
	}
	token := tokenSource.AccessToken

	dc, err := shared.GetMailChimpServerPrefix(token)
	if err != nil {
		return nil, fmt.Errorf("unable to get mailchimp server prefix: %w", err)
	}

	var fromDate string
	lastRunTime, err := ctx.GetMetadata("lastRun")
	if err != nil {
		return nil, err
	}

	if lastRunTime != nil {
		fromDate = lastRunTime.(*time.Time).UTC().Format(time.RFC3339)
	}

	var unsubscribes interface{}

	unsubscribes, err = shared.ListRecentUnSubscribers(token, dc, input.ListID, fromDate)
	if err != nil {
		return nil, err
	}

	return unsubscribes, nil
}

func (t *UnsubscriberTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *UnsubscriberTrigger) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewUnsubscriberTrigger() sdk.Trigger {
	return &UnsubscriberTrigger{}
}
