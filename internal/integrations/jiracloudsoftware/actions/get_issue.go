package actions

import (
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/jiracloudsoftware/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type getIssueActionProps struct {
	ProjectID string `json:"projectId"`
	IssueID   string `json:"issueId"`
}

type GetIssueAction struct{}

func (a *GetIssueAction) Name() string {
	return "Get Issue"
}

func (a *GetIssueAction) Description() string {
	return "Get details of a specific Jira issue"
}

func (a *GetIssueAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *GetIssueAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getIssueDocs,
	}
}

func (a *GetIssueAction) Icon() *string {
	icon := "mdi:ticket"
	return &icon
}

func (a *GetIssueAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"projectId": shared.GetProjectsInput(),
		"issueId":   shared.GetIssuesInput(),
	}
}

func (a *GetIssueAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getIssueActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	instanceURL := ctx.Auth.Extra["instance-url"] + "/rest/api/3/issue/" + input.IssueID

	resp, err := shared.JiraRequest(ctx.Auth.Extra["email"], ctx.Auth.Extra["api-token"], instanceURL, http.MethodGet, "", nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (a *GetIssueAction) Auth() *sdk.Auth {
	return nil
}

func (a *GetIssueAction) SampleData() sdkcore.JSON {
	return map[string]any{
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
	}
}

func (a *GetIssueAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewGetIssueAction() sdk.Action {
	return &GetIssueAction{}
}
