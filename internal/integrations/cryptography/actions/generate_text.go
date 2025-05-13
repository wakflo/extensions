package actions

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

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

// Metadata returns metadata about the action
func (a *GenerateTextAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "generate_text",
		DisplayName:   "Generate Text",
		Description:   "Generates text based on user-defined templates and variables, allowing you to create dynamic and personalized content for various use cases.",
		Type:          core.ActionTypeAction,
		Documentation: generateTextDocs,
		SampleOutput: map[string]any{
			"name":              "cryptography-generate-random-text",
			"usage_mode":        "operation",
			"has_digits":        true,
			"has_uppercase":     true,
			"has_lowercase":     true,
			"has_special_chars": true,
			"generated_text":    ":L)Z{VD_u6=w,=AwQ_Q_?A4xXfJn126G",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *GenerateTextAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("generate_text", "Generate Text")

	// Add hasDigits field
	form.CheckboxField("hasDigits", "Numeric digits").
		Required(false).
		DefaultValue(true).
		HelpText("Text can have digits (0-9)")

	// Add hasUppercase field
	form.CheckboxField("hasUppercase", "Uppercase letters").
		Required(false).
		DefaultValue(true).
		HelpText("Text can have uppercase letters (A-Z)")

	// Add hasLowercase field
	form.CheckboxField("hasLowercase", "Lowercase letters").
		Required(false).
		DefaultValue(true).
		HelpText("Text can have lowercase letters (a-z)")

	// Add hasSpecialChars field
	form.CheckboxField("hasSpecialChars", "Symbols").
		Required(false).
		DefaultValue(false).
		HelpText(fmt.Sprintf("Text can have symbols (%s)", Symbols))

	// Add length field
	form.NumberField("length", "Length").
		Placeholder("Enter text length").
		Required(true).
		DefaultValue(16).
		HelpText("Text length")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *GenerateTextAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *GenerateTextAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[generateTextActionProps](ctx)
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

func NewGenerateTextAction() sdk.Action {
	return &GenerateTextAction{}
}
