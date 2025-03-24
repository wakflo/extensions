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

package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/wakflo/extensions/internal/integrations/asana/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type listTasksActionProps struct {
	Project   *string `json:"project"`
	Workspace *string `json:"workspace"`
	Assignee  *string `json:"assignee"`
	Completed *bool   `json:"completed"`
	Limit     *int    `json:"limit"`
}

type ListTasksAction struct{}

func (l *ListTasksAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func (l ListTasksAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (l ListTasksAction) Name() string {
	return "List All Tasks"
}

func (l ListTasksAction) Description() string {
	return "Retrieve a list of tasks from Asana"
}

func (l ListTasksAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &listTasksDocs,
	}
}

func (l ListTasksAction) Icon() *string {
	return nil
}

func (l ListTasksAction) SampleData() sdkcore.JSON {
	return nil
}

func (l ListTasksAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"project": shared.GetProjectsInput(),
		"limit": autoform.NewNumberField().
			SetDisplayName("Limit").
			SetDescription("Maximum number of tasks to return").
			SetRequired(false).
			SetDefaultValue(50).
			Build(),
	}
}

func (l ListTasksAction) Auth() *sdk.Auth {
	return &sdk.Auth{
		Inherit: true,
	}
}

func (l ListTasksAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[listTasksActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	if input.Project == nil {
		return nil, errors.New("you must specify a project")
	}

	baseURL := shared.BaseAPI + "/tasks"
	params := url.Values{}

	if input.Project != nil {
		params.Add("project", *input.Project)
	}

	if input.Limit != nil {
		params.Add("limit", strconv.Itoa(*input.Limit))
	}

	params.Add("opt_fields", "gid,name,completed,due_on,notes,assignee,assignee_status,created_at,modified_at,projects,workspace")

	reqURL := baseURL
	if len(params) > 0 {
		reqURL = baseURL + "?" + params.Encode()
	}

	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+ctx.Auth.AccessToken)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		fmt.Println("Error response from Asana:", string(body))
		return nil, fmt.Errorf("request failed with status code: %d, response: %s", res.StatusCode, string(body))
	}

	var response interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func NewListTasksAction() sdk.Action {
	return &ListTasksAction{}
}
