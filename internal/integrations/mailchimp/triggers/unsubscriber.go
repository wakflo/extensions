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

type unsubscriberTriggerProps struct {
	ListID string `json:"list_id"`
}

type UnsubscriberTrigger struct{}

func (t *UnsubscriberTrigger) Name() string {
	return "Unsubscriber"
}

func (t *UnsubscriberTrigger) Description() string {
	return "The Unsubscriber integration trigger is designed to automatically remove subscribers from your workflow when they unsubscribe from a specific email list or service. This trigger can be used in conjunction with other automation workflows to ensure that unsubscribed contacts are no longer targeted for marketing campaigns, surveys, or other automated tasks. By integrating the Unsubscriber trigger with your existing workflows, you can maintain data accuracy and compliance with anti-spam laws by promptly removing unsubscribed contacts from your workflow."
}

func (t *UnsubscriberTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *UnsubscriberTrigger) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &unsubscriberDocs,
	}
}

func (t *UnsubscriberTrigger) Icon() *string {
	return nil
}

func (t *UnsubscriberTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"list_id": autoform.NewShortTextField().
			SetDisplayName(" List ID").
			SetDescription("").
			SetRequired(true).
			Build(),
	}
}

// Start initializes the unsubscriberTrigger, required for event and webhook triggers in a lifecycle context.
func (t *UnsubscriberTrigger) Start(ctx sdk.LifecycleContext) error {
	// Required for event and webhook triggers
	return nil
}

// Stop shuts down the unsubscriberTrigger, cleaning up resources and performing necessary teardown operations.
func (t *UnsubscriberTrigger) Stop(ctx sdk.LifecycleContext) error {
	return nil
}

// Execute performs the main action logic of unsubscriberTrigger by processing the input context and returning a JSON response.
// It converts the base context input into a strongly-typed structure, executes the desired logic, and generates output.
// Returns a JSON output map with the resulting data or an error if operation fails. required for Pooling triggers
func (t *UnsubscriberTrigger) Execute(ctx sdk.ExecuteContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[unsubscriberTriggerProps](ctx.BaseContext)
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

	var unsubscribes interface{}

	unsubscribes, err = shared.ListRecentUnSubscribers(accessToken, dc, input.ListID, fromDate)
	if err != nil {
		return nil, err
	}

	return unsubscribes, nil
}

func (t *UnsubscriberTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *UnsubscriberTrigger) Auth() *sdk.Auth {
	return nil
}

func (t *UnsubscriberTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func NewUnsubscriberTrigger() sdk.Trigger {
	return &UnsubscriberTrigger{}
}
