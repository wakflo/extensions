package actions

import (
	"errors"
	"fmt"
	"strings"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/linear/shared"

	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type findIssuesActionProps struct {
	TeamID      string `json:"team-id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	AssigneeID  string `json:"assignee-id"`
	LabelID     string `json:"label-id"`
	StateID     string `json:"state-id"`
}

type FindIssuesAction struct{}

func (a *FindIssuesAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "find_issues",
		DisplayName:   "Find Issues",
		Description:   "Find Issues: Automatically identifies and extracts issues from various data sources, such as logs, tickets, or databases, to provide a centralized view of problems and errors in your workflow.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: findIssuesDocs,
		SampleOutput: map[string]any{
			"nodes": []map[string]any{
				{
					"id":          "issue-123",
					"title":       "Sample Issue",
					"description": "Issue description",
					"state": map[string]any{
						"name": "In Progress",
					},
					"priority":  1,
					"createdAt": "2024-01-01T00:00:00.000Z",
					"assignee": map[string]any{
						"name": "John Doe",
					},
				},
			},
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *FindIssuesAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("find_issues", "Find Issues")

	shared.GetTeamsProp(form)

	form.TextField("title", "Filter by Issue Name").
		Placeholder("Issue Name").
		HelpText("Filter issue by name").
		Required(false)

	form.TextareaField("description", "Filter by Description").
		Placeholder("Filter by Description").
		HelpText("Return issues with certain keywords in the issue description").
		Required(false)

	shared.GetLabelsProp("label-id", "Filter by label", "Filter issue by label", form)

	shared.GetAssigneesProp("assignee-id", "Filter by assignees", "Filter issue by assignees", form)

	shared.GetIssueStatesProp("state-id", "Filter by state", "filter issue by state", false, form)

	schema := form.Build()

	return schema
}

func (a *FindIssuesAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[findIssuesActionProps](ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input: %w", err)
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, fmt.Errorf("failed to get auth context: %w", err)
	}

	apiKEY := authCtx.Key

	// Validate API key format
	if !strings.HasPrefix(apiKEY, "lin_api_") {
		return nil, errors.New("invalid Linear API key: must start with 'lin_api_'")
	}

	var filters []string

	// Add team filter if provided
	if input.TeamID != "" {
		filters = append(filters, fmt.Sprintf(`team: {
			id: {
				eq: "%s"
			}
		}`, input.TeamID))
	}

	if input.AssigneeID != "" {
		filters = append(filters, fmt.Sprintf(`assignee: {
			id: {
				eq: "%s"
			}
		}`, input.AssigneeID))
	}

	if input.Description != "" {
		filters = append(filters, fmt.Sprintf(`description: {
			containsIgnoreCase: "%s"
		}`, escapeQuotes(input.Description)))
	}

	if input.Title != "" {
		filters = append(filters, fmt.Sprintf(`title: {
			containsIgnoreCase: "%s"
		}`, escapeQuotes(input.Title)))
	}

	if input.LabelID != "" {
		filters = append(filters, fmt.Sprintf(`labels: {
			some: {
				id: {
					eq: "%s"
				}
			}
		}`, input.LabelID))
	}

	if input.StateID != "" {
		filters = append(filters, fmt.Sprintf(`state: {
			id: {
				eq: "%s"
			}
		}`, input.StateID))
	}

	// Build query
	filterString := ""
	if len(filters) > 0 {
		filterString = strings.Join(filters, ", ")
	}

	query := fmt.Sprintf(`{
		issues(filter: {%s}) {
			nodes {
				id
				title
				description
				state {
					id
					name
				}
				priority
				priorityLabel
				createdAt
				updatedAt
				assignee {
					id
					name
				}
				labels {
					nodes {
						id
						name
					}
				}
				team {
					id
					name
				}
			}
		}
	}`, filterString)

	response, err := shared.MakeGraphQLRequest(apiKEY, query)
	if err != nil {
		return nil, fmt.Errorf("error making GraphQL request: %w", err)
	}

	nodes, ok := response["data"].(map[string]interface{})["issues"]
	if !ok {
		return nil, errors.New("failed to extract issues from response")
	}

	return nodes, nil
}

func (a *FindIssuesAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewFindIssuesAction() sdk.Action {
	return &FindIssuesAction{}
}
