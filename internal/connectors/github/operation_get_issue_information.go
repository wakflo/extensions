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

type getIssueInfoOperationProps struct {
	Repository  string `json:"repository"`
	IssueNumber int    `json:"issue_number"`
}

type GetIssueInfoOperation struct {
	options *sdk.OperationInfo
}

func NewGetIssueInfoOperation() *GetIssueInfoOperation {
	return &GetIssueInfoOperation{
		options: &sdk.OperationInfo{
			Name:        "Get Issue information",
			Description: "Get information about a specific issue",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"repository": getRepositoryInput(),
				"issue_number": autoform.NewNumberField().
					SetDisplayName("Issue Number").
					SetDescription("The issue number").
					SetRequired(true).
					Build(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *GetIssueInfoOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	input, err := sdk.InputToTypeSafely[getIssueInfoOperationProps](ctx)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`
		query {
            node(id: "%s") {
				... on Repository {
				    issue(number:%d ){
					    title
						body
						createdAt
						updatedAt
						number
					author {
					  login
					}
					assignees(first:10){
					    nodes {
						     name
							 login
					   }
					}
				}
			}
		}
}`, input.Repository, input.IssueNumber)

	response, err := githubGQL(ctx.Auth.AccessToken, query)
	if err != nil {
		return nil, errors.New("error making graphQL request")
	}

	issue, ok := response["data"].(map[string]interface{})["node"].(map[string]interface{})["issue"].(map[string]interface{})
	if !ok {
		return nil, errors.New("failed to extract issue from response")
	}

	return issue, nil
}

func (c *CreateIssueCommentOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CreateIssueCommentOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
