package actions

import (
	"errors"
	"fmt"
	"strings"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/linear/shared"

	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
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

func (a *UpdateIssueAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "update_issue",
		DisplayName:   "Update Issue",
		Description:   "Update Issue: Automatically updates an existing issue in your project management tool with the latest information from other connected systems, ensuring seamless data synchronization and reducing manual errors.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: updateIssueDocs,
		SampleOutput: map[string]any{
			"success": true,
			"issue": map[string]any{
				"id":            "issue-123",
				"title":         "Updated Issue Title",
				"priorityLabel": "High",
				"priority":      1,
			},
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *UpdateIssueAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("update_issue", "Update Issue")

	form.TextField("title", "Issue Name").
		Placeholder("Issue Name").
		Required(false).
		HelpText("The name of the issue.")

	form.TextField("description", "Description").
		Placeholder("Description").
		Required(false).
		HelpText("Issue description.")

	shared.GetTeamsProp(form)

	shared.GetIssuesProp("issue-id", "Issue", "select an issue to update", form)

	shared.GetAssigneesProp("assignee-id", "Assignee", "select an assignee for the issue", form)

	shared.GetPriorityProp("priority", "Priority", "select issue priority", form)

	shared.GetTeamLabelsProp("label-id", "Labels", "select an issue label", form)

	shared.GetIssueStatesProp("state-id", "Issue State", "select issue state", false, form)

	schema := form.Build()

	return schema
}

func (a *UpdateIssueAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	// Parse input
	input, err := sdk.InputToTypeSafely[updateIssueActionProps](ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input: %w", err)
	}

	// Validate required fields
	if input.IssueID == "" {
		return nil, errors.New("issue ID is required")
	}

	// Check if at least one field to update is provided
	if input.Title == "" && input.StateID == "" && input.Description == "" &&
		input.LabelID == "" && input.Priority == "" && input.AssigneeID == "" {
		return nil, errors.New("at least one field must be provided to update")
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, fmt.Errorf("failed to get auth context: %w", err)
	}

	// Validate API key exists
	apiKEY, ok := authCtx.Extra["api-key"]
	if !ok || apiKEY == "" {
		return nil, errors.New("Linear API key not found in auth context")
	}

	if !strings.HasPrefix(apiKEY, "lin_api_") {
		return nil, errors.New("invalid Linear API key: must start with 'lin_api_'")
	}

	// Build mutation
	mutation := fmt.Sprintf(`mutation IssueUpdate {
      issueUpdate(
        id: "%s",
        input: {`, input.IssueID)

	updateFields := []string{}

	if input.Title != "" {
		updateFields = append(updateFields, fmt.Sprintf(`title: "%s"`, escapeQuotes(input.Title)))
	}
	if input.StateID != "" {
		updateFields = append(updateFields, fmt.Sprintf(`stateId: "%s"`, input.StateID))
	}
	if input.Description != "" {
		updateFields = append(updateFields, fmt.Sprintf(`description: "%s"`, escapeQuotes(input.Description)))
	}
	if input.LabelID != "" {
		updateFields = append(updateFields, fmt.Sprintf(`labelIds: ["%s"]`, input.LabelID))
	}
	if input.Priority != "" {
		updateFields = append(updateFields, fmt.Sprintf(`priority: %s`, input.Priority))
	}
	if input.AssigneeID != "" {
		updateFields = append(updateFields, fmt.Sprintf(`assigneeId: "%s"`, input.AssigneeID))
	}

	mutation += strings.Join(updateFields, ", ")
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

	// Make GraphQL request
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

func (a *UpdateIssueAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

// Helper function to escape quotes in strings for GraphQL
func escapeQuotes(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "\n", `\n`)
	s = strings.ReplaceAll(s, "\r", `\r`)
	s = strings.ReplaceAll(s, "\t", `\t`)
	return s
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func NewUpdateIssueAction() sdk.Action {
	return &UpdateIssueAction{}
}
