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

package mailchimp

import (
	"errors"
	"log"
	"time"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type newSubscriberProps struct {
	ListID string `json:"list_id"`
}

type TriggerNewSubscriber struct {
	options *sdk.TriggerInfo
}

func NewTriggerNewSubscriber() *TriggerNewSubscriber {
	return &TriggerNewSubscriber{
		options: &sdk.TriggerInfo{
			Name:        "New Subscriber",
			Description: "Runs when an Audience subscriber is added.",
			RequireAuth: true,
			Auth:        sharedAuth,
			Type:        sdkcore.TriggerTypePolling,
			Input: map[string]*sdkcore.AutoFormSchema{
				"list_id": autoform.NewShortTextField().
					SetDisplayName("List ID").
					SetDescription("").
					SetRequired(true).
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

func (t *TriggerNewSubscriber) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.AccessToken == "" {
		return nil, errors.New("missing mailchimp auth token")
	}
	accessToken := ctx.Auth.AccessToken

	input := sdk.InputToType[newSubscriberProps](ctx)
	dc, err := getMailChimpServerPrefix(accessToken)
	if err != nil {
		log.Fatalf("Error getting MailChimp server prefix: %v", err)
	}

	var fromDate string
	lastRunTime := ctx.Metadata.LastRun

	if lastRunTime != nil {
		fromDate = lastRunTime.UTC().Format(time.RFC3339)
	}

	var newSubscribers interface{}

	newSubscribers, err = listRecentSubscribers(accessToken, dc, input.ListID, fromDate)
	if err != nil {
		return nil, err
	}
	return newSubscribers, nil
}

func (t *TriggerNewSubscriber) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return t.Run(ctx)
}

func (t *TriggerNewSubscriber) OnEnabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *TriggerNewSubscriber) OnDisabled(ctx *sdk.RunContext) error {
	return nil
}

func (t *TriggerNewSubscriber) GetInfo() *sdk.TriggerInfo {
	return t.options
}
