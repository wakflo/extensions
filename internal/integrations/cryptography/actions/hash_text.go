package actions

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"hash"

	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type CryptographyAlgorithm struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var supportedAlgorithms = []*sdkcore.AutoFormSchema{
	{
		Const: "MD5",
		Title: "MD5",
	},
	{
		Const: "SHA1",
		Title: "SHA1",
	},
	{
		Const: "SHA256",
		Title: "SHA256",
	},
	{
		Const: "SHA512",
		Title: "SHA512",
	},
	// RIPEMD160 removed due to deprecation
}

type hashTextActionProps struct {
	Algorithm string `json:"algorithm"`
	Text      string `json:"text"`
}

type HashTextAction struct{}

func (a *HashTextAction) Name() string {
	return "Hash Text"
}

func (a *HashTextAction) Description() string {
	return "Hashes the input text and returns a unique digital fingerprint (hash value) that can be used to verify the integrity of the original text."
}

func (a *HashTextAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *HashTextAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &hashTextDocs,
	}
}

func (a *HashTextAction) Icon() *string {
	return nil
}

func (a *HashTextAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"algorithm": autoform.NewSelectField().
			SetDisplayName("Algorithm").
			SetDescription("Hashing algorithm").
			SetOptions(supportedAlgorithms).
			SetRequired(true).
			Build(),
		"text": autoform.NewShortTextField().
			SetDisplayName("Text").
			SetDescription("Text to be hashed").
			SetRequired(true).
			Build(),
	}
}

func (a *HashTextAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[hashTextActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	algo := input.Algorithm
	inputText := input.Text

	var hashedString hash.Hash

	switch algo {
	case supportedAlgorithms[0].Const:
		// #nosec
		hashedString = md5.New()
		hashedString.Write([]byte(inputText))
		break
	case supportedAlgorithms[1].Const:
		// #nosec
		hashedString = sha1.New()
		hashedString.Write([]byte(inputText))
		break
	case supportedAlgorithms[2].Const:
		hashedString = sha256.New()
		hashedString.Write([]byte(inputText))
		break
	case supportedAlgorithms[3].Const:
		hashedString = sha512.New()
		hashedString.Write([]byte(inputText))
		break
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

func (a *HashTextAction) Auth() *sdk.Auth {
	return nil
}

func (a *HashTextAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"name":         "cryptography-hash-text",
		"usage_mode":   "operation",
		"algorithm_id": "MD5",
		"text":         "This is the string to be hashed",
		"hashed_text":  "b2ba2dca5dd2ae6a18bb3f72ec1b9389",
	}
}

func (a *HashTextAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewHashTextAction() sdk.Action {
	return &HashTextAction{}
}
