package actions

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"hash"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type hashTextActionProps struct {
	Algorithm string `json:"algorithm"`
	Text      string `json:"text"`
}

type HashTextAction struct{}

// Metadata returns metadata about the action
func (a *HashTextAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "hash_text",
		DisplayName:   "Hash Text",
		Description:   "Hashes the input text and returns a unique digital fingerprint (hash value) that can be used to verify the integrity of the original text.",
		Type:          core.ActionTypeAction,
		Documentation: hashTextDocs,
		SampleOutput: map[string]any{
			"name":         "cryptography-hash-text",
			"usage_mode":   "operation",
			"algorithm_id": "MD5",
			"text":         "This is the string to be hashed",
			"hashed_text":  "b2ba2dca5dd2ae6a18bb3f72ec1b9389",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *HashTextAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("hash_text", "Hash Text")

	// Add algorithm field
	form.SelectField("algorithm", "Algorithm").
		Placeholder("Select a hashing algorithm").
		Required(true).
		HelpText("Hashing algorithm").
		AddOption("MD5", "MD5").
		AddOption("SHA1", "SHA1").
		AddOption("SHA256", "SHA256").
		AddOption("SHA512", "SHA512")

	// Add text field
	form.TextField("text", "Text").
		Placeholder("Enter text to hash").
		Required(true).
		HelpText("Text to be hashed")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *HashTextAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *HashTextAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[hashTextActionProps](ctx)
	if err != nil {
		return nil, err
	}

	algo := input.Algorithm
	inputText := input.Text

	var hashedString hash.Hash

	switch algo {
	case "MD5":
		// #nosec
		hashedString = md5.New()
		hashedString.Write([]byte(inputText))
	case "SHA1":
		// #nosec
		hashedString = sha1.New()
		hashedString.Write([]byte(inputText))
	case "SHA256":
		hashedString = sha256.New()
		hashedString.Write([]byte(inputText))
	case "SHA512":
		hashedString = sha512.New()
		hashedString.Write([]byte(inputText))
	default:
		return nil, errors.New("given hashing algorithm isn't supported")
	}

	return map[string]interface{}{
		"name":         "cryptography-hash-text",
		"usage_mode":   "operation",
		"algorithm_id": algo,
		"text":         inputText,
		"hashed_text":  hex.EncodeToString(hashedString.Sum(nil)),
	}, nil
}

func NewHashTextAction() sdk.Action {
	return &HashTextAction{}
}
