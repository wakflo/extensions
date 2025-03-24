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
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/wakflo/extensions/internal/integrations/asana/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type getTaskActionProps struct {
	TaskID    string `json:"task_id"`
	ProjectID string `json:"project_id"`
}

type GetTaskAction struct{}

func (g *GetTaskAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func (g GetTaskAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (g GetTaskAction) Name() string {
	return "Get Task"
}

func (g GetTaskAction) Description() string {
	return "Retrieve details for a specific task by ID"
}

func (g GetTaskAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getTaskDocs,
	}
}

func (g GetTaskAction) Icon() *string {
	return nil
}

func (g GetTaskAction) SampleData() sdkcore.JSON {
	return nil
}

func (g GetTaskAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"task_id":    shared.GetTasksInput(),
		"project_id": shared.GetProjectsInput(),
	}
}

func (g GetTaskAction) Auth() *sdk.Auth {
	return &sdk.Auth{
		Inherit: true,
	}
}

func (g GetTaskAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getTaskActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	baseURL := shared.BaseAPI + "/tasks/" + input.TaskID
	params := url.Values{}

	params.Add("opt_fields", "gid,name,completed,due_on,notes,assignee,description,assignee_status,created_at,modified_at,projects,subtasks,parent,tags,workspace")

	reqURL := baseURL
	if len(params) > 0 {
		reqURL = baseURL + "?" + params.Encode()
	}

	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}

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
		return nil, fmt.Errorf("request failed with status code: %d", res.StatusCode)
	}

	var response interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func NewGetTaskAction() sdk.Action {
	return &GetTaskAction{}
}
