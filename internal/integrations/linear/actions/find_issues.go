package actions

import (
	"errors"
	"fmt"
	"strings"

	"github.com/wakflo/extensions/internal/integrations/linear/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type findIssuesActionProps struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	AssigneeID  string `json:"assignee-id"`
	LabelID     string `json:"label-id"`
	StateID     string `json:"state-id"`
}

type FindIssuesAction struct{}

func (a *FindIssuesAction) Name() string {
	return "Find Issues"
}

func (a *FindIssuesAction) Description() string {
	return "Automatically identifies and extracts issues from various data sources, such as logs, tickets, or databases, to provide a centralized view of problems and errors in your workflow."
}

func (a *FindIssuesAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *FindIssuesAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &findIssuesDocs,
	}
}

func (a *FindIssuesAction) Icon() *string {
	return nil
}

func (a *FindIssuesAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"title": autoform.NewShortTextField().
			SetDisplayName("Filter by Issue Name").
			SetDescription("Returns issues specified by the name").
			Build(),
		"description": autoform.NewLongTextField().
			SetDisplayName("Filter by Description").
			SetDescription("Return issues wth certain keywords in the issue description").
			Build(),
		"label-id":    shared.GetLabelsInput("Filter by label", "Filter issue by label"),
		"assignee-id": shared.GetAssigneesInput("Filter by assignees", "Filter issue by assignees"),
		"state-id":    shared.GetIssueStatesInput("Filter by state", "filter issue by state", false),
	}
}

func (a *FindIssuesAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[findIssuesActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing linear api key")
	}
	apiKEY := ctx.Auth.Extra["api-key"]

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

func (a *FindIssuesAction) Auth() *sdk.Auth {
	return nil
}

func (a *FindIssuesAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *FindIssuesAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewFindIssuesAction() sdk.Action {
	return &FindIssuesAction{}
}
