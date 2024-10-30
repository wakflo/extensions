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
	"net/http"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type editIssueOperationProps struct {
	IssueTypeID string `json:"IssueTypeId,omitempty"`
	ProjectID   string `json:"projectId"`
	IssueID     string `json:"issueId"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Assignee    string `json:"assignee"`
	Priority    string `json:"priority"`
	ParentKey   string `json:"parentKey"`
}

type EditIssueOperation struct {
	options *sdk.OperationInfo
}

func NewEditIssueOperation() *EditIssueOperation {
	return &EditIssueOperation{
		options: &sdk.OperationInfo{
			Name:        "Update Issue",
			Description: "Update an issue in a project",
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"projectId":   getProjectsInput(),
				"issueId":     getIssuesInput(),
				"issueTypeId": getIssueTypesInput(),
				"summary": autoform.NewShortTextField().
					SetDisplayName("Summary").
					SetRequired(false).
					Build(),
				"description": autoform.NewLongTextField().
					SetDisplayName("Description").
					SetRequired(false).
					Build(),
				"assignee": getUsersInput(),
				"priority": autoform.NewSelectField().
					SetDisplayName("Priority").
					SetDescription("Priority level of issue").
					SetOptions(PriorityLevels).
					SetRequired(false).
					Build(),
				"parentKey": autoform.NewShortTextField().
					SetDisplayName("Parent Key").
					SetDescription("If this issue is a subtask, insert the parent issue key").
					SetRequired(false).
					Build(),
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

func (c *EditIssueOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	input := sdk.InputToType[editIssueOperationProps](ctx)
	instanceURL := ctx.Auth.Extra["instance-url"] + "/rest/api/3/issue/" + input.IssueID

	payload := map[string]interface{}{
		"fields": map[string]interface{}{},
	}

	if input.Description != "" {
		payload["fields"].(map[string]interface{})["description"] = map[string]interface{}{
			"type":    "doc",
			"version": 1,
			"content": []map[string]interface{}{
				{
					"type": "paragraph",
					"content": []map[string]interface{}{
						{
							"type": "text",
							"text": input.Description,
						},
					},
				},
			},
		}
	}

	if input.Assignee != "" {
		payload["fields"].(map[string]interface{})["assignee"] = map[string]interface{}{
			"id": input.Assignee,
		}
	}

	if input.IssueTypeID != "" {
		payload["fields"].(map[string]interface{})["issueTypeId"] = map[string]interface{}{
			"id": input.Assignee,
		}
	}

	if input.ProjectID != "" {
		payload["fields"].(map[string]interface{})["project"] = map[string]interface{}{
			"id": input.ProjectID,
		}
	}

	if input.Summary != "" {
		payload["fields"] = map[string]interface{}{
			"summary": input.Summary,
		}
	}

	if input.Priority != "" {
		payload["fields"].(map[string]interface{})["priority"] = map[string]interface{}{
			"id": input.Priority,
		}
	}

	if input.ParentKey != "" {
		payload["fields"].(map[string]interface{})["parent"] = map[string]interface{}{
			"key": input.ParentKey,
		}
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := jiraRequest(ctx.Auth.Extra["email"], ctx.Auth.Extra["api-token"], instanceURL, http.MethodPut, "Issue Updated", data)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *EditIssueOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *EditIssueOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
