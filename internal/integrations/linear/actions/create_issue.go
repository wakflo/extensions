package actions

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/wakflo/extensions/internal/integrations/linear/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type createIssueActionProps struct {
	TeamID      string `json:"team-id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	AssigneeID  string `json:"assignee-id"`
	LabelID     string `json:"label-id"`
	StateID     string `json:"state-id"`
}

type CreateIssueAction struct{}

func (a *CreateIssueAction) Name() string {
	return "Create Issue"
}

func (a *CreateIssueAction) Description() string {
	return "Create Issue: Automatically generates a new issue in your project management tool (e.g., Jira, Trello) based on specific conditions or triggers, ensuring timely and organized tracking of tasks and projects."
}

func (a *CreateIssueAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *CreateIssueAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createIssueDocs,
	}
}

func (a *CreateIssueAction) Icon() *string {
	return nil
}

func (a *CreateIssueAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"title": autoform.NewShortTextField().
			SetDisplayName("Issue Name").
			SetDescription("The issue name").
			SetRequired(true).
			Build(),
		"description": autoform.NewLongTextField().
			SetDisplayName("Description").
			SetDescription("Issue description").
			Build(),
		"team-id":     shared.GetTeamsInput(),
		"priority":    shared.GetPriorityInput("Priority", "select issue priority"),
		"assignee-id": shared.GetAssigneesInput("Assignee", "select an assignee for the issue"),
		"label-id":    shared.GetTeamLabelsInput("Labels", "select an issue label"),
	}
}

func (a *CreateIssueAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createIssueActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	apiKEY := ctx.Auth.Extra["api-key"]

	if !strings.HasPrefix(apiKEY, "lin_api_") {
		return nil, errors.New("invalid Linear API key: must start with 'lin_api_'")
	}

	// Create a map to store fields conditionally
	fields := make(map[string]string)
	fields["title"] = fmt.Sprintf(`"%s"`, input.Title)
	fields["teamId"] = fmt.Sprintf(`"%s"`, input.TeamID)

	if input.Description != "" {
		fields["description"] = fmt.Sprintf(`"%s"`, input.Description)
	}
	if input.AssigneeID != "" {
		fields["assigneeId"] = fmt.Sprintf(`"%s"`, input.AssigneeID)
	}
	if input.LabelID != "" {
		fields["labelIds"] = fmt.Sprintf(`"%s"`, input.LabelID)
	}

	if input.StateID != "" {
		fields["stateId"] = fmt.Sprintf(`"%s"`, input.StateID)
	}

	if input.Priority != "" {
		priority, err := strconv.Atoi(input.Priority)
		if err != nil {
			return nil, err
		}
		fields["priority"] = strconv.Itoa(priority)
	}

	fieldStrings := make([]string, 0, len(fields))
	for key, value := range fields {
		fieldStrings = append(fieldStrings, fmt.Sprintf("%s: %s", key, value))
	}

	query := fmt.Sprintf(`
		mutation IssueCreate {
			issueCreate(input: {
				%s
			}) {
				success
				issue {
					id
					title
					description
					priorityLabel
				}
			}
		}`, strings.Join(fieldStrings, "\n"))

	response, err := shared.MakeGraphQLRequest(apiKEY, query)
	if err != nil {
		return nil, fmt.Errorf("error making GraphQL request: %w", err)
	}

	issue, ok := response["data"].(map[string]interface{})["issueCreate"].(map[string]interface{})["issue"].(map[string]interface{})
	if !ok {
		return nil, errors.New("failed to extract issue from response")
	}

	return issue, nil
}

func (a *CreateIssueAction) Auth() *sdk.Auth {
	return nil
}

func (a *CreateIssueAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *CreateIssueAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreateIssueAction() sdk.Action {
	return &CreateIssueAction{}
}
