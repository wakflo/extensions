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
	"fmt"
	"log"
	"strings"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type createIssueOperationProps struct {
	Repository string `json:"repository"`
	Title      string `json:"title"`
	Body       string `json:"body"`
	Labels     string `json:"labels"`
}

type CreateIssueOperation struct {
	options *sdk.OperationInfo
}

func NewCreateIssueOperation() *CreateIssueOperation {
	return &CreateIssueOperation{
		options: &sdk.OperationInfo{
			Name:        "Create Issue",
			Description: "Create an issue",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"repository": getRepositoryInput(),
				"title": autoform.NewShortTextField().
					SetDisplayName("Issue Name").
					SetDescription("The issue name").
					SetRequired(true).
					Build(),
				"body": autoform.NewLongTextField().
					SetDisplayName("Description").
					SetDescription("Issue description").
					Build(),
				"labels": getLabelInput(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *CreateIssueOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	input, err := sdk.InputToTypeSafely[createIssueOperationProps](ctx)
	if err != nil {
		return nil, err
	}

	// Create a map to store fields conditionally
	fields := make(map[string]string)
	fields["title"] = fmt.Sprintf(`"%s"`, input.Title)
	fields["repositoryId"] = fmt.Sprintf(`"%s"`, input.Repository)

	if input.Body != "" {
		fields["body"] = fmt.Sprintf(`"%s"`, input.Body)
	}

	if input.Labels != "" {
		fields["labelIds"] = fmt.Sprintf(`"%s"`, input.Labels)
	}

	fieldStrings := make([]string, 0, len(fields))
	for key, value := range fields {
		fieldStrings = append(fieldStrings, fmt.Sprintf("%s: %s", key, value))
	}

	mutation := fmt.Sprintf(`
		mutation CreateNewIssue {
			createIssue(input: {
				%s
			}) {
				issue {
					id
					title
					url
				}
			}
		}`, strings.Join(fieldStrings, "\n"))

	response, err := githubGQL(ctx.Auth.AccessToken, mutation)
	if err != nil {
		log.Fatalf("Error making GraphQL request: %v", err)
	}

	return response, nil
}

func (c *CreateIssueOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CreateIssueOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
