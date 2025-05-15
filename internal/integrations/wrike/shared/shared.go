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
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gookit/goutil/arrutil"
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"

	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

var (
	// #nosec
	tokenURL = "https://login.wrike.com/oauth2/token"
)

var form = smartform.NewAuthForm("wrike-auth", "Wrike Oauth", smartform.AuthStrategyOAuth2)
var _ = form.OAuthField("oauth", "Wrike Oauth").
	AuthorizationURL("https://login.wrike.com/oauth2/authorize/v4").
	TokenURL(tokenURL).
	Scopes([]string{"Default, wsReadWrite, wsReadOnly"}).
	Build()

var (
	WrikeSharedAuth = form.Build()
)

func GetFoldersProp(form *smartform.FormBuilder) *smartform.FieldBuilder {
	getFolders := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		baseURL := WrikeAPIBaseURL + "/folders"
		params := url.Values{}

		reqURL := baseURL + "?" + params.Encode()

		req, err := http.NewRequest(http.MethodGet, reqURL, nil)
		if err != nil {
			return nil, err
		}

		tokenSource := ctx.Auth().Token
		if tokenSource == nil {
			return nil, errors.New("missing authentication token")
		}
		token := tokenSource.AccessToken

		req.Header.Add("Accept", "application/json")
		req.Header.Add("Authorization", "Bearer "+token)

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if res.StatusCode < 200 || res.StatusCode >= 300 {
			return nil, errors.New("request failed with status: " + strconv.Itoa(res.StatusCode))
		}

		type Folder struct {
			ID      string `json:"id"`
			Title   string `json:"title"`
			Color   string `json:"color,omitempty"`
			Space   bool   `json:"space,omitempty"`
			Project struct {
				Status       string `json:"status,omitempty"`
				CustomStatus struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				} `json:"customStatus,omitempty"`
			} `json:"project,omitempty"`
		}

		type FolderResponse struct {
			Kind string   `json:"kind"`
			Data []Folder `json:"data"`
		}

		var response FolderResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}

		folders := response.Data

		items := arrutil.Map[Folder, map[string]any](folders, func(input Folder) (target map[string]any, find bool) {
			return map[string]any{
				"value": input.ID,
				"label": input.Title,
			}, true
		})

		return ctx.Respond(items, len(items))
	}

	return form.SelectField("folderId", "Folder Id").
		Placeholder("Select a folder").
		Required(false).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getFolders)).
				WithSearchSupport().
				WithPagination(10).
				End().
				GetDynamicSource(),
		).
		HelpText("Select a folder")
}

func GetTaskProp(form *smartform.FormBuilder) *smartform.FieldBuilder {
	getTask := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		input := sdk.DynamicInputToType[struct {
			FolderID string `json:"folder_id,omitempty"`
		}](ctx)

		baseURL := WrikeAPIBaseURL + "/tasks"
		params := url.Values{}

		if input.FolderID != "" {
			params.Add("folder", input.FolderID)
		}

		reqURL := baseURL + "?" + params.Encode()

		req, err := http.NewRequest(http.MethodGet, reqURL, nil)
		if err != nil {
			return nil, err
		}

		tokenSource := ctx.Auth().Token
		if tokenSource == nil {
			return nil, errors.New("missing authentication token")
		}
		token := tokenSource.AccessToken

		req.Header.Add("Accept", "application/json")
		req.Header.Add("Authorization", "Bearer "+token)

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if res.StatusCode < 200 || res.StatusCode >= 300 {
			return nil, errors.New("request failed with status: " + strconv.Itoa(res.StatusCode))
		}

		type Task struct {
			ID         string `json:"id"`
			Title      string `json:"title"`
			Status     string `json:"status,omitempty"`
			Importance string `json:"importance,omitempty"`
			Dates      struct {
				Due   string `json:"due,omitempty"`
				Start string `json:"start,omitempty"`
			} `json:"dates,omitempty"`
		}

		type TaskResponse struct {
			Data []Task `json:"data"`
		}

		var response TaskResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, err
		}

		tasks := response.Data

		items := arrutil.Map[Task, map[string]any](tasks, func(input Task) (target map[string]any, find bool) {
			return map[string]any{
				"id":   input.ID,
				"name": input.Title,
			}, true
		})

		return ctx.Respond(items, len(items))
	}

	return form.SelectField("task_id", "Task").
		Placeholder("Select a task").
		Required(true).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(
					sdk.WithDynamicFunctionCalling(&getTask)).
				WithSearchSupport().
				WithPagination(10).
				End().
				GetDynamicSource(),
		).
		HelpText("Select a task")
}
