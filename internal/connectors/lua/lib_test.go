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

package lua

import (
	"testing"

	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

func TestNewConnector(t *testing.T) {
	opsId := "run-lua-script"
	testCases := []struct {
		name          string
		operationName string
		wantErr       bool
	}{
		{
			name:          "Success",
			operationName: opsId,
			wantErr:       false,
		},
	}

	ctx := &sdk.RunContext{
		Workflow: nil,
		Step: &sdkcore.ConnectorStep{
			Label:     "label",
			Icon:      "logo",
			Name:      "step-1",
			IsTrigger: false,
			Path:      nil,
			NodeIndex: 0,
			Type:      "normal",
			Data: sdkcore.ConnectorStepData{
				OperationID:      &opsId,
				AuthConnectionID: nil,
				Properties: sdkcore.ConnectorProperties{
					Input:  map[string]interface{}{},
					Output: map[string]interface{}{},
				},
			},
			Children:      nil,
			Reference:     nil,
			Metadata:      sdkcore.ConnectorStepMetadata{},
			ParentID:      nil,
			ErrorSettings: sdkcore.StepErrorSettings{},
			Valid:         false,
		},
		Auth: &sdkcore.AuthContext{
			AccessToken: "",
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

			if connector.DisplayName != "Lua" {
				t.Errorf("NewConnector() Name = %s, want %s", connector.DisplayName, "Lua")
			}

			if connector.Description != "Lua connector for running lua codes" {
				t.Errorf("NewConnector() Description = %s, want %s", connector.Description, "Lua connector for running lua codes")
			}

			if connector.Logo != "logos:lua" {
				t.Errorf("NewConnector() Logo = %s, want %s", connector.Logo, "logos:google-drive")
			}

			if connector.Version != "0.0.1" {
				t.Errorf("NewConnector() Version = %s, want %s", connector.Version, "0.0.1")
			}

			if connector.Category != sdk.Tools {
				t.Errorf("NewConnector() Category = %v, want %v", connector.Category, sdk.Tools)
			}

			if len(connector.Authors) != 1 || connector.Authors[0] != "Wakflo <integrations@wakflo.com>" {
				t.Errorf("NewConnector() Authors = %v, want %v", connector.Authors, []string{"Wakflo <integrations@wakflo.com>"})
			}

			if len(spider.Operations()) != 1 {
				t.Errorf("NewConnector() Operations() count = %d, want %d", len(spider.Operations()), 1)
			}

			_, err = spider.RunOperation(testCase.operationName, ctx)
			if err != nil {
				t.Errorf("NewConnector() RunOperation() with name %v threw an error = %v", testCase.operationName, err)
			}

			//if trigger!= "some data" {
			//	t.Errorf("NewConnector() RunOperation() response = %v, want %v", trigger["data"], "some data")
			//}
		})
	}
}
