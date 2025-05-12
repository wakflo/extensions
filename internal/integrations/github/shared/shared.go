// Copyright 2022-present Wakflo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package shared

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gookit/goutil/arrutil"
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/autoform"
	"github.com/wakflo/go-sdk/core"

	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

var (
	form = smartform.NewAuthForm("github-auth", "GitHub OAuth", smartform.AuthStrategyOAuth2)
	_    = form.
		OAuthField("oauth", "GitHub OAuth").
		AuthorizationURL("https://github.com/login/oauth/authorize").
		TokenURL("https://github.com/login/oauth/access_token").
		Scopes([]string{"admin:repo_hook", "admin:org", "repo"}).
		Build()

	SharedGithubAuth = form.Build()
)

const baseURL = "https://api.github.com/graphql"

func GithubGQL(accessToken, query string) (map[string]interface{}, error) {
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

	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", "Bearer "+accessToken)

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
	if errs := json.Unmarshal(body, &result); errs != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", errs)
	}

	return result, nil
}

func GetRepositoryInput() *sdkcore.AutoFormSchema {
	getRepository := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		query := `{
		  viewer {
		    repositories(first: 100) {
		      nodes {
		        name
		        id
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

		req.Header.Add("Accept", "application/vnd.github+json")
		req.Header.Add("Authorization", "Bearer "+ctx.Auth.AccessToken)

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
				Viewer struct {
					Repositories struct {
						Nodes []struct {
							Name string `json:"name"`
							ID   string `json:"id"`
						} `json:"nodes"`
					} `json:"repositories"`
				} `json:"viewer"`
			} `json:"data"`
		}

		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}

		repositories := response.Data.Viewer.Repositories.Nodes

		return ctx.Respond(repositories, len(repositories))
	}

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName("Repository").
		SetDescription("Select a repository").
		SetDynamicOptions(&getRepository).
		SetRequired(true).Build()
}

func GetLabelInput() *sdkcore.AutoFormSchema {
	getLabels := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		input := sdk.DynamicInputToType[struct {
			Repository string `json:"repository"`
		}](ctx)
		query := fmt.Sprintf(` {
		  node(id: "%s") {
		    ... on Repository {
		      labels(first: 100) {
		        nodes {
		          name
		          id
		        }
		      }
		    }
		  }
		}`, input.Repository)

		queryBody := map[string]interface{}{
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

		req.Header.Set("Authorization", "Bearer "+ctx.Auth.AccessToken)
		req.Header.Add("Accept", "application/vnd.github+json")

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
				Node struct {
					Labels struct {
						Nodes []struct {
							Name string `json:"name"`
							ID   string `json:"id"`
						} `json:"nodes"`
					} `json:"labels"`
				} `json:"node"`
			} `json:"data"`
		}

		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}

		labels := response.Data.Node.Labels.Nodes
		return ctx.Respond(labels, len(labels))
	}

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName("Labels").
		SetDescription("Select labels for the issue").
		SetDynamicOptions(&getLabels).
		SetRequired(false).
		Build()
}

func GetIssuesInput() *sdkcore.AutoFormSchema {
	getIssues := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		input := sdk.DynamicInputToType[struct {
			Repository string `json:"repository"`
		}](ctx)
		query := fmt.Sprintf(` {
		  node(id: "%s") {
		    ... on Repository {
		      issues(first:100){
  				nodes{
    				id
					title
  				}
			  }
		    }
		  }
		}`, input.Repository)

		queryBody := map[string]interface{}{
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

		req.Header.Set("Authorization", "Bearer "+ctx.Auth.AccessToken)
		req.Header.Add("Accept", "application/vnd.github+json")

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
				Node struct {
					Issues struct {
						Nodes []struct {
							Title string `json:"title"`
							ID    string `json:"id"`
						} `json:"nodes"`
					} `json:"issues"`
				} `json:"node"`
			} `json:"data"`
		}

		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}

		issues := arrutil.Map(response.Data.Node.Issues.Nodes, func(issue struct {
			Title string `json:"title"`
			ID    string `json:"id"`
		},
		) (map[string]any, bool) {
			return map[string]any{
				"id":   issue.ID,
				"name": issue.Title,
			}, true
		})

		return ctx.Respond(issues, len(issues))
	}

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName("Issues").
		SetDescription("Select issue").
		SetDynamicOptions(&getIssues).
		SetRequired(true).
		Build()
}

// func getAssigneeInput() *sdkcore.AutoFormSchema {
//	getAssignees := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
//		input := sdk.DynamicInputToType[struct {
//			Repository string `json:"repository"`
//		}](ctx)
//		query := fmt.Sprintf(` {
//		  node(id: "%s") {
//		    ... on Repository {
//		      assignableUsers(first: 100) {
//		        nodes {
//		          login
//		          id
//		        }
//		      }
//		    }
//		  }
//		}`, input.Repository)
//
//		queryBody := map[string]interface{}{
//			"query": query,
//		}
//		jsonQuery, err := json.Marshal(queryBody)
//		if err != nil {
//			return nil, err
//		}
//
//		req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewBuffer(jsonQuery))
//		if err != nil {
//			return nil, fmt.Errorf("failed to create request: %w", err)
//		}
//
//		req.Header.Set("Authorization", "Bearer "+ctx.Auth.AccessToken)
//		req.Header.Add("Accept", "application/vnd.github+json")
//
//		client := &http.Client{}
//		resp, err := client.Do(req)
//		if err != nil {
//			return nil, fmt.Errorf("failed to make request: %w", err)
//		}
//		defer resp.Body.Close()
//
//		body, err := io.ReadAll(resp.Body)
//		if err != nil {
//			return nil, fmt.Errorf("failed to read response body: %w", err)
//		}
//
//		var response struct {
//			Data struct {
//				Node struct {
//					AssignableUsers struct {
//						Nodes []struct {
//							Login string `json:"login"`
//							ID    string `json:"id"`
//						} `json:"nodes"`
//					} `json:"assignableUsers"`
//				} `json:"node"`
//			} `json:"data"`
//		}
//
//		err = json.Unmarshal(body, &response)
//		if err != nil {
//			return nil, err
//		}
//
//		assignees := arrutil.Map(response.Data.Node.AssignableUsers.Nodes, func(user struct {
//			Login string `json:"login"`
//			ID    string `json:"id"`
//		},
//		) (map[string]any, bool) {
//			return map[string]any{
//				"id":   user.ID,
//				"name": user.Login,
//			}, true
//		})
//
//		return &assignees, nil
//	}
//
//	return autoform.NewDynamicField(sdkcore.String).
//		SetDisplayName("Assignees").
//		SetDescription("Select assignees for the issue").
//		SetDynamicOptions(&getAssignees).
//		SetRequired(false).
//		Build()
//  }

var LockIssueReason = []*sdkcore.AutoFormSchema{
	{Const: " OFF_TOPIC", Title: "Off topic"},
	{Const: "TOO_HEATED", Title: "Too heated"},
	{Const: "RESOLVED", Title: "resolved"},
	{Const: "SPAM", Title: "Spam"},
}

// GetRepositories is a dynamic field function that retrieves a list of repositories
func GetRepositories(ctx sdkcontext.DynamicFieldContext) (*core.DynamicOptionsResponse, error) {
	query := `{
	  viewer {
	    repositories(first: 100) {
	      nodes {
	        name
	        id
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

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewBuffer(jsonQuery))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", "Bearer "+authCtx.AccessToken)

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
			Viewer struct {
				Repositories struct {
					Nodes []struct {
						Name string `json:"name"`
						ID   string `json:"id"`
					} `json:"nodes"`
				} `json:"repositories"`
			} `json:"viewer"`
		} `json:"data"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	repositories := response.Data.Viewer.Repositories.Nodes

	return ctx.Respond(repositories, len(repositories))
}

// GetIssues is a dynamic field function that retrieves a list of issues for a repository
func GetIssues(ctx sdkcontext.DynamicFieldContext) (*core.DynamicOptionsResponse, error) {
	input := sdk.DynamicInputToType[struct {
		Repository string `json:"repository"`
	}](ctx)

	query := fmt.Sprintf(` {
	  node(id: "%s") {
	    ... on Repository {
	      issues(first: 100) {
	        nodes {
	          title
	          id
	        }
	      }
	    }
	  }
	}`, input.Repository)

	queryBody := map[string]interface{}{
		"query": query,
	}
	jsonQuery, err := json.Marshal(queryBody)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewBuffer(jsonQuery))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+authCtx.AccessToken)
	req.Header.Add("Accept", "application/vnd.github+json")

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
			Node struct {
				Issues struct {
					Nodes []struct {
						Title string `json:"title"`
						ID    string `json:"id"`
					} `json:"nodes"`
				} `json:"issues"`
			} `json:"node"`
		} `json:"data"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	issues := response.Data.Node.Issues.Nodes

	return ctx.Respond(issues, len(issues))
}
