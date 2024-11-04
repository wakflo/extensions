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

package hubspot

import (
	"net/http"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type listTicketsOperationProps struct {
	Subject string `json:"subject"`
}

type ListTicketsOperation struct {
	options *sdk.OperationInfo
}

func NewListTicketsOperation() *ListTicketsOperation {
	return &ListTicketsOperation{
		options: &sdk.OperationInfo{
			Name:        "List Tickets",
			Description: "List tickets",
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"subject": autoform.NewShortTextField().
					SetDisplayName("").
					SetDescription("").
					SetRequired(false).Build(),
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

func (c *ListTicketsOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	_ = sdk.InputToType[listTicketsOperationProps](ctx)

	reqURL := "https://api.hubapi.com/crm/v3/objects/tickets"

	resp, err := hubspotClient(reqURL, ctx.Auth.AccessToken, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *ListTicketsOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *ListTicketsOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
