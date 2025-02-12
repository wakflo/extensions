package triggers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/wakflo/extensions/internal/integrations/mailchimp/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type newSubscriberTriggerProps struct {
	ListID string `json:"list_id"`
}

type NewSubscriberTrigger struct{}

func (t *NewSubscriberTrigger) Name() string {
	return "New Subscriber"
}

func (t *NewSubscriberTrigger) Description() string {
	return "Triggers when a new subscriber is added to your application or service, allowing you to automate tasks and workflows immediately after subscription."
}

func (t *NewSubscriberTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *NewSubscriberTrigger) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &newSubscriberDocs,
	}
}

func (t *NewSubscriberTrigger) Icon() *string {
	return nil
}

func (t *NewSubscriberTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"list_id": autoform.NewShortTextField().
			SetDisplayName("List ID").
			SetDescription("").
			SetRequired(true).
			Build(),
	}
}

// Start initializes the newSubscriberTrigger, required for event and webhook triggers in a lifecycle context.
func (t *NewSubscriberTrigger) Start(ctx sdk.LifecycleContext) error {
	// Required for event and webhook triggers
	return nil
}

// Stop shuts down the newSubscriberTrigger, cleaning up resources and performing necessary teardown operations.
func (t *NewSubscriberTrigger) Stop(ctx sdk.LifecycleContext) error {
	return nil
}

// Execute performs the main action logic of newSubscriberTrigger by processing the input context and returning a JSON response.
// It converts the base context input into a strongly-typed structure, executes the desired logic, and generates output.
// Returns a JSON output map with the resulting data or an error if operation fails. required for Pooling triggers
func (t *NewSubscriberTrigger) Execute(ctx sdk.ExecuteContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[newSubscriberTriggerProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing mailchimp auth token")
	}
	accessToken := ctx.Auth.AccessToken

	dc, err := shared.GetMailChimpServerPrefix(accessToken)
	if err != nil {
		return nil, fmt.Errorf("unable to get mailchimp server prefix: %w", err)
	}

	var fromDate string
	lastRunTime := ctx.Metadata().LastRun

	if lastRunTime != nil {
		fromDate = lastRunTime.UTC().Format(time.RFC3339)
	}

	var newSubscribers interface{}

	newSubscribers, err = shared.ListRecentSubscribers(accessToken, dc, input.ListID, fromDate)
	if err != nil {
		return nil, err
	}
	return newSubscribers, nil
}

func (t *NewSubscriberTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *NewSubscriberTrigger) Auth() *sdk.Auth {
	return nil
}

func (t *NewSubscriberTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func NewNewSubscriberTrigger() sdk.Trigger {
	return &NewSubscriberTrigger{}
}
