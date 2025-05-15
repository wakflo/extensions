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

type updateIssueActionProps struct {
	IssueTypeID string `json:"IssueTypeId,omitempty"`
	ProjectID   string `json:"projectId"`
	IssueID     string `json:"issueId"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Assignee    string `json:"assignee"`
	Priority    string `json:"priority"`
	ParentKey   string `json:"parentKey"`
}

type UpdateIssueAction struct{}

// Metadata returns metadata about the action
func (a *UpdateIssueAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "update_issue",
		DisplayName:   "Update Issue",
		Description:   "Update an existing issue in Jira",
		Type:          core.ActionTypeAction,
		Documentation: updateIssueDocs,
		Icon:          "mingcute:edit-line",
		SampleOutput: map[string]any{
			"Result": "Issue updated successfully",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *UpdateIssueAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("update_issue", "Update Issue")

	// Register project selection field
	shared.RegisterProjectsProps(form)

	// Register issue selection field
	shared.RegisterIssuesProps(form)

	// Register issue type selection field
	shared.RegisterIssueTypeProps(form, false)

	form.TextField("summary", "Summary").
		Required(false)

	form.TextField("description", "Description").
		Required(false)

	// Register users selection field
	shared.RegisterUsersProps(form)

	form.TextField("parentKey", "Parent Key").
		Required(false).
		HelpText("If this issue is a subtask, insert the parent issue key")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *UpdateIssueAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *UpdateIssueAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[updateIssueActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	instanceURL := authCtx.Extra["instance-url"] + "/rest/api/3/issue/" + input.IssueID

	fields := make(map[string]interface{})

	if input.Description != "" {
		fields["description"] = map[string]interface{}{
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
		fields["assignee"] = map[string]interface{}{
			"id": input.Assignee,
		}
	}

	if input.IssueTypeID != "" {
		fields["issuetype"] = map[string]interface{}{
			"id": input.IssueTypeID,
		}
	}

	if input.ProjectID != "" {
		fields["project"] = map[string]interface{}{
			"id": input.ProjectID,
		}
	}

	if input.Summary != "" {
		fields["summary"] = input.Summary
	}

	if input.ParentKey != "" {
		fields["parent"] = map[string]interface{}{
			"key": input.ParentKey,
		}
	}

	payload := map[string]interface{}{
		"fields": fields,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := shared.JiraRequest(authCtx.Extra["email"], authCtx.Extra["api-token"], instanceURL, http.MethodPut, "Issue Updated", data)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func NewUpdateIssueAction() sdk.Action {
	return &UpdateIssueAction{}
}
