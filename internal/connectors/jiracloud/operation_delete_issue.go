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
	"net/http"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
	sdk2 "github.com/wakflo/go-sdk/sdk"
)

type deleteIssueOperationProps struct {
	IssueID string `json:"IssueId,omitempty"`
}

type DeleteIssueOperation struct {
	options *sdk.OperationInfo
}

func NewDeleteIssueOperation() *DeleteIssueOperation {
	return &DeleteIssueOperation{
		options: &sdk.OperationInfo{
			Name:        "Delete an Issue",
			Description: "Delete an issue in a project",
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"issueId": autoform.NewShortTextField().
					SetDisplayName("Issue ID or Key").
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

func (c *DeleteIssueOperation) Run(ctx *sdk.RunContext) (sdk2.JSON, error) {
	input := sdk.InputToType[deleteIssueOperationProps](ctx)
	instanceURL := ctx.Auth.Extra["instance-url"] + "/rest/api/3/issue/" + input.IssueID

	resp, err := jiraRequest(ctx.Auth.Extra["email"], ctx.Auth.Extra["api-token"], instanceURL, http.MethodDelete, "Issue deleted", nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *DeleteIssueOperation) Test(ctx *sdk.RunContext) (sdk2.JSON, error) {
	return c.Run(ctx)
}

func (c *DeleteIssueOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
