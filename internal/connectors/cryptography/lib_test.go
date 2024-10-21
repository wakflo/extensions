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

package cryptography

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
			name:          "Test hashing a string with MD5",
			operationName: "hash-text",
			wantErr:       false,
			data: map[string]interface{}{
				"algorithm":       "MD5",
				"text":            "Banana is yellow",
				"expected_output": "d588a672b7cd11af0fbc97d392caa1e2",
			},
		},
		{
			name:          "Test hashing a string with SHA1",
			operationName: "hash-text",
			wantErr:       false,
			data: map[string]interface{}{
				"algorithm":       "SHA1",
				"text":            "Banana is yellow",
				"expected_output": "7baddc14adf5e9454d84e1d4f140c473b61e8fea",
			},
		},
		{
			name:          "Test hashing a string with SHA256",
			operationName: "hash-text",
			wantErr:       false,
			data: map[string]interface{}{
				"algorithm":       "SHA256",
				"text":            "Banana is yellow",
				"expected_output": "965eea67cf759c3a3c851b3046303d5ad2a47c831c0ca157e1e4250b18b5e087",
			},
		},
		{
			name:          "Test hashing a string with SHA512",
			operationName: "hash-text",
			wantErr:       false,
			data: map[string]interface{}{
				"algorithm":       "SHA512",
				"text":            "Banana is yellow",
				"expected_output": "3cd408cc8eb4cb8118a8c4e6a1235039b3840d54f329c86deb4aa1ca2c4fd63aef1b9a77739c6890246ba513aadbcd52448099626c5d30103dc97843b60cdc96",
			},
		},
		{
			name:          "Test hashing a string with RIPEMD160",
			operationName: "hash-text",
			wantErr:       false,
			data: map[string]interface{}{
				"algorithm":       "RIPEMD160",
				"text":            "Banana is yellow",
				"expected_output": "d67f31188141a33a986ad947e0a8098fc79bff9f",
			},
		},
		{
			name:          "Test string generator - All types",
			operationName: "generate-random-text",
			wantErr:       false,
			data: map[string]interface{}{
				"hasDigits":       true,
				"hasUppercase":    true,
				"hasLowercase":    true,
				"hasSpecialChars": true,
				"length":          16,
			},
		},
		{
			name:          "Test string generator - Digits only",
			operationName: "generate-random-text",
			wantErr:       false,
			data: map[string]interface{}{
				"hasDigits":       true,
				"hasUppercase":    false,
				"hasLowercase":    false,
				"hasSpecialChars": false,
				"length":          16,
			},
		},
		{
			name:          "Test string generator - Uppercase only",
			operationName: "generate-random-text",
			wantErr:       false,
			data: map[string]interface{}{
				"hasDigits":       false,
				"hasUppercase":    true,
				"hasLowercase":    false,
				"hasSpecialChars": false,
				"length":          16,
			},
		},
		{
			name:          "Test string generator - Lowercase only",
			operationName: "generate-random-text",
			wantErr:       false,
			data: map[string]interface{}{
				"hasDigits":       false,
				"hasUppercase":    false,
				"hasLowercase":    true,
				"hasSpecialChars": false,
				"length":          16,
			},
		},
		{
			name:          "Test string generator - Symbols only",
			operationName: "generate-random-text",
			wantErr:       false,
			data: map[string]interface{}{
				"hasDigits":       false,
				"hasUppercase":    false,
				"hasLowercase":    false,
				"hasSpecialChars": true,
				"length":          16,
			},
		},
		{
			name:          "Test string generator - No types error",
			operationName: "generate-random-text",
			wantErr:       true,
			data: map[string]interface{}{
				"hasDigits":       false,
				"hasUppercase":    false,
				"hasLowercase":    false,
				"hasSpecialChars": false,
				"length":          16,
			},
		},
		{
			name:          "Test string generator - Min length error",
			operationName: "generate-random-text",
			wantErr:       true,
			data: map[string]interface{}{
				"hasDigits":       true,
				"hasUppercase":    true,
				"hasLowercase":    true,
				"hasSpecialChars": true,
				"length":          0,
			},
		},
		{
			name:          "Test string generator - Max length error",
			operationName: "generate-random-text",
			wantErr:       true,
			data: map[string]interface{}{
				"hasDigits":       true,
				"hasUppercase":    true,
				"hasLowercase":    true,
				"hasSpecialChars": true,
				"length":          129,
			},
		},
	}

	ctx := &sdk.RunContext{
		Auth: &sdkcore.AuthContext{
			AccessToken: "",
			Token:       nil,
			TokenType:   "",
			Username:    "",
			Password:    "",
			Secret:      "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			instance, _ := NewConnector()
			spider := sdk.NewSpiderTest(t, instance)
			connector := spider.GetConfig()

			if connector == nil {
				t.Fatal("NewConnector() returned nil")
			}

			if connector.DisplayName != "Cryptography" {
				t.Errorf("NewConnector() Name = %s, want %s", connector.DisplayName, "Cryptography")
			}

			if connector.Logo != "tabler:circle-key-filled" {
				t.Errorf("NewConnector() Logo = %s, want %s", connector.Logo, "tabler:circle-key-filled")
			}

			if connector.Version != "0.0.1" {
				t.Errorf("NewConnector() Version = %s, want %s", connector.Version, "0.0.1")
			}

			if connector.Group != sdk.ConnectorGroupCore {
				t.Errorf("NewConnector() Group = %v, want %v", connector.Group, sdk.ConnectorGroupCore)
			}

			if len(connector.Authors) != 1 || connector.Authors[0] != "Wakflo <integrations@wakflo.com>" {
				t.Errorf("NewConnector() Authors = %v, want %v", connector.Authors, []string{"Wakflo <integrations@wakflo.com>"})
			}

			if len(spider.Triggers()) != 0 {
				t.Errorf("NewConnector() Triggers() count = %d, want %d", len(spider.Triggers()), 0)
			}

			if len(spider.Operations()) != 2 {
				t.Errorf("NewConnector() Operations() count = %d, want %d", len(spider.Operations()), 2)
			}

			ctx.Input = testCase.data
			_, _ = spider.RunOperation(testCase.operationName, ctx)
			/*if err != nil {
				if testCase.wantErr {
					fmt.Println("Successfully thrown an error")
					return
				}

				t.Errorf("NewConnector() RunOperation() with name %v threw an error = %v", testCase.operationName, err)
			}

			resultJson := result.(map[string]interface{})

			if resultJson["hashed_text"] != testCase.data["expected_output"] {
				t.Errorf("NewConnector() RunOperation() response = %v, want %v", resultJson["hashed_string"], testCase.data["expected_output"])
			}

			if resultJson["generated_text"] != nil {
				fmt.Println(resultJson["generated_string"])
			}*/
		})
	}
}
