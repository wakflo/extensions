package actions

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"

	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

/*
const (

	Symbols         = "!@#$%^&*()_+-=[]{}\\|;':\",.<>/?`~"
	Digits          = "0123456789"
	UpperLetters    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	LowerLetters    = "abcdefghijklmnopqrstuvwxyz"
	MaxStringLength = 128
	MinStringLength = 1

)
*/
const (
	// Symbols definition
	Symbols = "!@#$%^&*()_+-=[]{}\\|;':\",.<>/?`~"

	// Digits definition
	Digits = "0123456789"

	// UpperLetters Upper-case letters
	UpperLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// LowerLetters Lower-case letters
	LowerLetters = "abcdefghijklmnopqrstuvwxyz"

	// MaxStringLength Maximum string length allowed
	MaxStringLength = 128

	// MinStringLength Minimum string length allowed
	MinStringLength = 1
)

type generateTextActionProps struct {
	HasDigits       bool `json:"hasDigits"`
	HasUppercase    bool `json:"hasUppercase"`
	HasLowercase    bool `json:"hasLowercase"`
	HasSpecialChars bool `json:"hasSpecialChars"`
	Length          int  `json:"length"`
}

type GenerateTextAction struct{}

func (a *GenerateTextAction) Name() string {
	return "Generate Text"
}

func (a *GenerateTextAction) Description() string {
	return "Generates text based on user-defined templates and variables, allowing you to create dynamic and personalized content for various use cases."
}

func (a *GenerateTextAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *GenerateTextAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &generateTextDocs,
	}
}

func (a *GenerateTextAction) Icon() *string {
	return nil
}

func (a *GenerateTextAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
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
			//nolint:mnd
			SetDefaultValue(16).
			SetRequired(true).
			Build(),
	}
}

func (a *GenerateTextAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[generateTextActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

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

	for i := range randomBytes {
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

func (a *GenerateTextAction) Auth() *sdk.Auth {
	return nil
}

func (a *GenerateTextAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"name":              "cryptography-generate-random-text",
		"usage_mode":        "operation",
		"has_digits":        true,
		"has_uppercase":     true,
		"has_lowercase":     true,
		"has_special_chars": true,
		"generated_text":    ":L)Z{VD_u6=w,=AwQ_Q_?A4xXfJn126G",
	}
}

func (a *GenerateTextAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewGenerateTextAction() sdk.Action {
	return &GenerateTextAction{}
}
