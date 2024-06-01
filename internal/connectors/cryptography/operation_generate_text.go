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
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

const (
	Symbols         = "!@#$%^&*()_+-=[]{}\\|;':\",.<>/?`~"
	Digits          = "0123456789"
	UpperLetters    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	LowerLetters    = "abcdefghijklmnopqrstuvwxyz"
	MaxStringLength = 128
	MinStringLength = 1
)

type generateTextOperationProps struct {
	HasDigits       bool `json:"hasDigits"`
	HasUppercase    bool `json:"hasUppercase"`
	HasLowercase    bool `json:"hasLowercase"`
	HasSpecialChars bool `json:"hasSpecialChars"`
	Length          int  `json:"length"`
}

type GenerateTextOperation struct {
	options *sdk.OperationInfo
}

func NewGenerateString() *GenerateTextOperation {
	return &GenerateTextOperation{
		options: &sdk.OperationInfo{
			Name:        "Generate random text",
			Description: "Generates a random text / password (with a specified length)",
			RequireAuth: false,
			Input: map[string]*sdkcore.AutoFormSchema{
				"hasDigits": autoform.NewBooleanField().
					SetDisplayName("Numeric digits").
					SetDescription("Text can have digits (0-9)").
					SetDefaultValue(true).
					Build(),
				"hasUppercase": autoform.NewBooleanField().
					SetDisplayName("Uppercase letters").
					SetDescription("Text can have uppercase letters (A-Z)").
					SetDefaultValue(true).
					Build(),
				"hasLowercase": autoform.NewBooleanField().
					SetDisplayName("Lowercase letters").
					SetDescription("Text can have lowercase letters (a-z)").
					SetDefaultValue(true).
					Build(),
				"hasSpecialChars": autoform.NewBooleanField().
					SetDisplayName("Symbols").
					SetDescription(fmt.Sprintf("Text can have symbols (%s)", Symbols)).
					SetDefaultValue(false).
					Build(),
				"length": autoform.NewNumberField().
					SetDisplayName("Length").
					SetDescription("Text length").
					SetMaximum(MaxStringLength).
					SetMinimum(MinStringLength).
					SetDefaultValue(16).
					SetRequired(true).
					Build(),
			},
			SampleOutput: map[string]interface{}{
				"name":              "cryptography-generate-random-text",
				"usage_mode":        "operation",
				"has_digits":        true,
				"has_uppercase":     true,
				"has_lowercase":     true,
				"has_special_chars": true,
				"generated_text":    ":L)Z{VD_u6=w,=AwQ_Q_?A4xXfJn126G",
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *GenerateTextOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	input := sdk.InputToType[generateTextOperationProps](ctx)

	var charset string

	hasDigits := input.HasDigits
	hasUpper := input.HasUppercase
	hasLower := input.HasLowercase
	hasSymbols := input.HasSpecialChars
	length := input.Length

	if !hasDigits && !hasUpper && !hasLower && !hasSymbols {
		return nil, errors.New("text must at least have one type of charset")
	}

	if length < MinStringLength || length > MaxStringLength {
		return nil, fmt.Errorf("text length must be between %d and %d", MinStringLength, MaxStringLength)
	}

	if hasDigits {
		charset += Digits
	}

	if hasUpper {
		charset += UpperLetters
	}

	if hasLower {
		charset += LowerLetters
	}

	if hasSymbols {
		charset += Symbols
	}

	randomBytes := make([]byte, length)

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return nil, err
		}
		randomBytes[i] = charset[num.Int64()]
	}

	return map[string]interface{}{
		"name":              "cryptography-generate-random-text",
		"usage_mode":        "operation",
		"has_digits":        hasDigits,
		"has_uppercase":     hasUpper,
		"has_lowercase":     hasLower,
		"has_special_chars": hasSymbols,
		"generated_text":    string(randomBytes),
	}, nil
}

func (c *GenerateTextOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *GenerateTextOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
