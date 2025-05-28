package shared

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gookit/goutil/arrutil"

	"github.com/juicycleff/smartform/v1"

	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

var (
	form = smartform.NewAuthForm("linear-auth", "Linear Oauth", smartform.AuthStrategyAPIKey)

	_ = form.APIKeyField("key", "Api Key").
		HelpText("The api key used to authenticate linear.").
		Required(true).
		Build()

	SharedAuth = form.Build()
)

const baseURL = "https://api.linear.app/graphql"

func MakeGraphQLRequest(apiKEY, query string) (map[string]interface{}, error) {
	payload := map[string]string{
		"query": query,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", apiKEY)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return result, nil
}

func GetTeamsProp(form *smartform.FormBuilder) *smartform.FieldBuilder {
	getTeams := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		query := `{
		 teams {
			nodes {
		     id
		     name
		   }
		 }
		}`

		queryBody := map[string]string{
			"query": query,
		}
		jsonQuery, err := json.Marshal(queryBody)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewBuffer(jsonQuery))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Get the auth context
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", authCtx.Extra["api-key"])

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to make request: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		var response struct {
			Data struct {
				Team struct {
					Nodes []Team `json:"nodes"`
				} `json:"teams"`
			} `json:"data"`
		}

		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}

		teams := response.Data.Team.Nodes
		return ctx.Respond(teams, len(teams))
	}

	return form.SelectField("team-id", "Teams").
		Placeholder("Select a team").
		Required(true).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getTeams)).
				WithSearchSupport().
				WithPagination(10).
				End().
				GetDynamicSource(),
		).
		HelpText("Select a team")
}

func GetIssuesProp(id string, title string, description string, form *smartform.FormBuilder) *smartform.FieldBuilder {
	getIssues := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		input := sdk.DynamicInputToType[struct {
			TeamID string `json:"team-id"`
		}](ctx)

		query := fmt.Sprintf(`{
		   team(id: "%s") {
			issues {
			  nodes {
			   id
			   title
			  }
			}
		   }
        }`, input.TeamID)

		queryBody := map[string]string{
			"query": query,
		}
		jsonQuery, err := json.Marshal(queryBody)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewBuffer(jsonQuery))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Get the auth context
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", authCtx.Extra["api-key"])

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to make request: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		var response struct {
			Data struct {
				Team struct {
					Issues struct {
						Nodes []Issue `json:"nodes"`
					} `json:"issues"`
				} `json:"team"`
			} `json:"data"`
		}

		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}

		issues := response.Data.Team.Issues.Nodes

		items := arrutil.Map[Issue, map[string]any](issues, func(input Issue) (target map[string]any, find bool) {
			return map[string]any{
				"id":   input.ID,
				"name": input.Title,
			}, true
		})
		return ctx.Respond(items, len(items))
	}

	return form.SelectField(id, title).
		Placeholder(description).
		Required(true).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getIssues)).
				WithSearchSupport().
				WithPagination(10).
				End().
				RefreshOn("team-id").
				GetDynamicSource(),
		).
		HelpText(description)
}

func GetPriorityProp(id string, title string, description string, form *smartform.FormBuilder) *smartform.FieldBuilder {
	getPriorities := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		query := `
		{
			issuePriorityValues {
				priority
				label
			}
		}`

		queryBody := map[string]string{
			"query": query,
		}
		jsonQuery, err := json.Marshal(queryBody)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewBuffer(jsonQuery))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Get the auth context
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", authCtx.Extra["api-key"])

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to make request: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		var response struct {
			Data struct {
				IssuePriorityValues []struct {
					Priority int    `json:"priority"`
					Label    string `json:"label"`
				} `json:"issuePriorityValues"`
			} `json:"data"`
		}

		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}

		priorities := response.Data.IssuePriorityValues

		items := arrutil.Map[struct {
			Priority int
			Label    string
		}, map[string]any]([]struct {
			Priority int
			Label    string
		}(priorities), func(input struct {
			Priority int
			Label    string
		},
		) (target map[string]any, find bool) {
			return map[string]any{
				"id":   strconv.Itoa(input.Priority),
				"name": input.Label,
			}, true
		})
		return ctx.Respond(items, len(items))
	}

	return form.SelectField(id, title).
		Placeholder(description).
		Required(false).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getPriorities)).
				WithSearchSupport().
				WithPagination(10).
				End().
				GetDynamicSource(),
		).
		HelpText(description)
}

func GetTeamLabelsProp(id string, title string, description string, form *smartform.FormBuilder) *smartform.FieldBuilder {
	getLabels := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		input := sdk.DynamicInputToType[struct {
			TeamID string `json:"team-id"`
		}](ctx)
		query := fmt.Sprintf(`
			query {
				team(id: "%s") {
					labels {
						nodes {
							id
							name
						}
					}
				}
			}
`, input.TeamID)

		queryBody := map[string]string{
			"query": query,
		}
		jsonQuery, err := json.Marshal(queryBody)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewBuffer(jsonQuery))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Get the auth context
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", authCtx.Extra["api-key"])

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to make request: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		var response struct {
			Data struct {
				Team struct {
					Labels struct {
						Nodes []Label `json:"nodes"`
					} `json:"labels"`
				} `json:"team"`
			} `json:"data"`
		}

		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}

		labels := response.Data.Team.Labels.Nodes
		return ctx.Respond(labels, len(labels))
	}

	return form.SelectField(id, title).
		Placeholder(description).
		Required(false).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getLabels)).
				WithSearchSupport().
				WithPagination(10).
				End().
				RefreshOn("team-id").
				GetDynamicSource(),
		).
		HelpText(description)
}

func GetAssigneesProp(id string, title string, description string, form *smartform.FormBuilder) *smartform.FieldBuilder {
	getAssignees := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		// Get the team ID from the form context
		input := sdk.DynamicInputToType[struct {
			TeamID string `json:"team-id"`
		}](ctx)

		// Query for team members instead of issue assignees
		query := fmt.Sprintf(`{
			team(id: "%s") {
				members {
					nodes {
						id
						name
						email
					}
				}
			}
		}`, input.TeamID)

		queryBody := map[string]string{
			"query": query,
		}
		jsonQuery, err := json.Marshal(queryBody)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewBuffer(jsonQuery))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Get the auth context
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", authCtx.Extra["api-key"])

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to make request: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		var response struct {
			Data struct {
				Team struct {
					Members struct {
						Nodes []struct {
							ID    string `json:"id"`
							Name  string `json:"name"`
							Email string `json:"email"`
						} `json:"nodes"`
					} `json:"members"`
				} `json:"team"`
			} `json:"data"`
		}

		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}

		// Convert to Assignee type
		var assignees []Assignee
		for _, member := range response.Data.Team.Members.Nodes {
			assignees = append(assignees, Assignee{
				ID:   member.ID,
				Name: member.Name,
			})
		}

		return ctx.Respond(assignees, len(assignees))
	}

	return form.SelectField(id, title).
		Placeholder(description).
		Required(false).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getAssignees)).
				WithSearchSupport().
				WithPagination(10).
				End().
				RefreshOn("team-id").
				GetDynamicSource(),
		).
		HelpText(description)
}

func GetIssueStatesProp(id string, title string, description string, required bool, form *smartform.FormBuilder) *smartform.FieldBuilder {
	getIssueStates := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		// Get the team ID from the form context
		input := sdk.DynamicInputToType[struct {
			TeamID string `json:"team-id"`
		}](ctx)

		// Query for team-specific workflow states
		query := fmt.Sprintf(`{
			team(id: "%s") {
				states {
					nodes {
						id
						name
						type
					}
				}
			}
		}`, input.TeamID)

		queryBody := map[string]string{
			"query": query,
		}
		jsonQuery, err := json.Marshal(queryBody)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewBuffer(jsonQuery))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Get the auth context
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", authCtx.Extra["api-key"])

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to make request: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		var response struct {
			Data struct {
				Team struct {
					States struct {
						Nodes []WorkflowState `json:"nodes"`
					} `json:"states"`
				} `json:"team"`
			} `json:"data"`
		}

		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}

		issueStates := response.Data.Team.States.Nodes
		return ctx.Respond(issueStates, len(issueStates))
	}

	return form.SelectField(id, title).
		Placeholder(description).
		Required(required).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getIssueStates)).
				WithSearchSupport().
				WithPagination(10).
				End().
				RefreshOn("team-id").
				GetDynamicSource(),
		).
		HelpText(description)
}

func GetLabelsProp(id string, title string, description string, form *smartform.FormBuilder) *smartform.FieldBuilder {
	getLabels := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		query := `{
          issueLabels {
    		nodes {
              id
              name
            }
		  }
        }`

		queryBody := map[string]string{
			"query": query,
		}
		jsonQuery, err := json.Marshal(queryBody)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewBuffer(jsonQuery))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Get the auth context
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", authCtx.Extra["api-key"])

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to make request: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		var response struct {
			Data struct {
				IssueLabels struct {
					Nodes []Label `json:"nodes"`
				} `json:"issueLabels"`
			} `json:"data"`
		}

		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}

		labels := response.Data.IssueLabels.Nodes
		return ctx.Respond(labels, len(labels))
	}

	return form.SelectField(id, title).
		Placeholder(description).
		Required(false).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getLabels)).
				WithSearchSupport().
				WithPagination(10).
				End().
				GetDynamicSource(),
		).
		HelpText(description)
}
