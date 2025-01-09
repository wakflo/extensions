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
	"strings"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type updateIssueOperationProps struct {
	TeamID      string `json:"team-id"`
	IssueID     string `json:"issue-id"`
	AssigneeID  string `json:"assignee-id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	StateID     string `json:"state-id"`
	LabelID     string `json:"label-id"`
}

type UpdateIssueOperation struct {
	options *sdk.OperationInfo
}

func NewUpdateIssueOperation() *UpdateIssueOperation {
	return &UpdateIssueOperation{
		options: &sdk.OperationInfo{
			Name:        "Update Issue",
			Description: "Update an issue in Linear Workspace",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"title": autoform.NewShortTextField().
					SetDisplayName("Issue Name").
					SetDescription("The issue name").
					SetRequired(false).
					Build(),
				"description": autoform.NewLongTextField().
					SetDisplayName("Description").
					SetDescription("Issue description").
					Build(),
				"team-id":     getTeamsInput(),
				"issue-id":    getIssuesInput("Select issue", "select an issue to update"),
				"assignee-id": getAssigneesInput("Assignee", "Select assignee"),
				"priority":    getPriorityInput("Priority", "Select priority"),
				"label-id":    getTeamLabelsInput("Label", "Select label"),
				"state-id":    getIssueStatesInput("Issue State", "select issue state", false),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *UpdateIssueOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing linear api key")
	}
	apiKEY := ctx.Auth.Extra["api-key"]

	if !strings.HasPrefix(apiKEY, "lin_api_") {
		return nil, errors.New("invalid Linear API key: must start with 'lin_api_'")
	}

	input := sdk.InputToType[updateIssueOperationProps](ctx)

	mutation := fmt.Sprintf(`mutation IssueUpdate {
      issueUpdate(
        id: "%s",
        input: {`, input.IssueID)

	if input.Title != "" {
		mutation += fmt.Sprintf(`title: "%s",`, input.Title)
	}
	if input.StateID != "" {
		mutation += fmt.Sprintf(`stateId: "%s",`, input.StateID)
	}
	if input.Description != "" {
		mutation += fmt.Sprintf(`description: "%s",`, input.Description)
	}
	if input.LabelID != "" {
		mutation += fmt.Sprintf(`labelIds: "%s",`, input.LabelID)
	}
	if input.Priority != "" {
		mutation += fmt.Sprintf(`priority: %s,`, input.Priority)
	}
	if input.AssigneeID != "" {
		mutation += fmt.Sprintf(`assigneeId: "%s",`, input.AssigneeID)
	}

	mutation = strings.TrimSuffix(mutation, ",")
	mutation += `}) {
        success
        issue {
          id
          title
          priorityLabel
          priority
        }
      }
    }`

	response, err := MakeGraphQLRequest(apiKEY, mutation)
	if err != nil {
		log.Fatalf("Error making GraphQL request: %v", err)
	}

	issue, ok := response["data"].(map[string]interface{})["issueUpdate"]
	if !ok {
		return nil, errors.New("failed to extract issue from response")
	}

	return issue, nil
}

func (c *UpdateIssueOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *UpdateIssueOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
