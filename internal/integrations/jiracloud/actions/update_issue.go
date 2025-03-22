package actions

import (
	"encoding/json"
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/jiracloud/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
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

func (a *UpdateIssueAction) Name() string {
	return "Update Issue"
}

func (a *UpdateIssueAction) Description() string {
	return "Update an existing issue in Jira"
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
	icon := "mdi:ticket-edit"
	return &icon
}

func (a *UpdateIssueAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"projectId":   shared.GetProjectsInput(),
		"issueId":     shared.GetIssuesInput(),
		"issueTypeId": shared.GetIssueTypesInput(),
		"summary": autoform.NewShortTextField().
			SetDisplayName("Summary").
			SetRequired(false).
			Build(),
		"description": autoform.NewLongTextField().
			SetDisplayName("Description").
			SetRequired(false).
			Build(),
		"assignee": shared.GetUsersInput(),
		"priority": autoform.NewSelectField().
			SetDisplayName("Priority").
			SetDescription("Priority level of issue").
			SetOptions(shared.PriorityLevels).
			SetRequired(false).
			Build(),
		"parentKey": autoform.NewShortTextField().
			SetDisplayName("Parent Key").
			SetDescription("If this issue is a subtask, insert the parent issue key").
			SetRequired(false).
			Build(),
	}
}

func (a *UpdateIssueAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[updateIssueActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	instanceURL := ctx.Auth.Extra["instance-url"] + "/rest/api/3/issue/" + input.IssueID

	payload := map[string]interface{}{
		"fields": map[string]interface{}{},
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

	if input.IssueTypeID != "" {
		payload["fields"].(map[string]interface{})["issueTypeId"] = map[string]interface{}{
			"id": input.Assignee,
		}
	}

	if input.ProjectID != "" {
		payload["fields"].(map[string]interface{})["project"] = map[string]interface{}{
			"id": input.ProjectID,
		}
	}

	if input.Summary != "" {
		payload["fields"] = map[string]interface{}{
			"summary": input.Summary,
		}
	}

	if input.Priority != "" {
		payload["fields"].(map[string]interface{})["priority"] = map[string]interface{}{
			"id": input.Priority,
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

	resp, err := shared.JiraRequest(ctx.Auth.Extra["email"], ctx.Auth.Extra["api-token"], instanceURL, http.MethodPut, "Issue Updated", data)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (a *UpdateIssueAction) Auth() *sdk.Auth {
	return nil
}

func (a *UpdateIssueAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"Result": "Issue updated successfully",
	}
}

func (a *UpdateIssueAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewUpdateIssueAction() sdk.Action {
	return &UpdateIssueAction{}
}
