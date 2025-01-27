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
	"time"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type TaskCreated struct {
	options *sdk.TriggerInfo
}

func NewTaskCreated() *TicketCreated {
	return &TicketCreated{
		options: &sdk.TriggerInfo{
			Name:        "New Task Added",
			Description: "triggers workflow when a new task has been added",
			RequireAuth: true,
			Auth:        sharedAuth,
			Type:        sdkcore.TriggerTypePolling,
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

func (t *TaskCreated) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	reqURL := "/crm/v3/objects/tasks?limit=50&archived=false&properties=createdAt,updatedAt"

	if ctx.Metadata.LastRun != nil {
		createdAfter := ctx.Metadata.LastRun.UTC().Format(time.RFC3339)
		reqURL += "&createdAfter=" + createdAfter
	}

	resp, err := hubspotClient(reqURL, ctx.Auth.AccessToken, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (t *TaskCreated) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return t.Run(ctx)
}

func (t *TaskCreated) OnEnabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *TaskCreated) OnDisabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *TaskCreated) GetInfo() *sdk.TriggerInfo {
	return t.options
}
