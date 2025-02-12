package shared

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gookit/goutil/arrutil"

	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

var SharedAuth = autoform.NewAuth().NewCustomAuth().
	SetFields(map[string]*sdkcore.AutoFormSchema{
		"api-key": autoform.NewShortTextField().SetDisplayName("Api Key").
			SetDescription("The api key used to authenticate linear.").
			SetRequired(true).
			Build(),
	}).
	Build()

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

func GetTeamsInput() *sdkcore.AutoFormSchema {
	getTeams := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
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

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", ctx.Auth.Extra["api-key"])

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

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName("Teams").
		SetDescription("Select a team").
		SetDynamicOptions(&getTeams).
		SetRequired(true).Build()
}

func GetIssuesInput(title, description string) *sdkcore.AutoFormSchema {
	getIssues := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
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

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", ctx.Auth.Extra["api-key"])

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

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName(title).
		SetDescription(description).
		SetDynamicOptions(&getIssues).
		SetRequired(true).Build()
}

func GetPriorityInput(title, description string) *sdkcore.AutoFormSchema {
	getPriorities := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
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

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", ctx.Auth.Extra["api-key"])

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
				"id":   input.Priority,
				"name": input.Label,
			}, true
		})
		return ctx.Respond(items, len(items))
	}

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName(title).
		SetDescription(description).
		SetDynamicOptions(&getPriorities).
		SetRequired(false).Build()
}

func GetTeamLabelsInput(title, description string) *sdkcore.AutoFormSchema {
	getLabels := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
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

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", ctx.Auth.Extra["api-key"])

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

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName(title).
		SetDescription(description).
		SetDynamicOptions(&getLabels).
		SetRequired(false).Build()
}

func GetAssigneesInput(title, description string) *sdkcore.AutoFormSchema {
	getAssignees := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		query := `
		query Assignees {
		  issues {
		    nodes {
		      assignee {
		        id
		        name
		      }
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

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", ctx.Auth.Extra["api-key"])

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
				Issues struct {
					Nodes []struct {
						Assignee Assignee `json:"assignee"`
					} `json:"nodes"`
				} `json:"issues"`
			} `json:"data"`
		}

		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}

		// Create a unique set of assignees to avoid duplicates
		assigneeMap := make(map[string]Assignee)
		for _, node := range response.Data.Issues.Nodes {
			if node.Assignee.ID != "" {
				assigneeMap[node.Assignee.ID] = node.Assignee
			}
		}

		// Convert map to slice
		var assignees []Assignee
		for _, assignee := range assigneeMap {
			assignees = append(assignees, assignee)
		}
		return ctx.Respond(assignees, len(assignees))
	}

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName(title).
		SetDescription(description).
		SetDynamicOptions(&getAssignees).
		SetRequired(false).Build()
}

func GetIssueStatesInput(title, description string, required bool) *sdkcore.AutoFormSchema {
	getIssueStates := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		query := `{
          workflowStates {
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

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", ctx.Auth.Extra["api-key"])

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
				WorkflowStates struct {
					Nodes []WorkflowState `json:"nodes"`
				} `json:"workflowStates"`
			} `json:"data"`
		}

		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}

		issueStates := response.Data.WorkflowStates.Nodes
		return ctx.Respond(issueStates, len(issueStates))
	}

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName(title).
		SetDescription(description).
		SetDynamicOptions(&getIssueStates).
		SetRequired(required).Build()
}

func GetLabelsInput(title, description string) *sdkcore.AutoFormSchema {
	getLabels := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
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

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", ctx.Auth.Extra["api-key"])

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

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName(title).
		SetDescription(description).
		SetDynamicOptions(&getLabels).
		SetRequired(false).Build()
}
