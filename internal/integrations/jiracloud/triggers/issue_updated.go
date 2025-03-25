package triggers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/wakflo/extensions/internal/integrations/jiracloud/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type issueUpdatedTriggerProps struct {
	ProjectID string `json:"projectId,omitempty"`
	IssueType string `json:"issueType,omitempty"`
}

type IssueUpdatedTrigger struct{}

func (t *IssueUpdatedTrigger) Name() string {
	return "Issue Updated"
}

func (t *IssueUpdatedTrigger) Description() string {
	return "Trigger a workflow when an issue is updated in Jira"
}

func (t *IssueUpdatedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *IssueUpdatedTrigger) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &issueUpdatedDocs,
	}
}

func (t *IssueUpdatedTrigger) Icon() *string {
	icon := "mdi:ticket-confirmation"
	return &icon
}

func (t *IssueUpdatedTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"projectId": shared.GetProjectsInput(),
		"issueType": shared.GetIssueTypesInput(),
	}
}

func (t *IssueUpdatedTrigger) Start(ctx sdk.LifecycleContext) error {
	return nil
}

func (t *IssueUpdatedTrigger) Stop(ctx sdk.LifecycleContext) error {
	return nil
}

func (t *IssueUpdatedTrigger) Execute(ctx sdk.ExecuteContext) (sdkcore.JSON, error) {
	lastRunTime := ctx.Metadata().LastRun

	email := ctx.Auth.Extra["email"]
	apiToken := ctx.Auth.Extra["api-token"]
	instanceURL := ctx.Auth.Extra["instance-url"]

	var jql string
	if lastRunTime != nil {
		jql = fmt.Sprintf("updated >= '%s' AND updated > created", lastRunTime.Format("2006-01-02 15:04"))
	} else {
		jql = fmt.Sprintf("updated >= '%s' AND updated > created", "")
	}

	input, err := sdk.InputToTypeSafely[issueUpdatedTriggerProps](ctx.BaseContext)
	if err == nil {
		if input.ProjectID != "" {
			jql += fmt.Sprintf(" AND project = %s", input.ProjectID)
		}

		if input.IssueType != "" {
			jql += fmt.Sprintf(" AND issuetype = %s", input.IssueType)
		}
	}

	jql += " ORDER BY updated DESC"

	requestBody := map[string]interface{}{
		"jql":        jql,
		"maxResults": 50,
		"fields":     []string{"summary", "description", "status", "creator", "created", "updated", "priority", "assignee"},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	response, err := shared.JiraRequest(
		email,
		apiToken,
		instanceURL+"/rest/api/3/search",
		"POST",
		"Issues retrieved successfully",
		jsonBody,
	)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}

	return response, nil
}

func (t *IssueUpdatedTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *IssueUpdatedTrigger) Auth() *sdk.Auth {
	return nil
}

func (t *IssueUpdatedTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"issues": []map[string]any{
			{
				"id":   "12345",
				"key":  "PRJ-123",
				"self": "https://yourcompany.atlassian.net/rest/api/3/issue/12345",
				"fields": map[string]any{
					"summary": "Updated issue",
					"description": map[string]any{
						"type": "doc",
						"content": []map[string]any{
							{
								"type": "paragraph",
								"content": []map[string]any{
									{
										"type": "text",
										"text": "This issue has been updated",
									},
								},
							},
						},
					},
					"status": map[string]any{
						"name": "In Progress",
					},
					"creator": map[string]any{
						"displayName":  "John Doe",
						"emailAddress": "john.doe@example.com",
					},
					"created": "2023-05-01T09:12:34.567Z",
					"updated": "2023-05-05T14:23:45.678Z",
					"priority": map[string]any{
						"name": "High",
					},
					"assignee": map[string]any{
						"displayName":  "Jane Smith",
						"emailAddress": "jane.smith@example.com",
					},
				},
			},
		},
		"total":      1,
		"maxResults": 50,
		"startAt":    0,
	}
}

func NewIssueUpdatedTrigger() sdk.Trigger {
	return &IssueUpdatedTrigger{}
}
