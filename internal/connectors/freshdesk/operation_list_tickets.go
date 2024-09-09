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

package freshdesk

import (
	"errors"
	"fmt"
	"log"

	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type ListTicketsOperation struct {
	options *sdk.OperationInfo
}

func NewListTicketsOperation() *ListTicketsOperation {
	return &ListTicketsOperation{
		options: &sdk.OperationInfo{
			Name:        "Get a list of tickets",
			Description: "Retrieves a list of tickets",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input:       map[string]*sdkcore.AutoFormSchema{},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *ListTicketsOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing freshdesk api key")
	}

	domain := ctx.Auth.Extra["domain"]
	freshdeskDomain := "https://" + domain + ".freshdesk.com"
	fmt.Println(freshdeskDomain)

	response, err := GetTickets(freshdeskDomain, ctx.Auth.Extra["api-key"])
	if err != nil {
		log.Fatalf("error fetching data: %v", err)
	}

	return response, nil
}

func (c *ListTicketsOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *ListTicketsOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
