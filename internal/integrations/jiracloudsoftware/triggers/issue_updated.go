package triggers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/jiracloudsoftware/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type issueUpdatedTriggerProps struct {
	ProjectID string `json:"projectId,omitempty"`
	IssueType string `json:"issueType,omitempty"`
}

type IssueUpdatedTrigger struct{}

// Metadata returns metadata about the trigger
func (t *IssueUpdatedTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "issue_updated",
		DisplayName:   "Issue Updated",
		Description:   "Trigger a workflow when an issue is updated in Jira",
		Type:          core.TriggerTypePolling,
		Documentation: issueUpdatedDocs,
		Icon:          "mdi:ticket-confirmation",
		SampleOutput: map[string]any{
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
		},
	}
}

// Auth returns the authentication requirements for the trigger
func (t *IssueUpdatedTrigger) Auth() *core.AuthMetadata {
	return nil
}

// GetType returns the type of the trigger
func (t *IssueUpdatedTrigger) GetType() core.TriggerType {
	return core.TriggerTypePolling
}

// Props returns the schema for the trigger's input configuration
func (t *IssueUpdatedTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("issue_updated", "Issue Updated")

	shared.RegisterProjectsProps(form)

	shared.RegisterIssueTypeProps(form, false)

	schema := form.Build()

	return schema
}

// Start initializes the trigger, required for event and webhook triggers in a lifecycle context
func (t *IssueUpdatedTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Stop shuts down the trigger, cleaning up resources and performing necessary teardown operations
func (t *IssueUpdatedTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Execute performs the trigger logic
func (t *IssueUpdatedTrigger) Execute(ctx sdkcontext.ExecuteContext) (core.JSON, error) {
	// Get the last run time
	var lastRunTime *string
	lr, err := ctx.GetMetadata("lastRun")
	if err == nil && lr != nil {
		formatted := lr.(*time.Time).Format("2006-01-02 15:04")
		lastRunTime = &formatted
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	email := authCtx.Extra["email"]
	apiToken := authCtx.Extra["api-token"]
	instanceURL := authCtx.Extra["instance-url"]

	var jql string
	if lastRunTime != nil {
		jql = fmt.Sprintf("updated >= '%s' AND updated > created", *lastRunTime)
	} else {
		jql = "updated >= '' AND updated > created"
	}

	input, err := sdk.InputToTypeSafely[issueUpdatedTriggerProps](ctx)
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

// Criteria returns the criteria for triggering this trigger
func (t *IssueUpdatedTrigger) Criteria(ctx context.Context) core.TriggerCriteria {
	return core.TriggerCriteria{}
}

// SampleData returns sample data for this trigger
func (t *IssueUpdatedTrigger) SampleData() core.JSON {
	return map[string]any{
		"issues": []map[string]any{
			{
				"id":  "12345",
				"key": "PRJ-123",
				"fields": map[string]any{
					"summary":  "Updated issue",
					"status":   map[string]any{"name": "In Progress"},
					"created":  "2023-05-01T09:12:34.567Z",
					"updated":  "2023-05-05T14:23:45.678Z",
					"priority": map[string]any{"name": "High"},
				},
			},
		},
	}
}

func NewIssueUpdatedTrigger() sdk.Trigger {
	return &IssueUpdatedTrigger{}
}
