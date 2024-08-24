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

type getIssuesByConditionOperationProps struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	AssigneeID  string `json:"assignee-id"`
	LabelID     string `json:"labelID"`
	StateID     string `json:"state-id"`
}

type GetIssuesByConditionOperation struct {
	options *sdk.OperationInfo
}

func NewGetIssuesByConditionOperation() *GetIssuesByConditionOperation {
	return &GetIssuesByConditionOperation{
		options: &sdk.OperationInfo{
			Name:        "Get issues by condition",
			Description: "Get issues that meet specific conditions",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"title": autoform.NewShortTextField().
					SetDisplayName("Issue Name").
					SetDescription("Filter by the issue name").
					Build(),
				"description": autoform.NewLongTextField().
					SetDisplayName("Description").
					SetDescription("Filter Issue by description").
					Build(),
				"label-id":    getLabelsInput("Filter by label", "Filter issue by label"),
				"assignee-id": getAssigneesInput("Filter by assignees", "Filter issue by assignees"),
				"state-id":    getIssueStatesInput("Filter by state", "filter issue by state", false),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *GetIssuesByConditionOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing linear api key")
	}
	apiKEY := ctx.Auth.Extra["api-key"]

	if !strings.HasPrefix(apiKEY, "lin_api_") {
		return nil, errors.New("invalid Linear API key: must start with 'lin_api_'")
	}

	input, err := sdk.InputToTypeSafely[getIssuesByConditionOperationProps](ctx)
	if err != nil {
		return nil, err
	}

	var filters []string

	if input.AssigneeID != "" {
		filters = append(filters, fmt.Sprintf(`assignee: {
      id: {
        eq:"%s"
      }
    }`, input.AssigneeID))
	}
	if input.Description != "" {
		filters = append(filters, fmt.Sprintf(`description:{
      containsIgnoreCase: "%s"
    }`, input.Description))
	}
	if input.Title != "" {
		filters = append(filters, fmt.Sprintf(`title:{
      containsIgnoreCase: "%s"
    }`, input.Title))
	}
	if input.LabelID != "" {
		filters = append(filters, fmt.Sprintf(`labelIds: ["%s"]`, input.LabelID))
	}
	if input.StateID != "" {
		filters = append(filters, fmt.Sprintf(`  state: {
        id: {
			eq: "%s"
		}


      }`, input.StateID))
	}

	filterString := strings.Join(filters, ", ")

	query := fmt.Sprintf(`
	{
		issues(filter: {%s}) {
			nodes {
				id
				title
				description
				state {
					name
				}
				priority
				createdAt
				assignee {
					name
				}
			}
		}
	}`, filterString)

	response, err := MakeGraphQLRequest(apiKEY, query)
	if err != nil {
		log.Fatalf("Error making GraphQL request: %v", err)
	}

	return map[string]interface{}{
		"Result": response,
	}, nil
}

func (c *GetIssuesByConditionOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *GetIssuesByConditionOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
