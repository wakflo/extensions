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
	LabelID     string `json:"label-id"`
	StateID     string `json:"state-id"`
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
				"priority":    getPriorityInput("Priority", "select issue priority"),
				"assignee-id": getAssigneesInput("Assignee", "select an assignee for the issue"),
				"label-id":    getTeamLabelsInput("Labels", "select an issue label"),
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

	// Create a map to store fields conditionally
	fields := make(map[string]string)
	fields["title"] = fmt.Sprintf(`"%s"`, input.Title)
	fields["teamId"] = fmt.Sprintf(`"%s"`, input.TeamID)

	if input.Description != "" {
		fields["description"] = fmt.Sprintf(`"%s"`, input.Description)
	}
	if input.AssigneeID != "" {
		fields["assigneeId"] = fmt.Sprintf(`"%s"`, input.AssigneeID)
	}
	if input.LabelID != "" {
		fields["labelIds"] = fmt.Sprintf(`"%s"`, input.LabelID)
	}

	if input.StateID != "" {
		fields["stateId"] = fmt.Sprintf(`"%s"`, input.StateID)
	}

	if input.Priority != "" {
		priority, err := strconv.Atoi(input.Priority)
		if err != nil {
			return nil, err
		}
		fields["priority"] = strconv.Itoa(priority)
	}

	fieldStrings := make([]string, 0, len(fields))
	for key, value := range fields {
		fieldStrings = append(fieldStrings, fmt.Sprintf("%s: %s", key, value))
	}

	query := fmt.Sprintf(`
		mutation IssueCreate {
			issueCreate(input: {
				%s
			}) {
				success
				issue {
					id
					title
					description
					priorityLabel
				}
			}
		}`, strings.Join(fieldStrings, "\n"))

	response, err := MakeGraphQLRequest(apiKEY, query)
	if err != nil {
		log.Fatalf("Error making GraphQL request: %v", err)
	}

	issue, ok := response["data"].(map[string]interface{})["issueCreate"].(map[string]interface{})["issue"].(map[string]interface{})
	if !ok {
		return nil, errors.New("failed to extract issue from response")
	}

	return issue, nil
}

func (c *CreateIssueOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CreateIssueOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
