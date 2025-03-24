package triggers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/wakflo/extensions/internal/integrations/jiracloud/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type issueCreatedTriggerProps struct {
	ProjectID string `json:"projectId,omitempty"`
	IssueType string `json:"issueType,omitempty"`
}

type IssueCreatedTrigger struct{}

func (t *IssueCreatedTrigger) Name() string {
	return "Issue Created"
}

func (t *IssueCreatedTrigger) Description() string {
	return "Trigger a workflow when a new issue is created in Jira"
}

func (t *IssueCreatedTrigger) GetType() sdkcore.TriggerType {
	return sdkcore.TriggerTypePolling
}

func (t *IssueCreatedTrigger) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &issueCreatedDocs,
	}
}

func (t *IssueCreatedTrigger) Icon() *string {
	icon := "mdi:ticket-account"
	return &icon
}

func (t *IssueCreatedTrigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"projectId": shared.GetProjectsInput(),
		"issueType": shared.GetIssueTypesInput(),
	}
}

func (t *IssueCreatedTrigger) Start(ctx sdk.LifecycleContext) error {
	return nil
}

// Stop shuts down the trigger
func (t *IssueCreatedTrigger) Stop(ctx sdk.LifecycleContext) error {
	return nil
}

// Execute performs the trigger logic
func (t *IssueCreatedTrigger) Execute(ctx sdk.ExecuteContext) (sdkcore.JSON, error) {
	lastRunTime := ctx.Metadata().LastRun

	email := ctx.Auth.Extra["email"]
	apiToken := ctx.Auth.Extra["api-token"]
	instanceURL := ctx.Auth.Extra["instance-url"]

	var jql string
	if lastRunTime != nil {
		jql = fmt.Sprintf("created >= '%s'", lastRunTime.Format("2006-01-02 15:04"))
	} else {
		jql = "created >= ''"
	}

	// Get input properties if any are provided
	input, err := sdk.InputToTypeSafely[issueCreatedTriggerProps](ctx.BaseContext)
	if err == nil {
		if input.ProjectID != "" {
			jql += fmt.Sprintf(" AND project = %s", input.ProjectID)
		}

		if input.IssueType != "" {
			jql += fmt.Sprintf(" AND issuetype = %s", input.IssueType)
		}
	}

	// Order by creation date in descending order (newest first)
	jql += " ORDER BY created DESC"

	// Build the request payload for search
	requestBody := map[string]interface{}{
		"jql":        jql,
		"maxResults": 50,
		"fields":     []string{"summary", "description", "status", "creator", "created", "priority", "assignee"},
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

func (t *IssueCreatedTrigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *IssueCreatedTrigger) Auth() *sdk.Auth {
	return nil
}

func (t *IssueCreatedTrigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"issues": []map[string]any{
			{
				"id":   "12345",
				"key":  "PRJ-123",
				"self": "https://yourcompany.atlassian.net/rest/api/3/issue/12345",
				"fields": map[string]any{
					"summary": "New issue created",
					"description": map[string]any{
						"type": "doc",
						"content": []map[string]any{
							{
								"type": "paragraph",
								"content": []map[string]any{
									{
										"type": "text",
										"text": "This is a newly created issue",
									},
								},
							},
						},
					},
					"status": map[string]any{
						"name": "To Do",
					},
					"creator": map[string]any{
						"displayName":  "John Doe",
						"emailAddress": "john.doe@example.com",
					},
					"created": "2023-05-05T12:34:56.789Z",
					"priority": map[string]any{
						"name": "Medium",
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

func NewIssueCreatedTrigger() sdk.Trigger {
	return &IssueCreatedTrigger{}
}
