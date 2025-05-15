package triggers

import (
	"context"
	"time"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/googlemail/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type newEmailTriggerProps struct {
	Subject    string `json:"subject"`
	From       string `json:"from"`
	MaxResults int    `json:"maxResults"`
}

type NewEmailTrigger struct{}

func (t *NewEmailTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "new_email",
		DisplayName:   "New Email",
		Description:   "Triggered when a new email is received in your inbox or mailbox, allowing you to automate workflows based on incoming emails.",
		Type:          sdkcore.TriggerTypePolling,
		Documentation: newEmailDocs,
		SampleOutput: map[string]any{
			"emails": []map[string]any{
				{
					"id":      "12345abcde",
					"subject": "Important Message",
					"from":    "sender@example.com",
					"date":    "Mon, 01 Jan 2025 10:00:00 -0700",
				},
			},
		},
	}
}

func (t *NewEmailTrigger) Auth() *sdkcore.AuthMetadata {
	return nil
}

func (t *NewEmailTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *NewEmailTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("google-mail-new-email", "New Email")

	form.TextField("subject", "subject").
		Placeholder("Subject").
		HelpText("Subject of the email to trigger on").
		Required(true)

	form.TextField("from", "from").
		Placeholder("From").
		HelpText("Sender of the email to trigger on").
		Required(false)

	form.NumberField("maxResults", "maxResults").
		Placeholder("Max Results").
		HelpText("Maximum number of emails to return").
		DefaultValue(50).
		Required(false)

	schema := form.Build()

	return schema
}

// Start initializes the newEmailTrigger, required for event and webhook triggers in a lifecycle context.
func (t *NewEmailTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	// Required for event and webhook triggers
	return nil
}

// Stop shuts down the newEmailTrigger, cleaning up resources and performing necessary teardown operations.
func (t *NewEmailTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Execute performs the main action logic of newEmailTrigger by processing the input context and returning a JSON response.
// It converts the base context input into a strongly-typed structure, executes the desired logic, and generates output.
// Returns a JSON output map with the resulting data or an error if operation fails. required for Pooling triggers
func (t *NewEmailTrigger) Execute(ctx sdkcontext.ExecuteContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[newEmailTriggerProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	gmailService, err := gmail.NewService(context.Background(), option.WithTokenSource(*authCtx.TokenSource))
	if err != nil {
		return nil, err
	}

	var lastRunTime time.Time
	lastRun, err := ctx.GetMetadata("lastRun")
	if err == nil && lastRun != nil {
		lastRunTime = *lastRun.(*time.Time)
	}

	query := "after:" + lastRunTime.Format(time.DateOnly)

	// Add subject filter if provided
	if input.Subject != "" {
		query += " subject:" + input.Subject
	}

	// Add sender filter if provided
	if input.From != "" {
		query += " from:" + input.From
	}

	maxResults := int64(50)
	if input.MaxResults > 0 {
		maxResults = int64(input.MaxResults)
	}

	messages, err := gmailService.Users.Messages.List("me").
		MaxResults(maxResults).
		Q(query).Do()
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

	// Update lastRun metadata to current time for next execution
	now := time.Now()
	ctx.SetMetadata("lastRun", &now)

	return newEmails, nil
}

func (t *NewEmailTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func NewNewEmailTrigger() sdk.Trigger {
	return &NewEmailTrigger{}
}
