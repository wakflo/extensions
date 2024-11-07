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

package github

import (
	"errors"
	"fmt"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type createIssueCommentOperationProps struct {
	Repository  string `json:"repository"`
	Body        string `json:"body"`
	IssueNumber string `json:"issue_number"`
}

type CreateIssueCommentOperation struct {
	options *sdk.OperationInfo
}

func NewCreateIssueCommentOperation() *CreateIssueCommentOperation {
	return &CreateIssueCommentOperation{
		options: &sdk.OperationInfo{
			Name:        "Create comment on a issue",
			Description: "Adds a comment to the specified issue (also works with pull requests)",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"repository":   getRepositoryInput(),
				"issue_number": getIssuesInput(),
				"body": autoform.NewLongTextField().
					SetDisplayName("Comment").
					SetDescription("Issue comment").
					Build(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *CreateIssueCommentOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	input, err := sdk.InputToTypeSafely[createIssueCommentOperationProps](ctx)
	if err != nil {
		return nil, err
	}

	mutation := fmt.Sprintf(`
			mutation {
			  addComment(input: { subjectId: "%s", body: "%s" }) {
				commentEdge {
				  node {
					id
					body
					createdAt
				  }
				}
			  }
			}`, input.IssueNumber, input.Body)

	response, err := githubGQL(ctx.Auth.AccessToken, mutation)
	if err != nil {
		return nil, errors.New("error making graphQL request")
	}

	issue, ok := response["data"].(map[string]interface{})["addComment"].(map[string]interface{})["commentEdge"].(map[string]interface{})
	if !ok {
		return nil, errors.New("failed to extract issue from response")
	}

	return issue, nil
}

func (c *GetIssueInfoOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *GetIssueInfoOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
