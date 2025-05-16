package triggers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/juicycleff/smartform/v1"

	"github.com/wakflo/extensions/internal/integrations/linear/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type issueCreatedTriggerProps struct {
	Name string `json:"name"`
}

type IssueCreatedTrigger struct{}

func (t *IssueCreatedTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "issue_created",
		DisplayName:   "Issue Created",
		Description:   "Triggered when a new issue is created in your project management tool, this integration allows you to automate workflows and tasks immediately after an issue is reported, streamlining your team's response time and ensuring prompt attention to new issues.",
		Type:          sdkcore.TriggerTypePolling,
		Documentation: issueCreatedDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
	}
}

func (t *IssueCreatedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *IssueCreatedTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("issue_created", "Issue Created")

	schema := form.Build()

	return schema
}

// Start initializes the issueCreatedTrigger, required for event and webhook triggers in a lifecycle context.
func (t *IssueCreatedTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	// Required for event and webhook triggers
	return nil
}

// Stop shuts down the issueCreatedTrigger, cleaning up resources and performing necessary teardown operations.
func (t *IssueCreatedTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Execute performs the main action logic of issueCreatedTrigger by processing the input context and returning a JSON response.
// It converts the base context input into a strongly-typed structure, executes the desired logic, and generates output.
// Returns a JSON output map with the resulting data or an error if operation fails. required for Pooling triggers
func (t *IssueCreatedTrigger) Execute(ctx sdkcontext.ExecuteContext) (sdkcore.JSON, error) {
	_, err := sdk.InputToTypeSafely[issueCreatedTriggerProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}
	apiKEY := authCtx.Extra["api-key"]

	if !strings.HasPrefix(apiKEY, "lin_api_") {
		return nil, errors.New("invalid Linear API key: must start with 'lin_api_'")
	}

	lastRunTime, err := ctx.GetMetadata("lastRun")
	if err != nil {
		return nil, err
	}

	var query string
	if lastRunTime == nil {
		query = `{
				issues {
					nodes {
						id
						title
						updatedAt
					}
				}
			}`
	} else {
		query = fmt.Sprintf(`{
				issues(filter: {createdAt: {gt: "%s"}}) {
					nodes {
						id
						title
						updatedAt
					}
				}
			}`, lastRunTime.(*time.Time).UTC().Format(time.RFC3339))
	}

	response, err := shared.MakeGraphQLRequest(apiKEY, query)
	if err != nil {
		return nil, fmt.Errorf("error making GraphQL request: %w", err)
	}

	return map[string]interface{}{
		"Result": response,
	}, nil
}

func (t *IssueCreatedTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *IssueCreatedTrigger) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewIssueCreatedTrigger() sdk.Trigger {
	return &IssueCreatedTrigger{}
}
