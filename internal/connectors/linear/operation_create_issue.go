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

package linear

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type createIssueOperationProps struct {
	TeamID      string `json:"team-id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	AssigneeID  string `json:"assignee-id"`
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
				"title": autoform.NewShortTextField().
					SetDisplayName("Issue Name").
					SetDescription("The issue name").
					SetRequired(true).
					Build(),
				"description": autoform.NewLongTextField().
					SetDisplayName("Description").
					SetDescription("Issue description").
					Build(),
				"team-id":     getTeamsInput(),
				"priority":    getPriorityInput("Select Priority", "select issue priority"),
				"assignee-id": getAssigneesInput("Select Assignee", "select an assignee for the issue"),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *CreateIssueOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing linear api key")
	}

	apiKEY := ctx.Auth.Extra["api-key"]

	if !strings.HasPrefix(apiKEY, "lin_api_") {
		return nil, errors.New("invalid Linear API key: must start with 'lin_api_'")
	}

	input, err := sdk.InputToTypeSafely[createIssueOperationProps](ctx)
	if err != nil {
		return nil, err
	}

	priority, err := strconv.Atoi(input.Priority)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`
    mutation IssueCreate{
      issueCreate(input: {
        title: "%s"
        description: "%s"
        teamId: "%s"
        assigneeId: "%s"
		priority: %d
      }) {
        success
        issue {
          id
          title
	      priorityLabel
          priority
        }
      }
    }`, input.Title, input.Description, input.TeamID, input.AssigneeID, priority)

	response, err := MakeGraphQLRequest(apiKEY, query)
	if err != nil {
		log.Fatalf("Error making GraphQL request: %v", err)
	}

	return map[string]interface{}{
		"Result": response,
	}, nil
}

func (c *CreateIssueOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CreateIssueOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
