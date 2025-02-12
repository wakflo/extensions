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
	"encoding/json"
	"fmt"
	"net/http"

	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
	sdk2 "github.com/wakflo/go-sdk/sdk"
)

type assignIssueOperationProps struct {
	ProjectID string `json:"projectId"`
	IssueID   string `json:"issueId"`
	Assignee  string `json:"assignee"`
}

type AssignIssueOperation struct {
	options *sdk.OperationInfo
}

func NewAssignIssueOperation() *AssignIssueOperation {
	return &AssignIssueOperation{
		options: &sdk.OperationInfo{
			Name:        "Assign Issue",
			Description: "Assigns an issue to a user.",
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"projectId": getProjectsInput(),
				"issueId":   getIssuesInput(),
				"assignee":  getUsersInput(),
			},
			SampleOutput: map[string]interface{}{},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
			RequireAuth: true,
		},
	}
}

func (c *AssignIssueOperation) Run(ctx *sdk.RunContext) (sdk2.JSON, error) {
	input := sdk.InputToType[assignIssueOperationProps](ctx)
	endpoint := fmt.Sprintf("/rest/api/3/issue/%s/assignee", input.IssueID)

	instanceURL := ctx.Auth.Extra["instance-url"] + endpoint

	payload := map[string]interface{}{
		"accountId": input.Assignee,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := jiraRequest(ctx.Auth.Extra["email"], ctx.Auth.Extra["api-token"], instanceURL, http.MethodPut, "Issue assigned", data)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *AssignIssueOperation) Test(ctx *sdk.RunContext) (sdk2.JSON, error) {
	return c.Run(ctx)
}

func (c *AssignIssueOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
