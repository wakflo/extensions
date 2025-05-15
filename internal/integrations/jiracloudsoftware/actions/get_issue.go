package actions

import (
	"net/http"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/jiracloudsoftware/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type getIssueActionProps struct {
	ProjectID string `json:"projectId"`
	IssueID   string `json:"issueId"`
}

type GetIssueAction struct{}

// Metadata returns metadata about the action
func (a *GetIssueAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_issue",
		DisplayName:   "Get Issue",
		Description:   "Get details of a specific Jira issue",
		Type:          core.ActionTypeAction,
		Documentation: getIssueDocs,
		Icon:          "mdi:ticket",
		SampleOutput: map[string]any{
			"id":   "12345",
			"key":  "PRJ-123",
			"self": "https://yourcompany.atlassian.net/rest/api/3/issue/12345",
			"fields": map[string]any{
				"summary": "Sample issue",
				"description": map[string]any{
					"type": "doc",
					"content": []map[string]any{
						{
							"type": "paragraph",
							"content": []map[string]any{
								{
									"type": "text",
									"text": "This is a sample issue description",
								},
							},
						},
					},
				},
				"status": map[string]any{
					"name": "To Do",
				},
				"priority": map[string]any{
					"name": "Medium",
				},
				"assignee": map[string]any{
					"displayName":  "John Doe",
					"emailAddress": "john.doe@example.com",
				},
			},
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *GetIssueAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_issue", "Get Issue")

	// Register project selection field
	shared.RegisterProjectsProps(form)

	// Register issue selection field
	shared.RegisterIssuesProps(form)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *GetIssueAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *GetIssueAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[getIssueActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	instanceURL := authCtx.Extra["instance-url"] + "/rest/api/3/issue/" + input.IssueID

	resp, err := shared.JiraRequest(authCtx.Extra["email"], authCtx.Extra["api-token"], instanceURL, http.MethodGet, "", nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func NewGetIssueAction() sdk.Action {
	return &GetIssueAction{}
}
