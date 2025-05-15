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

type issueCreatedTriggerProps struct {
	ProjectID string `json:"projectId,omitempty"`
	IssueType string `json:"issueType,omitempty"`
}

type IssueCreatedTrigger struct{}

// Metadata returns metadata about the trigger
func (t *IssueCreatedTrigger) Metadata() sdk.TriggerMetadata {
	return sdk.TriggerMetadata{
		ID:            "issue_created",
		DisplayName:   "Issue Created",
		Description:   "Trigger a workflow when a new issue is created in Jira",
		Type:          core.TriggerTypePolling,
		Documentation: issueCreatedDocs,
		Icon:          "mdi:ticket-account",
		SampleOutput: map[string]any{
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
		},
	}
}

// Auth returns the authentication requirements for the trigger
func (t *IssueCreatedTrigger) Auth() *core.AuthMetadata {
	return nil
}

// GetType returns the type of the trigger
func (t *IssueCreatedTrigger) GetType() core.TriggerType {
	return core.TriggerTypePolling
}

// Props returns the schema for the trigger's input configuration
func (t *IssueCreatedTrigger) Props() *smartform.FormSchema {
	form := smartform.NewForm("issue_created", "Issue Created")
	shared.RegisterProjectsProps(form)
	shared.RegisterIssueTypeProps(form, false)

	schema := form.Build()

	return schema
}

// Start initializes the trigger, required for event and webhook triggers in a lifecycle context
func (t *IssueCreatedTrigger) Start(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Stop shuts down the trigger, cleaning up resources and performing necessary teardown operations
func (t *IssueCreatedTrigger) Stop(ctx sdkcontext.LifecycleContext) error {
	return nil
}

// Execute performs the trigger logic
func (t *IssueCreatedTrigger) Execute(ctx sdkcontext.ExecuteContext) (core.JSON, error) {
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
		jql = fmt.Sprintf("created >= '%s'", *lastRunTime)
	} else {
		jql = "created >= ''"
	}

	// Get input properties if any are provided
	input, err := sdk.InputToTypeSafely[issueCreatedTriggerProps](ctx)
	if err == nil {
		if input.ProjectID != "" {
			jql += " AND project = " + input.ProjectID
		}

		if input.IssueType != "" {
			jql += " AND issuetype = " + input.IssueType
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

// Criteria returns the criteria for triggering this trigger
func (t *IssueCreatedTrigger) Criteria(ctx context.Context) core.TriggerCriteria {
	return core.TriggerCriteria{}
}

// SampleData returns sample data for this trigger
func (t *IssueCreatedTrigger) SampleData() core.JSON {
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
					"created": "2023-05-05T12:34:56.789Z",
				},
			},
		},
	}
}

func NewIssueCreatedTrigger() sdk.Trigger {
	return &IssueCreatedTrigger{}
}
