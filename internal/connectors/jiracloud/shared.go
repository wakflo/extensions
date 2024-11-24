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

package jiracloud

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gookit/goutil/arrutil"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

var sharedAuth = autoform.NewCustomAuthField().
	SetFields(map[string]*sdkcore.AutoFormSchema{
		"instance-url": autoform.NewShortTextField().
			SetDisplayName("Instance URL (Required)").
			SetDescription("The link of your Jira instance (e.g https://example.atlassian.net)").
			Build(),
		"email": autoform.NewShortTextField().SetDisplayName("Email (Required)").
			SetDescription("The email you use to login to Jira").
			SetRequired(true).
			Build(),
		"api-token": autoform.NewLongTextField().SetDisplayName("Api Token (Required)").
			SetDescription("Your Jira API Token").
			SetRequired(true).
			Build(),
	}).
	Build()

func jiraRequest(email, apiToken, reqURL, method, message string, request []byte) (interface{}, error) {
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

func getUsersInput() *sdkcore.AutoFormSchema {
	getUsers := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		auth := ctx.Auth.Extra["email"] + ":" + ctx.Auth.Extra["api-token"]

		encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))

		authHeader := "Basic " + encodedAuth

		baseAPI := ctx.Auth.Extra["instance-url"] + "/rest/api/2/users/search"

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

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName("Assignees").
		SetDescription("Select an assignee").
		SetDynamicOptions(&getUsers).
		SetRequired(false).Build()
}

func getProjectsInput() *sdkcore.AutoFormSchema {
	getProjects := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		auth := ctx.Auth.Extra["email"] + ":" + ctx.Auth.Extra["api-token"]

		encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))

		authHeader := "Basic " + encodedAuth

		baseAPI := ctx.Auth.Extra["instance-url"] + "/rest/api/2/project/search"

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

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName("Projects").
		SetDescription("Select a project").
		SetDynamicOptions(&getProjects).
		SetRequired(true).Build()
}

func getIssueTypesInput() *sdkcore.AutoFormSchema {
	getIssues := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		auth := ctx.Auth.Extra["email"] + ":" + ctx.Auth.Extra["api-token"]

		encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))

		authHeader := "Basic " + encodedAuth

		input := sdk.DynamicInputToType[struct {
			ProjectID string `json:"projectId"`
		}](ctx)

		baseAPI := ctx.Auth.Extra["instance-url"] + "/rest/api/3/issuetype/project?projectId=" + input.ProjectID

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

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName("Issue Type").
		SetDescription("Select an issue type").
		SetDynamicOptions(&getIssues).
		SetRequired(false).Build()
}

func getIssuesInput() *sdkcore.AutoFormSchema {
	getIssues := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		auth := ctx.Auth.Extra["email"] + ":" + ctx.Auth.Extra["api-token"]
		encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
		authHeader := "Basic " + encodedAuth

		baseAPI := ctx.Auth.Extra["instance-url"] + "/rest/api/3/search"

		input := sdk.DynamicInputToType[struct {
			ProjectID   string `json:"projectId"`
			CommentBody string `json:"commentText"`
			IssueID     string `json:"issueId"`
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

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName("Issues").
		SetDescription("Select an issue").
		SetDynamicOptions(&getIssues).
		SetRequired(false).Build()
}

var PriorityLevels = []*sdkcore.AutoFormSchema{
	{Const: "1", Title: "Higher"},
	{Const: "2", Title: "High"},
	{Const: "3", Title: "Medium"},
	{Const: "4", Title: "Lower"},
	{Const: "5", Title: "Lower"},
}

var OrderBy = []*sdkcore.AutoFormSchema{
	{Const: "-created", Title: "Created (Descending)"},
	{Const: "+created", Title: "Created (Ascending)"},
}
