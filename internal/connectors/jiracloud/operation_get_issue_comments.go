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
	"fmt"
	"net/http"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type getIssueCommentsOperationProps struct {
	IssueID string `json:"IssueId,omitempty"`
	Limit   int    `json:"limit,omitempty"`
	OrderBy string `json:"orderBy,omitempty"`
}

type GetIssueCommentsOperation struct {
	options *sdk.OperationInfo
}

func NewGetIssueCommentsOperation() *GetIssueCommentsOperation {
	return &GetIssueCommentsOperation{
		options: &sdk.OperationInfo{
			Name:        "List Issue Comments",
			Description: "Get a list of issue comments",
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"projectId": getProjectsInput(),
				"issueId":   getIssuesInput(),
				"orderBy": autoform.NewSelectField().
					SetDisplayName("Order by").
					SetOptions(OrderBy).
					SetDefaultValue("-created").
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

func (c *GetIssueCommentsOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	input := sdk.InputToType[getIssueCommentsOperationProps](ctx)
	endpoint := fmt.Sprintf("/rest/api/3/issue/%s/comment?orderBy=%s", input.IssueID, input.OrderBy)

	instanceURL := ctx.Auth.Extra["instance-url"] + endpoint

	resp, err := jiraRequest(ctx.Auth.Extra["email"], ctx.Auth.Extra["api-token"], instanceURL, http.MethodGet, "", nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *GetIssueCommentsOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *GetIssueCommentsOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
