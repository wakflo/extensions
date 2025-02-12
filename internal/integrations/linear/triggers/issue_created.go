package triggers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/wakflo/extensions/internal/integrations/linear/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type issueCreatedTriggerProps struct {
	Name string `json:"name"`
}

type IssueCreatedTrigger struct{}

func (t *IssueCreatedTrigger) Name() string {
	return "Issue Created"
}

func (t *IssueCreatedTrigger) Description() string {
	return "Triggered when a new issue is created in your project management tool, this integration allows you to automate workflows and tasks immediately after an issue is reported, streamlining your team's response time and ensuring prompt attention to new issues."
}

func (t *IssueCreatedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *IssueCreatedTrigger) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &issueCreatedDocs,
	}
}

func (t *IssueCreatedTrigger) Icon() *string {
	return nil
}

func (t *IssueCreatedTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		// "name": autoform.NewShortTextField().
		// 	SetLabel("Name").
		// 	SetRequired(true).
		// 	SetPlaceholder("Your name").
		// 	Build(),
	}
}

// Start initializes the issueCreatedTrigger, required for event and webhook triggers in a lifecycle context.
func (t *IssueCreatedTrigger) Start(ctx sdk.LifecycleContext) error {
	// Required for event and webhook triggers
	return nil
}

// Stop shuts down the issueCreatedTrigger, cleaning up resources and performing necessary teardown operations.
func (t *IssueCreatedTrigger) Stop(ctx sdk.LifecycleContext) error {
	return nil
}

// Execute performs the main action logic of issueCreatedTrigger by processing the input context and returning a JSON response.
// It converts the base context input into a strongly-typed structure, executes the desired logic, and generates output.
// Returns a JSON output map with the resulting data or an error if operation fails. required for Pooling triggers
func (t *IssueCreatedTrigger) Execute(ctx sdk.ExecuteContext) (sdkcore.JSON, error) {
	_, err := sdk.InputToTypeSafely[issueCreatedTriggerProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing linear api key")
	}
	apiKEY := ctx.Auth.Extra["api-key"]

	if !strings.HasPrefix(apiKEY, "lin_api_") {
		return nil, errors.New("invalid Linear API key: must start with 'lin_api_'")
	}

	lastRunTime := ctx.Metadata().LastRun

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
			}`, lastRunTime.UTC().Format(time.RFC3339))
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

func (t *IssueCreatedTrigger) Auth() *sdk.Auth {
	return nil
}

func (t *IssueCreatedTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func NewIssueCreatedTrigger() sdk.Trigger {
	return &IssueCreatedTrigger{}
}
