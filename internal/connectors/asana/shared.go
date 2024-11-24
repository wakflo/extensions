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

package asana

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/gookit/goutil/arrutil"
	fastshot "github.com/opus-domini/fast-shot"

	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
)

var (
	// #nosec
	tokenURL   = "https://app.asana.com/-/oauth_token"
	sharedAuth = autoform.NewOAuthField("https://app.asana.com/-/oauth_authorize", &tokenURL, []string{
		"default",
	}).SetRequired(true).Build()
)
var baseAPI = "https://app.asana.com/api/1.0"

func getWorkspacesInput() *sdkcore.AutoFormSchema {
	getWorkspaces := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		client := fastshot.NewClient(baseAPI).
			Auth().BearerToken(ctx.Auth.AccessToken).
			Header().
			AddAccept("application/json").
			Build()

		rsp, err := client.GET("/workspaces").Send()
		if err != nil {
			return nil, err
		}

		if rsp.Status().IsError() {
			return nil, errors.New(rsp.Status().Text())
		}

		bytes, err := io.ReadAll(rsp.Raw().Body) //nolint:bodyclose
		if err != nil {
			return nil, err
		}

		var workspaces WorkspaceResponse
		err = json.Unmarshal(bytes, &workspaces)
		if err != nil {
			return nil, err
		}

		workspace := workspaces.Data
		items := arrutil.Map[Workspace, map[string]any](workspace, func(input Workspace) (target map[string]any, find bool) {
			return map[string]any{
				"id":   input.GID,
				"name": input.Name,
			}, true
		})
		return ctx.Respond(items, len(items))
	}

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName("Workspace").
		SetDescription("Task workspace ID.").
		SetDependsOn([]string{"connection"}).
		SetDynamicOptions(&getWorkspaces).
		SetRequired(false).Build()
}

func getProjectsInput() *sdkcore.AutoFormSchema {
	getProjects := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		qu := fastshot.NewClient(baseAPI).
			Auth().BearerToken(ctx.Auth.AccessToken).
			Header().
			AddAccept("application/json").
			Build().GET("/projects")

		rsp, err := qu.Send()
		if err != nil {
			return nil, err
		}

		if rsp.Status().IsError() {
			return nil, errors.New(rsp.Status().Text())
		}
		defer rsp.Raw().Body.Close()

		bytes, err := io.ReadAll(rsp.Raw().Body) //nolint:bodyclose
		if err != nil {
			return nil, err
		}

		var response ProjectResponse
		err = json.Unmarshal(bytes, &response)
		if err != nil {
			return nil, err
		}

		projects := response.Data

		items := arrutil.Map[Project, map[string]any](projects, func(input Project) (target map[string]any, find bool) {
			return map[string]any{
				"id":   input.GID,
				"name": input.Name,
			}, true
		})

		return ctx.Respond(items, len(items))
	}

	return autoform.NewDynamicField(sdkcore.String).
		SetDisplayName("Project").
		SetDescription("A project to create the task under").
		SetDynamicOptions(&getProjects).
		SetRequired(false).Build()
}
