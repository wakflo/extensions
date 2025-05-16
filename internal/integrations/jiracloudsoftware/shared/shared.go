package shared

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gookit/goutil/arrutil"
	"github.com/juicycleff/smartform/v1"

	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"

	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

var (
	form = smartform.NewAuthForm("jiracloud-auth", "Jira Cloud Software Authentication", smartform.AuthStrategyCustom)

	_ = form.TextField("instance-url", "Instance URL (Required)").
		Required(true).
		HelpText("The link of your Jira instance (e.g https://example.atlassian.net)")

	_ = form.TextField("email", "Email (Required)").
		Required(true).
		HelpText("The email you use to login to Jira")

	_ = form.TextField("api-token", "Your Jira API Token").
		Required(true).
		HelpText("Your Jira API Token")

	JiraSharedAuth = form.Build()
)

func JiraRequest(email, apiToken, reqURL, method, message string, request []byte) (interface{}, error) {
	auth := email + ":" + apiToken

	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))

	authHeader := "Basic " + encodedAuth

	req, err := http.NewRequest(method, reqURL, bytes.NewBuffer(request))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Add("Authorization", authHeader)
	client := &http.Client{}
	res, errs := client.Do(req)
	if errs != nil {
		return nil, errs
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusNoContent {
		return map[string]interface{}{
			"Result": message,
		}, nil
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, fmt.Errorf("received status code %d: %s", res.StatusCode, string(body))
	}

	var response interface{}
	if newErrs := json.Unmarshal(body, &response); newErrs != nil {
		return nil, errors.New("error parsing response")
	}

	return response, nil
}

func RegisterUsersProps(form *smartform.FormBuilder) {
	getUsers := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		auth := authCtx.Extra["email"] + ":" + authCtx.Extra["api-token"]

		encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))

		authHeader := "Basic " + encodedAuth

		baseAPI := authCtx.Extra["instance-url"] + "/rest/api/2/users/search"

		req, err := http.NewRequest(http.MethodGet, baseAPI, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Authorization", authHeader)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		rsp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer rsp.Body.Close()

		newBytes, err := io.ReadAll(rsp.Body)
		if err != nil {
			return nil, err
		}

		var users []User
		err = json.Unmarshal(newBytes, &users)
		if err != nil {
			return nil, err
		}

		atlassianUsers := arrutil.Filter[User](users, func(input User) bool {
			return input.AccountType == "atlassian"
		})

		items := arrutil.Map[User, map[string]any](atlassianUsers, func(input User) (target map[string]any, find bool) {
			return map[string]any{
				"id":   input.AccountID,
				"name": input.DisplayName,
			}, true
		})

		return ctx.Respond(items, len(items))
	}

	form.SelectField("assignee", "Assignees").
		Placeholder("Select an assignee").
		Required(false).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getUsers)).
				WithSearchSupport().
				WithPagination(10).
				End().
				GetDynamicSource(),
		).
		HelpText("Select an assignee")
}

func RegisterProjectsProps(form *smartform.FormBuilder) {
	getProjects := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		auth := authCtx.Extra["email"] + ":" + authCtx.Extra["api-token"]

		encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))

		authHeader := "Basic " + encodedAuth

		baseAPI := authCtx.Extra["instance-url"] + "/rest/api/2/project/search"

		req, err := http.NewRequest(http.MethodGet, baseAPI, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Authorization", authHeader)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		rsp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer rsp.Body.Close()

		newBytes, err := io.ReadAll(rsp.Body)
		if err != nil {
			return nil, err
		}

		var projects ProjectResponse
		err = json.Unmarshal(newBytes, &projects)
		if err != nil {
			return nil, err
		}

		return ctx.Respond(projects.Values, len(projects.Values))
	}

	form.SelectField("projectId", "Projects").
		Placeholder("Select a project").
		Required(true).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getProjects)).
				WithSearchSupport().
				WithPagination(10).
				End().
				RefreshOn("connection").
				GetDynamicSource(),
		).
		HelpText("Select a project")
}

func RegisterIssueTypeProps(form *smartform.FormBuilder, required bool) *smartform.FieldBuilder {
	getIssues := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		auth := authCtx.Extra["email"] + ":" + authCtx.Extra["api-token"]

		encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))

		authHeader := "Basic " + encodedAuth

		input := sdk.DynamicInputToType[struct {
			ProjectID string `json:"projectId"`
		}](ctx)

		baseAPI := authCtx.Extra["instance-url"] + "/rest/api/3/issuetype/project?projectId=" + input.ProjectID

		req, err := http.NewRequest(http.MethodGet, baseAPI, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Authorization", authHeader)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		client := &http.Client{}
		rsp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer rsp.Body.Close()

		newBytes, err := io.ReadAll(rsp.Body)
		if err != nil {
			return nil, err
		}

		var issueTypes []IssueType
		err = json.Unmarshal(newBytes, &issueTypes)
		if err != nil {
			return nil, err
		}

		return ctx.Respond(issueTypes, len(issueTypes))
	}

	return form.SelectField("IssueTypeId", "Issue Type").
		Placeholder("Select an issue type").
		Required(required).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getIssues)).
				WithFieldReference("projectId", "projectId").
				WithSearchSupport().
				WithPagination(10).
				End().
				RefreshOn("projectId").
				GetDynamicSource(),
		).
		HelpText("Select an issue type")
}

func RegisterTransitionsProps(form *smartform.FormBuilder) *smartform.FieldBuilder {
	getTransitions := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		auth := authCtx.Extra["email"] + ":" + authCtx.Extra["api-token"]
		encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
		authHeader := "Basic " + encodedAuth

		input := sdk.DynamicInputToType[struct {
			ProjectID string `json:"projectId"`
			IssueID   string `json:"issueId"`
		}](ctx)

		if input.IssueID == "" {
			return ctx.Respond([]map[string]any{}, 0)
		}

		baseAPI := authCtx.Extra["instance-url"] + "/rest/api/3/issue/" + input.IssueID + "/transitions"

		req, err := http.NewRequest(http.MethodGet, baseAPI, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Authorization", authHeader)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		client := &http.Client{}
		rsp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer rsp.Body.Close()

		responseBytes, err := io.ReadAll(rsp.Body)
		if err != nil {
			return nil, err
		}

		var transitionResponse struct {
			Transitions []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
				To   struct {
					Name string `json:"name"`
				} `json:"to"`
			} `json:"transitions"`
		}

		err = json.Unmarshal(responseBytes, &transitionResponse)
		if err != nil {
			return nil, err
		}

		items := make([]map[string]any, 0)
		for _, transition := range transitionResponse.Transitions {
			items = append(items, map[string]any{
				"id":   transition.ID,
				"name": transition.Name + " → " + transition.To.Name,
			})
		}

		return ctx.Respond(items, len(items))
	}

	return form.SelectField("transitionId", "Transition").
		Placeholder("Select a status transition").
		Required(true).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getTransitions)).
				WithFieldReference("issueId", "issueId").
				WithFieldReference("projectId", "projectId").
				WithSearchSupport().
				End().
				RefreshOn("issueId").
				GetDynamicSource(),
		).
		HelpText("Select a status transition")
}

func RegisterIssuesProps(form *smartform.FormBuilder) *smartform.FieldBuilder {
	getIssues := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		auth := authCtx.Extra["email"] + ":" + authCtx.Extra["api-token"]
		encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
		authHeader := "Basic " + encodedAuth

		baseAPI := authCtx.Extra["instance-url"] + "/rest/api/3/search"

		input := sdk.DynamicInputToType[struct {
			ProjectID string `json:"projectId"`
		}](ctx)

		body := map[string]interface{}{
			"jql":        "project=" + input.ProjectID,
			"fields":     []string{"summary"},
			"maxResults": 50,
		}

		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest(http.MethodPost, baseAPI, bytes.NewBuffer(bodyBytes))
		if err != nil {
			return nil, err
		}

		req.Header.Set("Authorization", authHeader)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		client := &http.Client{}
		rsp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer rsp.Body.Close()

		responseBytes, err := io.ReadAll(rsp.Body)
		if err != nil {
			return nil, err
		}

		var searchResponse SearchIssuesResponse
		err = json.Unmarshal(responseBytes, &searchResponse)
		if err != nil {
			return nil, err
		}

		items := arrutil.Map[Issue, map[string]any](searchResponse.Issues, func(issue Issue) (map[string]any, bool) {
			return map[string]any{
				"id":   issue.ID,
				"name": fmt.Sprintf("[%s] %s", issue.Key, issue.Fields.Summary),
			}, true
		})

		return ctx.Respond(items, len(items))
	}

	return form.SelectField("issues", "Issues").
		Placeholder("Select an issue").
		Required(false).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getIssues)).
				WithFieldReference("projectId", "projectId").
				WithSearchSupport().
				WithPagination(50).
				End().
				RefreshOn("projectId").
				GetDynamicSource(),
		).
		HelpText("Select an issue")
}
