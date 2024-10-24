package calculator

import (
	"math"
	"testing"

	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

func isEqualWithTolerance(a float64, b float64) bool {
	return math.Abs(a-b) <= 1e-9
}

func TestNewConnector(t *testing.T) {
	testCases := []struct {
		name          string
		operationName string
		wantErr       bool
		data          map[string]interface{}
	}{
		{
			name:          "Test addition - 1",
			operationName: "addition",
			wantErr:       false,
			data: map[string]interface{}{
				"firstNumber":     5,
				"secondNumber":    10,
				"expected_output": 15.0,
			},
		},
		{
			name:          "Test addition - 2",
			operationName: "addition",
			wantErr:       false,
			data: map[string]interface{}{
				"firstNumber":     -5,
				"secondNumber":    10,
				"expected_output": 5.0,
			},
		},
		{
			name:          "Test addition - 3",
			operationName: "addition",
			wantErr:       false,
			data: map[string]interface{}{
				"firstNumber":     23948320480932,
				"secondNumber":    48324328940923,
				"expected_output": 7.2272649421855e+13,
			},
		},
		{
			name:          "Test subtraction - 1",
			operationName: "subtraction",
			wantErr:       false,
			data: map[string]interface{}{
				"firstNumber":     329,
				"secondNumber":    29,
				"expected_output": 300.0,
			},
		},
		{
			name:          "Test subtraction - 2",
			operationName: "subtraction",
			wantErr:       false,
			data: map[string]interface{}{
				"firstNumber":     5,
				"secondNumber":    10,
				"expected_output": -5.0,
			},
		},
		{
			name:          "Test subtraction - 3",
			operationName: "subtraction",
			wantErr:       false,
			data: map[string]interface{}{
				"firstNumber":     -5,
				"secondNumber":    -5,
				"expected_output": 0.0,
			},
		},
		{
			name:          "Test division - 1",
			operationName: "division",
			wantErr:       false,
			data: map[string]interface{}{
				"firstNumber":     1,
				"secondNumber":    2,
				"expected_output": 0.5,
			},
		},
		{
			name:          "Test division - 2",
			operationName: "division",
			wantErr:       false,
			data: map[string]interface{}{
				"firstNumber":     329,
				"secondNumber":    0,
				"expected_output": "Unknown",
			},
		},
		{
			name:          "Test division - 3",
			operationName: "division",
			wantErr:       false,
			data: map[string]interface{}{
				"firstNumber":     0,
				"secondNumber":    29,
				"expected_output": 0.0,
			},
		},
		{
			name:          "Test division - 4",
			operationName: "division",
			wantErr:       false,
			data: map[string]interface{}{
				"firstNumber":     683,
				"secondNumber":    29,
				"expected_output": 23.5517241379,
			},
		},
		{
			name:          "Test multiplication - 1",
			operationName: "multiplication",
			wantErr:       false,
			data: map[string]interface{}{
				"firstNumber":     5,
				"secondNumber":    5,
				"expected_output": 25.0,
			},
		},
		{
			name:          "Test multiplication - 2",
			operationName: "multiplication",
			wantErr:       false,
			data: map[string]interface{}{
				"firstNumber":     5,
				"secondNumber":    -1,
				"expected_output": -5.0,
			},
		},
		{
			name:          "Test multiplication - 3",
			operationName: "multiplication",
			wantErr:       false,
			data: map[string]interface{}{
				"firstNumber":     25,
				"secondNumber":    502,
				"expected_output": 12550.0,
			},
		},
		{
			name:          "Test multiplication - 3",
			operationName: "multiplication",
			wantErr:       false,
			data: map[string]interface{}{
				"firstNumber":     25,
				"secondNumber":    0,
				"expected_output": 0.0,
			},
		},
		{
			name:          "Test modulo - 1",
			operationName: "modulo",
			wantErr:       false,
			data: map[string]interface{}{
				"firstNumber":     4,
				"secondNumber":    2,
				"expected_output": 0.0,
			},
		},
		{
			name:          "Test modulo - 2",
			operationName: "modulo",
			wantErr:       false,
			data: map[string]interface{}{
				"firstNumber":     5,
				"secondNumber":    2,
				"expected_output": 1.0,
			},
		},
		{
			name:          "Test modulo - 3",
			operationName: "modulo",
			wantErr:       false,
			data: map[string]interface{}{
				"firstNumber":     50,
				"secondNumber":    3,
				"expected_output": 2.0,
			},
		},
		{
			name:          "Test modulo - 4",
			operationName: "modulo",
			wantErr:       false,
			data: map[string]interface{}{
				"firstNumber":     2,
				"secondNumber":    0.0,
				"expected_output": 0.0,
			},
		},
		{
			name:          "Test modulo - 4",
			operationName: "modulo",
			wantErr:       false,
			data: map[string]interface{}{
				"firstNumber":     0,
				"secondNumber":    2,
				"expected_output": 0.0,
			},
		},
	}

	ctx := &sdk.RunContext{
		Input: map[string]interface{}{},
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

			if connector.DisplayName != "Calculator" {
				t.Errorf("NewConnector() Name = %s, want %s", connector.DisplayName, "Cryptography")
			}

			if connector.Logo != "ion:calculator" {
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

			if len(spider.Operations()) != 5 {
				t.Errorf("NewConnector() Operations() count = %d, want %d", len(spider.Operations()), 2)
			}

			ctx.Input = testCase.data

			result, err := spider.RunOperation(testCase.operationName, ctx)
			if err != nil {
				t.Errorf("NewConnector() RunOperation() with name %v threw an error = %v", testCase.operationName, err)
			}

			resultJSON := result.(map[string]interface{})

			switch resultJSON["result"].(type) {
			case string:
				if resultJSON["result"] != testCase.data["expected_output"] {
					t.Errorf("NewConnector() RunOperation() response = %v, want %v", resultJSON["result"], testCase.data["expected_output"])
				}
			default:
				if !isEqualWithTolerance(resultJSON["result"].(float64), testCase.data["expected_output"].(float64)) {
					t.Errorf("NewConnector() RunOperation() response = %v, want %v", resultJSON["result"], testCase.data["expected_output"])
				}
			}
		})
	}
}
