package actions

import (
	"errors"
	"fmt"
	"strings"

	"github.com/wakflo/extensions/internal/integrations/linear/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type updateIssueActionProps struct {
	TeamID      string `json:"team-id"`
	IssueID     string `json:"issue-id"`
	AssigneeID  string `json:"assignee-id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	StateID     string `json:"state-id"`
	LabelID     string `json:"label-id"`
}

type UpdateIssueAction struct{}

func (a *UpdateIssueAction) Name() string {
	return "Update Issue"
}

func (a *UpdateIssueAction) Description() string {
	return "Updates an existing issue in your project management tool with the latest information from other connected systems, ensuring seamless data synchronization and reducing manual errors."
}

func (a *UpdateIssueAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *UpdateIssueAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &updateIssueDocs,
	}
}

func (a *UpdateIssueAction) Icon() *string {
	return nil
}

func (a *UpdateIssueAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"title": autoform.NewShortTextField().
			SetDisplayName("Issue Name").
			SetDescription("The issue name").
			SetRequired(false).
			Build(),
		"description": autoform.NewLongTextField().
			SetDisplayName("Description").
			SetDescription("Issue description").
			Build(),
		"team-id":     shared.GetTeamsInput(),
		"issue-id":    shared.GetIssuesInput("Select issue", "select an issue to update"),
		"assignee-id": shared.GetAssigneesInput("Assignee", "Select assignee"),
		"priority":    shared.GetPriorityInput("Priority", "Select priority"),
		"label-id":    shared.GetTeamLabelsInput("Label", "Select label"),
		"state-id":    shared.GetIssueStatesInput("Issue State", "select issue state", false),
	}
}

func (a *UpdateIssueAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[updateIssueActionProps](ctx.BaseContext)
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

	mutation := fmt.Sprintf(`mutation IssueUpdate {
      issueUpdate(
        id: "%s",
        input: {`, input.IssueID)

	if input.Title != "" {
		mutation += fmt.Sprintf(`title: "%s",`, input.Title)
	}
	if input.StateID != "" {
		mutation += fmt.Sprintf(`stateId: "%s",`, input.StateID)
	}
	if input.Description != "" {
		mutation += fmt.Sprintf(`description: "%s",`, input.Description)
	}
	if input.LabelID != "" {
		mutation += fmt.Sprintf(`labelIds: "%s",`, input.LabelID)
	}
	if input.Priority != "" {
		mutation += fmt.Sprintf(`priority: %s,`, input.Priority)
	}
	if input.AssigneeID != "" {
		mutation += fmt.Sprintf(`assigneeId: "%s",`, input.AssigneeID)
	}

	mutation = strings.TrimSuffix(mutation, ",")
	mutation += `}) {
        success
        issue {
          id
          title
          priorityLabel
          priority
        }
      }
    }`

	response, err := shared.MakeGraphQLRequest(apiKEY, mutation)
	if err != nil {
		return nil, fmt.Errorf("error making GraphQL request: %w", err)
	}

	issue, ok := response["data"].(map[string]interface{})["issueUpdate"]
	if !ok {
		return nil, errors.New("failed to extract issue from response")
	}

	return issue, nil
}

func (a *UpdateIssueAction) Auth() *sdk.Auth {
	return nil
}

func (a *UpdateIssueAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *UpdateIssueAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewUpdateIssueAction() sdk.Action {
	return &UpdateIssueAction{}
}
