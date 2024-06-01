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

package slack

import (
	"testing"

	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

func TestNewConnector(t *testing.T) {
	testCases := []struct {
		name          string
		operationName string
		wantErr       bool
		data          map[string]interface{}
	}{
		{
			name:          "Message a public channel",
			operationName: "send-public-channel-message",
			wantErr:       false,
			data: map[string]interface{}{
				"channel": "C06RC8P6133",
				"message": "Hello everyone! Did you know snakes smell with their tongue?",
			},
		},
		{
			name:          "Message a private channel",
			operationName: "send-private-channel-message",
			wantErr:       false,
			data: map[string]interface{}{
				"channel": "C06SZ3SRNN6",
				"message": "Hello everyone! Did you know ketchup was once sold as medicine?",
			},
		},
		{
			name:          "Direct message a user",
			operationName: "send-direct-message",
			wantErr:       false,
			data: map[string]interface{}{
				"user":    "U06QT00GCCE",
				"message": "Hello! Did you know lemons float but limes sink?",
			},
		},
	}

	ctx := &sdk.RunContext{
		Workflow: nil,
		Step: &sdkcore.ConnectorStep{
			Label:     "",
			Icon:      "",
			Name:      "",
			IsTrigger: false,
			Path:      nil,
			Type:      "",
			Data: sdkcore.ConnectorStepData{
				OperationID:      nil,
				AuthConnectionId: nil,
				Properties: sdkcore.ConnectorProperties{
					Input:  map[string]interface{}{},
					Output: nil,
				},
			},
			Children:      nil,
			Reference:     nil,
			Metadata:      sdkcore.ConnectorStepMetadata{},
			ParentId:      nil,
			ErrorSettings: sdkcore.StepErrorSettings{},
			Valid:         false,
		},
		Auth: &sdkcore.AuthContext{
			AccessToken: "xoxb-6849639878738-6862302315249-tRBRRh6W5SabZO0nS0YyPK5d",
			Token:       nil,
			TokenType:   "",
			Username:    "",
			Password:    "",
			Secret:      "",
		},
		State: &sdkcore.StepsState{
			Steps: map[string]*sdkcore.StepState{
				"step-1": {
					ConnectorName: "",
					Version:       "",
					Input:         map[string]interface{}{},
					Output:        map[string]interface{}{},
					Logs:          nil,
					Status:        "",
				},
			},
			CurrentStepID: "step-1",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			instance, err := NewConnector()
			spider := sdk.NewSpiderTest(t, instance)
			connector := spider.GetConfig()

			if (err != nil) != testCase.wantErr {
				t.Errorf("NewConnector() error = %v, wantErr %v", err, testCase.wantErr)
				return
			}

			if connector == nil {
				t.Fatal("NewConnector() returned nil")
			}

			if connector.DisplayName != "Slack" {
				t.Errorf("NewConnector() Name = %s, want %s", connector.DisplayName, "Google Drive")
			}

			if connector.Description != "Send one message to a chosen user, public channel or private channel" {
				t.Errorf("NewConnector() Description = %s, want %s", connector.Description, "Send one message to a chosen user, public channel or private channel")
			}

			if connector.Logo != "logos:slack-icon" {
				t.Errorf("NewConnector() Logo = %s, want %s", connector.Logo, "logos:google-drive")
			}

			if connector.Version != "0.0.1" {
				t.Errorf("NewConnector() Version = %s, want %s", connector.Version, "0.0.1")
			}

			if connector.Category != sdk.Apps {
				t.Errorf("NewConnector() Category = %v, want %v", connector.Category, sdk.Apps)
			}

			if len(connector.Authors) != 1 || connector.Authors[0] != "Wakflo <integrations@wakflo.com>" {
				t.Errorf("NewConnector() Authors = %v, want %v", connector.Authors, []string{"Wakflo <integrations@wakflo.com>"})
			}

			if len(spider.Triggers()) != 0 {
				t.Errorf("NewConnector() Triggers() count = %d, want %d", len(spider.Triggers()), 0)
			}

			if len(spider.Operations()) != 3 {
				t.Errorf("NewConnector() Operations() count = %d, want %d", len(spider.Operations()), 3)
			}

			ctx.Step.Data.Properties.Input = testCase.data

			_, _ = spider.RunOperation(testCase.operationName, ctx)
			//if err != nil {
			//	t.Errorf("NewConnector() RunOperation() with name %v threw an error = %v", testCase.operationName, err)
			//}

			//result := trigger.(map[string]interface{})
			//if result["usage_mode"] != "operation" {
			//	t.Errorf("NewConnector() RunOperation() response = %v, want %v", result["usage_mode"], "operation")
			//}
		})
	}
}
