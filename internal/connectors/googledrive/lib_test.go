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

package googledrive

import (
	"testing"
	"time"

	"golang.org/x/oauth2"

	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

var (
	testToken     = oauth2.StaticTokenSource(&oauth2.Token{})
	lastFourHours = time.Now().Add(-4 * time.Hour)
)

func TestNewConnector(t *testing.T) {
	testCases := []struct {
		name          string
		operationName string
		wantErr       bool
		isTrigger     bool
		ctx           *sdk.RunContext
	}{
		{
			name:          "Success",
			operationName: "create-new-file",
			wantErr:       false,
			isTrigger:     true,
			ctx: &sdk.RunContext{
				Auth: &sdkcore.AuthContext{
					AccessToken: "some-token",
					Token:       nil,
					TokenType:   "",
					Username:    "",
					Password:    "",
					Secret:      "",
					TokenSource: &testToken,
				},
				Metadata: &sdkcore.FlowRunMetadata{
					LastRun: &lastFourHours,
				},
			},
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

			if connector.DisplayName != "Google Drive" {
				t.Errorf("NewConnector() Name = %s, want %s", connector.DisplayName, "Google Drive")
			}

			if connector.Description != "some google drive connector" {
				t.Errorf("NewConnector() Description = %s, want %s", connector.Description, "some google drive connector")
			}

			if connector.Logo != "logos:google-drive" {
				t.Errorf("NewConnector() Logo = %s, want %s", connector.Logo, "logos:google-drive")
			}

			if connector.Version != "0.0.1" {
				t.Errorf("NewConnector() Version = %s, want %s", connector.Version, "0.0.1")
			}

			if connector.Group != sdk.ConnectorGroupApps {
				t.Errorf("NewConnector() Group = %v, want %v", connector.Group, sdk.ConnectorGroupApps)
			}

			if len(connector.Authors) != 1 || connector.Authors[0] != "Wakflo <integrations@wakflo.com>" {
				t.Errorf("NewConnector() Authors = %v, want %v", connector.Authors, []string{"Wakflo <integrations@wakflo.com>"})
			}

			if len(spider.Triggers()) != 2 {
				t.Errorf("NewConnector() Triggers() count = %d, want %d", len(spider.Triggers()), 2)
			}

			if len(spider.Operations()) != 6 {
				t.Errorf("NewConnector() ActionsMap() count = %d, want %d", len(spider.Operations()), 6)
			}

			// if testCase.isTrigger {
			//	_, err = spider.RunTrigger(testCase.operationName, testCase.ctx)
			//	if err != nil {
			//		t.Errorf("NewConnector() RunTrigger() with name %v threw an error = %v", testCase.operationName, err)
			//	}
			// } else {
			//	_, err = spider.RunOperation(testCase.operationName, testCase.ctx)
			//	if err != nil {
			//		t.Errorf("NewConnector() RunOperation() with name %v threw an error = %v", testCase.operationName, err)
			//	}
			//}

			/*
				/*_, err = spider.RunOperation(testCase.operationName, ctx)
				if err != nil {
					t.Errorf("NewConnector() RunOperation() with name %v threw an error = %v", testCase.operationName, err)
				}

				if trigger!= "some data" {
					t.Errorf("NewConnector() RunOperation() response = %v, want %v", trigger["data"], "some data")
				}*/
		})
	}
}
