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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/asana/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/integration"
)

type createTaskActionProps struct {
	Name      string  `json:"name"`
	Project   *string `json:"projects"`
	Workspace *string `json:"workspace"`
}

type CreateTaskAction struct{}

func (c CreateTaskAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (c CreateTaskAction) Name() string {
	return "Create Task"
}

func (c CreateTaskAction) Description() string {
	return "Create a new task"
}

func (c CreateTaskAction) Documentation() *integration.OperationDocumentation {
	return &integration.OperationDocumentation{
		Documentation: &createTaskDocs,
	}
}

func (c CreateTaskAction) Icon() *string {
	return nil
}

func (c CreateTaskAction) SampleData() (sdkcore.JSON, error) {
	return nil, nil
}

func (c CreateTaskAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"name": autoform.NewShortTextField().
			SetDisplayName("Task Name").
			SetDescription("The name task's Name.").
			SetRequired(true).
			Build(),
		"workspace": shared.GetWorkspacesInput(),
		"projects":  shared.GetProjectsInput(),
	}
}

func (c CreateTaskAction) Auth() *integration.Auth {
	return &integration.Auth{
		Inherit: true,
	}
}

func (c CreateTaskAction) Perform(ctx integration.PerformContext) (sdkcore.JSON, error) {
	input, err := integration.InputToTypeSafely[createTaskActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	var projects []string
	if input.Project != nil {
		projects = append(projects, *input.Project)
	}

	taskData := map[string]interface{}{
		"data": map[string]interface{}{
			"name":      input.Name,
			"workspace": input.Workspace,
			"projects":  projects,
		},
	}

	taskJSON, err := json.Marshal(taskData)
	if err != nil {
		return nil, err
	}

	reqURL := "https://app.asana.com/api/1.0/tasks"

	req, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewBuffer(taskJSON))
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
		return nil, nil
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil, nil
	}

	fmt.Println(string(body))

	return map[string]interface{}{
		"Result": "Task created Successfully",
	}, nil
}

func NewCreateTaskAction() integration.Action {
	return &CreateTaskAction{}
}
