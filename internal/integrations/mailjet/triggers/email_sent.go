package triggers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/mailjet/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type emailSentTriggerProps struct {
	Limit int `json:"limit,omitempty"`
}

type EmailSentTrigger struct{}

func (t *EmailSentTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "email_sent",
		DisplayName:   "Email Sent",
		Description:   "Triggers a workflow when an email is successfully sent through Mailjet.",
		Type:          sdkcore.TriggerTypePolling,
		Documentation: emailSentDocs,
		SampleOutput: map[string]any{
			"Count": "2",
			"Data": []map[string]any{
				{
					"ID":                "12345678901234500",
					"MessageUUID":       "1a2b3c4d-5e6f-7g8h-9i0j-1k2l3m4n5o6p",
					"ArrivedAt":         "2023-03-20T14:25:36Z",
					"Campaign":          "Newsletter March 2023",
					"ContactID":         "987654321",
					"Delay":             "0.5",
					"Destination":       "recipient@example.com",
					"FilterTime":        "0.01",
					"From":              "sender@company.com",
					"MessageSize":       "5432",
					"SpamassassinScore": "0.1",
					"Status":            "sent",
					"Subject":           "Your Monthly Newsletter",
					"CustomID":          "newsletter-2023-03",
				},
				{
					"ID":                "12345678901234501",
					"MessageUUID":       "2b3c4d5e-6f7g-8h9i-0j1k-2l3m4n5o6p7",
					"ArrivedAt":         "2023-03-20T14:26:42Z",
					"Campaign":          "Welcome Series",
					"ContactID":         "987654322",
					"Delay":             "0.3",
					"Destination":       "newuser@example.com",
					"FilterTime":        "0.01",
					"From":              "support@company.com",
					"MessageSize":       "3254",
					"SpamassassinScore": " 0.05",
					"Status":            "sent",
					"Subject":           "Welcome to Our Service",
					"CustomID":          "welcome-flow-1",
				},
			},
			"Total": "2",
		},
	}
}

func (t *EmailSentTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("email_sent", "Email Sent")

	form.NumberField("limit", "Limit").
		Placeholder("10").
		Required(false).
		HelpText("Maximum number of emails to process per poll (default: 50, max: 1000)")

	schema := form.Build()

	return schema
}

func (t *EmailSentTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	return nil
}

func (t *EmailSentTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

func (t *EmailSentTrigger) Execute(ctx sdkcontext.ExecuteContext) (sdkcore.JSON, error) {
	lr, err := ctx.GetMetadata("lastRun")
	if err != nil {
		return nil, err
	}

	lastRunTime := lr.(*time.Time)

	input, err := sdk.InputToTypeSafely[emailSentTriggerProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	client, err := shared.GetMailJetClient(authCtx.Extra["api_key"], authCtx.Extra["secret_key"])
	if err != nil {
		return nil, err
	}
	limit := 50
	maxLimit := 1000
	if input.Limit > 0 && input.Limit <= 1000 {
		limit = input.Limit
	} else if input.Limit > maxLimit {
		limit = 1000
	}

	baseURL := "/v3/REST/message"
	queryParams := fmt.Sprintf("?Limit=%d", limit)

	if lastRunTime != nil {
		formattedTime := lastRunTime.UTC().Format(time.RFC3339)
		queryParams += "&FromTS=" + formattedTime
	}

	var result map[string]interface{}
	err = client.Request(http.MethodGet, baseURL+queryParams, nil, &result)
	if err != nil {
		return nil, fmt.Errorf("error fetching sent emails: %v", err)
	}

	data, ok := result["Data"].([]interface{})
	if !ok || len(data) == 0 {
		return []interface{}{}, nil
	}

	return result, nil
}

func (t *EmailSentTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *EmailSentTrigger) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewEmailSentTrigger() sdk.Trigger {
	return &EmailSentTrigger{}
}
