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
			"message": "Hello World!",
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *FindIssuesAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("find_issues", "Find Issues")

	form.TextField("title", "Filter by Issue Name").
		Placeholder("Issue Name").
		HelpText("Filter issue by name").
		Required(false)

	form.TextareaField("description", "Filter by Description").
		Placeholder("Filter by Description").
		HelpText("Return issues wth certain keywords in the issue description").
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
		return nil, err
	}

	// Get the auth context
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	apiKEY := authCtx.Extra["api-key"]

	if !strings.HasPrefix(apiKEY, "lin_api_") {
		return nil, errors.New("invalid Linear API key: must start with 'lin_api_'")
	}

	var filters []string

	if input.AssigneeID != "" {
		filters = append(filters, fmt.Sprintf(`assignee: {
      id: {
        eq:"%s"
      }
    }`, input.AssigneeID))
	}
	if input.Description != "" {
		filters = append(filters, fmt.Sprintf(`description:{
      containsIgnoreCase: "%s"
    }`, input.Description))
	}
	if input.Title != "" {
		filters = append(filters, fmt.Sprintf(`title:{
      containsIgnoreCase: "%s"
    }`, input.Title))
	}
	if input.LabelID != "" {
		filters = append(filters, fmt.Sprintf(`labelIds: ["%s"]`, input.LabelID))
	}
	if input.StateID != "" {
		filters = append(filters, fmt.Sprintf(`  state: {
        id: {
			eq: "%s"
		}


      }`, input.StateID))
	}

	filterString := strings.Join(filters, ", ")

	query := fmt.Sprintf(`
	{
		issues(filter: {%s}) {
			nodes {
				id
				title
				description
				state {
					name
				}
				priority
				createdAt
				assignee {
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
