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

package zendesk

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type TicketCreated struct {
	options *sdk.TriggerInfo
}

func NewTicketCreated() *TicketCreated {
	return &TicketCreated{
		options: &sdk.TriggerInfo{
			Name:        "New Ticket Added",
			Description: "triggers workflow when a new ticket is added",
			RequireAuth: true,
			Auth:        sharedAuth,
			Type:        sdkcore.TriggerTypeCron,
			Input: map[string]*sdkcore.AutoFormSchema{
				"id": autoform.NewShortTextField().
					SetDisplayName("").
					SetDescription("").
					Build(),
			},
			Settings: &sdkcore.TriggerSettings{},
			ErrorSettings: &sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (t *TicketCreated) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["api-token"] == "" || ctx.Auth.Extra["email"] == "" || ctx.Auth.Extra["subdomain"] == "" {
		return nil, errors.New("missing zendesk api credentials")
	}

	email := ctx.Auth.Extra["email"]
	apiToken := ctx.Auth.Extra["api-token"]
	subdomain := "https://" + ctx.Auth.Extra["subdomain"] + ".zendesk.com/api/v2"

	fullURL := subdomain + "/search.json?query=type:ticket"

	if ctx.Metadata.LastRun != nil {
		createdAfter := ctx.Metadata.LastRun.UTC().Format("2006-01-02T15:04:05Z")
		fullURL = fmt.Sprintf("%s+created>=%s", fullURL, createdAfter)
	}

	response, err := zendeskRequest(http.MethodGet, fullURL, email, apiToken, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}

	return response, nil
}

func (t *TicketCreated) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return t.Run(ctx)
}

func (t *TicketCreated) OnEnabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *TicketCreated) OnDisabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *TicketCreated) GetInfo() *sdk.TriggerInfo {
	return t.options
}
