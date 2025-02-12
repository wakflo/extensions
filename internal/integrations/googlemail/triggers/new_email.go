package triggers

import (
	"context"
	"time"

	"github.com/wakflo/extensions/internal/integrations/googlemail/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type newEmailTriggerProps struct {
	Subject    string `json:"subject"`
	From       string `json:"from"`
	MaxResults int    `json:"maxResults"`
}

type NewEmailTrigger struct{}

func (t *NewEmailTrigger) Name() string {
	return "New Email"
}

func (t *NewEmailTrigger) Description() string {
	return "Triggered when a new email is received in your inbox or mailbox, allowing you to automate workflows based on incoming emails."
}

func (t *NewEmailTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *NewEmailTrigger) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &newEmailDocs,
	}
}

func (t *NewEmailTrigger) Icon() *string {
	return nil
}

func (t *NewEmailTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"subject": autoform.NewShortTextField().
			SetDisplayName("Subject").
			SetDescription("Subject of the email to trigger on").
			SetRequired(true).
			Build(),
		"from": autoform.NewShortTextField().
			SetDisplayName("From").
			SetDescription("Sender of the email to trigger on").
			SetRequired(false).
			Build(),
		"maxResults": autoform.NewNumberField().
			SetDisplayName("Max Results").
			SetDescription("Maximum number of emails to return").
			SetDefaultValue(50).
			SetRequired(false).
			Build(),
	}
}

// Start initializes the newEmailTrigger, required for event and webhook triggers in a lifecycle context.
func (t *NewEmailTrigger) Start(ctx sdk.LifecycleContext) error {
	// Required for event and webhook triggers
	return nil
}

// Stop shuts down the newEmailTrigger, cleaning up resources and performing necessary teardown operations.
func (t *NewEmailTrigger) Stop(ctx sdk.LifecycleContext) error {
	return nil
}

// Execute performs the main action logic of newEmailTrigger by processing the input context and returning a JSON response.
// It converts the base context input into a strongly-typed structure, executes the desired logic, and generates output.
// Returns a JSON output map with the resulting data or an error if operation fails. required for Pooling triggers
func (t *NewEmailTrigger) Execute(ctx sdk.ExecuteContext) (sdkcore.JSON, error) {
	_, err := sdk.InputToTypeSafely[newEmailTriggerProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	gmailService, err := gmail.NewService(context.Background(), option.WithTokenSource(*ctx.Auth.TokenSource))
	if err != nil {
		return nil, err
	}

	var lastRunTime time.Time
	if ctx.Metadata().LastRun != nil {
		lastRunTime = *ctx.Metadata().LastRun
	}

	query := "after:" + lastRunTime.Format(time.DateOnly)

	messages, err := gmailService.Users.Messages.List("me").Q(query).Do()
	if err != nil {
		return nil, err
	}

	newEmails := make([]map[string]interface{}, 0, len(messages.Messages))
	for _, msg := range messages.Messages {
		email, err := gmailService.Users.Messages.Get("me", msg.Id).Do()
		if err != nil {
			continue
		}

		emailData := map[string]interface{}{
			"id":      email.Id,
			"subject": shared.GetHeader(email.Payload.Headers, "Subject"),
			"from":    shared.GetHeader(email.Payload.Headers, "From"),
			"date":    shared.GetHeader(email.Payload.Headers, "Date"),
		}
		newEmails = append(newEmails, emailData)
	}

	return newEmails, nil
}

func (t *NewEmailTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *NewEmailTrigger) Auth() *sdk.Auth {
	return nil
}

func (t *NewEmailTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func NewNewEmailTrigger() sdk.Trigger {
	return &NewEmailTrigger{}
}
