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

type issueUpdatedTriggerProps struct {
	Name string `json:"name"`
}

type IssueUpdatedTrigger struct{}

func (t *IssueUpdatedTrigger) Name() string {
	return "Issue Updated"
}

func (t *IssueUpdatedTrigger) Description() string {
	return "Triggered when an issue is updated in your project management tool, such as Jira or Trello. This integration allows you to automate workflows and tasks based on changes made to issues, including status updates, comments, and attachments."
}

func (t *IssueUpdatedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *IssueUpdatedTrigger) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &issueUpdatedDocs,
	}
}

func (t *IssueUpdatedTrigger) Icon() *string {
	return nil
}

func (t *IssueUpdatedTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		// "name": autoform.NewShortTextField().
		// 	SetLabel("Name").
		// 	SetRequired(true).
		// 	SetPlaceholder("Your name").
		// 	Build(),
	}
}

// Start initializes the issueUpdatedTrigger, required for event and webhook triggers in a lifecycle context.
func (t *IssueUpdatedTrigger) Start(ctx sdk.LifecycleContext) error {
	// Required for event and webhook triggers
	return nil
}

// Stop shuts down the issueUpdatedTrigger, cleaning up resources and performing necessary teardown operations.
func (t *IssueUpdatedTrigger) Stop(ctx sdk.LifecycleContext) error {
	return nil
}

// Execute performs the main action logic of issueUpdatedTrigger by processing the input context and returning a JSON response.
// It converts the base context input into a strongly-typed structure, executes the desired logic, and generates output.
// Returns a JSON output map with the resulting data or an error if operation fails. required for Pooling triggers
func (t *IssueUpdatedTrigger) Execute(ctx sdk.ExecuteContext) (sdkcore.JSON, error) {
	_, err := sdk.InputToTypeSafely[issueUpdatedTriggerProps](ctx.BaseContext)
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
				issues(filter: {updatedAt: {gt: "%s"}}) {
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

func (t *IssueUpdatedTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *IssueUpdatedTrigger) Auth() *sdk.Auth {
	return nil
}

func (t *IssueUpdatedTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func NewIssueUpdatedTrigger() sdk.Trigger {
	return &IssueUpdatedTrigger{}
}
