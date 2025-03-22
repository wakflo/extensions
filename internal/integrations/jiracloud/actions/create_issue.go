package actions

import (
	"encoding/json"
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/jiracloud/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type createIssueActionProps struct {
	IssueTypeID string `json:"IssueTypeId,omitempty"`
	ProjectID   string `json:"projectId"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Assignee    string `json:"assignee"`
	Priority    string `json:"priority"`
	ParentKey   string `json:"parentKey"`
}

type CreateIssueAction struct{}

func (a *CreateIssueAction) Name() string {
	return "Create Issue"
}

func (a *CreateIssueAction) Description() string {
	return "Create a new issue in Jira with specified details"
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
	icon := "mdi:ticket-plus"
	return &icon
}

func (a *CreateIssueAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"projectId":   shared.GetProjectsInput(),
		"issueTypeId": shared.GetIssueTypesInput(),
		"summary": autoform.NewShortTextField().
			SetDisplayName("Summary").
			SetRequired(true).
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

func (a *CreateIssueAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createIssueActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	instanceURL := ctx.Auth.Extra["instance-url"] + "/rest/api/3/issue"

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

	reqURL := instanceURL
	resp, err := shared.JiraRequest(ctx.Auth.Extra["email"], ctx.Auth.Extra["api-token"], reqURL, http.MethodPost, "", data)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (a *CreateIssueAction) Auth() *sdk.Auth {
	return nil
}

func (a *CreateIssueAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"id":   "12345",
		"key":  "PRJ-123",
		"self": "https://yourcompany.atlassian.net/rest/api/3/issue/12345",
	}
}

func (a *CreateIssueAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreateIssueAction() sdk.Action {
	return &CreateIssueAction{}
}
