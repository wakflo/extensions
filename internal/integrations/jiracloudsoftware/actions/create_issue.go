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

type createIssueActionProps struct {
	IssueTypeID string `json:"IssueTypeId,omitempty"`
	ProjectID   string `json:"projectId"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Assignee    string `json:"assignee"`
	ParentKey   string `json:"parentKey"`
}

type CreateIssueAction struct{}

// Metadata returns metadata about the action
func (a *CreateIssueAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_issue",
		DisplayName:   "Create Issue",
		Description:   "Create a new issue in Jira with specified details",
		Type:          core.ActionTypeAction,
		Documentation: createIssueDocs,
		Icon:          "lucide:ticket-plus",
		SampleOutput: map[string]any{
			"id":   "12345",
			"key":  "PRJ-123",
			"self": "https://yourcompany.atlassian.net/rest/api/3/issue/12345",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *CreateIssueAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_issue", "Create Issue")

	shared.RegisterProjectsProps(form)

	shared.RegisterIssuesProps(form)

	shared.RegisterIssueTypeProps(form, false)

	form.TextField("summary", "Summary").
		Required(true)

	form.TextareaField("description", "Description").
		Required(false)

	shared.RegisterUsersProps(form)

	form.TextField("parentKey", "Parent Key").
		Required(false).
		HelpText("If this issue is a subtask, insert the parent issue key")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *CreateIssueAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *CreateIssueAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[createIssueActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	instanceURL := authCtx.Extra["instance-url"] + "/rest/api/3/issue"

	payload := map[string]interface{}{
		"fields": map[string]interface{}{
			"summary": input.Summary,
			"issuetype": map[string]interface{}{
				"id": input.IssueTypeID,
			},
			"project": map[string]interface{}{
				"id": input.ProjectID,
			},
		},
	}

	if input.Description != "" {
		payload["fields"].(map[string]interface{})["description"] = map[string]interface{}{
			"type":    "doc",
			"version": 1,
			"content": []map[string]interface{}{
				{
					"type": "paragraph",
					"content": []map[string]interface{}{
						{
							"type": "text",
							"text": input.Description,
						},
					},
				},
			},
		}
	}

	if input.Assignee != "" {
		payload["fields"].(map[string]interface{})["assignee"] = map[string]interface{}{
			"id": input.Assignee,
		}
	}

	if input.ParentKey != "" {
		payload["fields"].(map[string]interface{})["parent"] = map[string]interface{}{
			"key": input.ParentKey,
		}
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	reqURL := instanceURL
	resp, err := shared.JiraRequest(authCtx.Extra["email"], authCtx.Extra["api-token"], reqURL, http.MethodPost, "", data)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func NewCreateIssueAction() sdk.Action {
	return &CreateIssueAction{}
}
