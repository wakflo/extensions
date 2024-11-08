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

	"github.com/wakflo/go-sdk/autoform"

	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type addCommentIssueOperationProps struct {
	ProjectID   string `json:"projectId"`
	IssueID     string `json:"issueId"`
	CommentBody string `json:"commentText"`
}

type AddCommentToIssueOperation struct {
	options *sdk.OperationInfo
}

func NewAddCommentToIssueOperation() *AddCommentToIssueOperation {
	return &AddCommentToIssueOperation{
		options: &sdk.OperationInfo{
			Name:        "Add Comment to Issue",
			Description: "Adds comment to an issue.",
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"projectId": getProjectsInput(),
				"issueId":   getIssuesInput(),
				"commentText": autoform.NewLongTextField().
					SetDisplayName("Comment Body").
					SetRequired(true).
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

func (c *AddCommentToIssueOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	input, err := sdk.InputToTypeSafely[addCommentIssueOperationProps](ctx)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("/rest/api/3/issue/%s/comment", input.IssueID)

	instanceURL := ctx.Auth.Extra["instance-url"] + endpoint

	payload := map[string]interface{}{
		"body": map[string]interface{}{
			"type":    "doc",
			"version": 1,
			"content": []map[string]interface{}{
				{
					"type": "paragraph",
					"content": []map[string]interface{}{
						{
							"text": input.CommentBody,
							"type": "text",
						},
					},
				},
			},
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := jiraRequest(ctx.Auth.Extra["email"], ctx.Auth.Extra["api-token"], instanceURL, http.MethodPost, "Comment Added!", data)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *AddCommentToIssueOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *AddCommentToIssueOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
