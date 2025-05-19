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

type issueUpdatedTriggerProps struct {
	Name string `json:"name"`
}

type IssueUpdatedTrigger struct{}

func (t *IssueUpdatedTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "issue_updated",
		DisplayName:   "Issue Updated",
		Description:   "Triggered when an issue is updated in your project management tool, such as Jira or Trello. This integration allows you to automate workflows and tasks based on changes made to issues, including status updates, comments, and attachments.",
		Type:          sdkcore.TriggerTypePolling,
		Documentation: issueUpdatedDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
	}
}

func (t *IssueUpdatedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *IssueUpdatedTrigger) Props() *smartform.FormSchema {
	// return map[string]*sdkcore.AutoFormSchema{
	// "name": autoform.NewShortTextField().
	// 	SetLabel("Name").
	// 	SetRequired(true).
	// 	SetPlaceholder("Your name").
	// 	Build(),
	// }

	form := smartform.NewForm("issue_updated", "Issue Updated")

	schema := form.Build()

	return schema
}

// Start initializes the issueUpdatedTrigger, required for event and webhook triggers in a lifecycle context.
func (t *IssueUpdatedTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	// Required for event and webhook triggers
	return nil
}

// Stop shuts down the issueUpdatedTrigger, cleaning up resources and performing necessary teardown operations.
func (t *IssueUpdatedTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Execute performs the main action logic of issueUpdatedTrigger by processing the input context and returning a JSON response.
// It converts the base context input into a strongly-typed structure, executes the desired logic, and generates output.
// Returns a JSON output map with the resulting data or an error if operation fails. required for Pooling triggers
func (t *IssueUpdatedTrigger) Execute(ctx sdkcontext.ExecuteContext) (sdkcore.JSON, error) {
	_, err := sdk.InputToTypeSafely[issueUpdatedTriggerProps](ctx)
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
				issues(filter: {updatedAt: {gt: "%s"}}) {
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

func (t *IssueUpdatedTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *IssueUpdatedTrigger) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewIssueUpdatedTrigger() sdk.Trigger {
	return &IssueUpdatedTrigger{}
}
