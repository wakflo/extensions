package actions

import (
	"encoding/json"
	"net/http"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/jiracloudsoftware/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type listIssuesActionProps struct {
	ProjectID        string `json:"projectId"`
	MaxResults       int    `json:"maxResults,omitempty"`
	OnlyAssignedToMe bool   `json:"onlyAssignedToMe,omitempty"`
}

type ListIssuesAction struct{}

// Metadata returns metadata about the action
func (a *ListIssuesAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "list_issues",
		DisplayName:   "List Issues",
		Description:   "List issues from a Jira project",
		Type:          core.ActionTypeAction,
		Documentation: listIssuesDocs,
		Icon:          "mdi:ticket-outline",
		SampleOutput: map[string]any{
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
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *ListIssuesAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("list_issues", "List Issues")

	// Register project selection field
	shared.RegisterProjectsProps(form)

	form.NumberField("maxResults", "Max Results").
		Required(false).
		HelpText("Maximum number of results to return (default: 50)")

	form.CheckboxField("onlyAssignedToMe", "Only Assigned To Me").
		Required(false).
		DefaultValue(false).
		HelpText("Only show issues assigned to the current user")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *ListIssuesAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *ListIssuesAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[listIssuesActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	email := authCtx.Extra["email"]
	apiToken := authCtx.Extra["api-token"]
	instanceURL := authCtx.Extra["instance-url"]

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

func NewListIssuesAction() sdk.Action {
	return &ListIssuesAction{}
}
