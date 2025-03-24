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
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/asana/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type updateTaskActionProps struct {
	TaskID    string  `json:"task_id"`
	Name      *string `json:"name"`
	Notes     *string `json:"notes"`
	Completed *bool   `json:"completed"`
	DueOn     *string `json:"due_on"`
}

type UpdateTaskAction struct{}

func (u *UpdateTaskAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func (u UpdateTaskAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (u UpdateTaskAction) Name() string {
	return "Update Task"
}

func (u UpdateTaskAction) Description() string {
	return "Update an existing task in Asana"
}

func (u UpdateTaskAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &updateTaskDocs,
	}
}

func (u UpdateTaskAction) Icon() *string {
	return nil
}

func (u UpdateTaskAction) SampleData() sdkcore.JSON {
	return nil
}

func (u UpdateTaskAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"project_id": shared.GetProjectsInput(),
		"task_id":    shared.GetTasksInput(),
		"name": autoform.NewShortTextField().
			SetDisplayName("Name").
			SetDescription("New name for the task").
			SetRequired(false).
			Build(),
		"notes": autoform.NewLongTextField().
			SetDisplayName("Notes").
			SetDescription("Notes or description for the task").
			SetRequired(false).
			Build(),
		"completed": autoform.NewBooleanField().
			SetDisplayName("Completed").
			SetDescription("Mark the task as completed").
			SetRequired(false).
			Build(),
	}
}

func (u UpdateTaskAction) Auth() *sdk.Auth {
	return &sdk.Auth{
		Inherit: true,
	}
}

func (u UpdateTaskAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[updateTaskActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	taskData := map[string]interface{}{
		"data": map[string]interface{}{},
	}

	data := taskData["data"].(map[string]interface{})

	if input.Name != nil {
		data["name"] = *input.Name
	}

	if input.Notes != nil {
		data["notes"] = *input.Notes
	}

	if input.Completed != nil {
		data["completed"] = *input.Completed
	}

	if len(data) == 0 {
		return nil, errors.New("at least one field must be provided to update the task")
	}

	taskJSON, err := json.Marshal(taskData)
	if err != nil {
		return nil, err
	}

	reqURL := shared.BaseAPI + "/tasks/" + input.TaskID

	req, err := http.NewRequest(http.MethodPut, reqURL, bytes.NewBuffer(taskJSON))
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
		return nil, fmt.Errorf("request failed with status code: %d", res.StatusCode)
	}

	var response interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func NewUpdateTaskAction() sdk.Action {
	return &UpdateTaskAction{}
}
