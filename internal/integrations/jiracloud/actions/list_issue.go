package actions

import (
	"encoding/json"
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/jiracloud/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type listIssuesActionProps struct {
	ProjectID        string `json:"projectId"`
	MaxResults       int    `json:"maxResults,omitempty"`
	OnlyAssignedToMe bool   `json:"onlyAssignedToMe,omitempty"`
}

type ListIssuesAction struct{}

func (a *ListIssuesAction) Name() string {
	return "List Issues"
}

func (a *ListIssuesAction) Description() string {
	return "List issues from a Jira project"
}

func (a *ListIssuesAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *ListIssuesAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &listIssuesDocs,
	}
}

func (a *ListIssuesAction) Icon() *string {
	icon := "mdi:ticket-outline"
	return &icon
}

func (a *ListIssuesAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"projectId": shared.GetProjectsInput(),
		"maxResults": autoform.NewNumberField().
			SetDisplayName("Max Results").
			SetDescription("Maximum number of results to return (default: 50)").
			SetRequired(false).Build(),
		"onlyAssignedToMe": autoform.NewBooleanField().
			SetDisplayName("Only Assigned To Me").
			SetDescription("Only show issues assigned to the current user").
			SetDefaultValue(false).
			SetRequired(false).Build(),
	}
}

func (a *ListIssuesAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[listIssuesActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	email := ctx.Auth.Extra["email"]
	apiToken := ctx.Auth.Extra["api-token"]
	instanceURL := ctx.Auth.Extra["instance-url"]

	jql := "project=" + input.ProjectID

	if input.OnlyAssignedToMe {
		jql += " AND assignee = currentUser()"
	}

	maxResults := 50
	if input.MaxResults > 0 {
		maxResults = input.MaxResults
	}

	requestBody := map[string]interface{}{
		"jql":        jql,
		"maxResults": maxResults,
		"fields":     []string{"summary", "status", "priority", "assignee", "updated", "created"},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	response, err := shared.JiraRequest(
		email,
		apiToken,
		instanceURL+"/rest/api/3/search",
		http.MethodPost,
		"",
		jsonBody,
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (a *ListIssuesAction) Auth() *sdk.Auth {
	return nil
}

func (a *ListIssuesAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"issues": []map[string]any{
			{
				"id":   "12345",
				"key":  "PRJ-123",
				"self": "https://yourcompany.atlassian.net/rest/api/3/issue/12345",
				"fields": map[string]any{
					"summary": "Sample issue 1",
					"status": map[string]any{
						"name": "To Do",
					},
					"priority": map[string]any{
						"name": "Medium",
					},
					"assignee": map[string]any{
						"displayName": "John Doe",
					},
					"created": "2023-01-15T10:30:45.123+0000",
					"updated": "2023-01-16T14:22:33.456+0000",
				},
			},
			{
				"id":   "12346",
				"key":  "PRJ-124",
				"self": "https://yourcompany.atlassian.net/rest/api/3/issue/12346",
				"fields": map[string]any{
					"summary": "Sample issue 2",
					"status": map[string]any{
						"name": "In Progress",
					},
					"priority": map[string]any{
						"name": "High",
					},
					"assignee": map[string]any{
						"displayName": "Jane Smith",
					},
					"created": "2023-01-17T09:15:30.789+0000",
					"updated": "2023-01-18T11:45:12.345+0000",
				},
			},
		},
		"total":      "2",
		"maxResults": "50",
		"startAt":    "0",
	}
}

func (a *ListIssuesAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewListIssuesAction() sdk.Action {
	return &ListIssuesAction{}
}
