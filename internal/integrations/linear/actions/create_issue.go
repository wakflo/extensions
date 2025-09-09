package actions

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/linear/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
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

func (a *CreateIssueAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_issue",
		DisplayName:   "Create Issue",
		Description:   "Create Issue: Automatically generates a new issue in your project management tool (e.g., Jira, Trello) based on specific conditions or triggers, ensuring timely and organized tracking of tasks and projects.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: createIssueDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *CreateIssueAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_issue", "Create Issue")

	form.TextField("title", "Issue Name").
		Placeholder("Issue Name").
		Required(true).
		HelpText("The name of the issue.")

	form.TextField("description", "Description").
		Placeholder("Description").
		Required(true).
		HelpText("Issue description.")

	shared.GetTeamsProp(form)

	shared.GetPriorityProp("priority", "Priority", "select issue priority", form)

	shared.GetAssigneesProp("assignee-id", "Assignee", "select an assignee for the issue", form)

	shared.GetTeamLabelsProp("label-id", "Labels", "select an issue label", form)

	schema := form.Build()

	return schema
}

func (a *CreateIssueAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createIssueActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	apiKEY := authCtx.Key

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

func (a *CreateIssueAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewCreateIssueAction() sdk.Action {
	return &CreateIssueAction{}
}
