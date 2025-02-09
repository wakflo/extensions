package cryptography

import (
	// #nosec
	"crypto/md5"
	// #nosec
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"hash"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type CryptographyAlgorithm struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var suppportedAlgorithms = []*sdkcore.AutoFormSchema{
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

type hashTextOperationProps struct {
	Algorithm string `json:"algorithm"`
	Text      string `json:"text"`
}

type HashTextOperation struct {
	options *sdk.OperationInfo
}

func NewHashString() *HashTextOperation {
	return &HashTextOperation{
		options: &sdk.OperationInfo{
			Name:        "Hash text",
			Description: "Hashes given text with chosen algorithm",
			RequireAuth: false,
			Input: map[string]*sdkcore.AutoFormSchema{
				"algorithm": autoform.NewSelectField().
					SetDisplayName("Algorithm").
					SetDescription("Hashing algorithm").
					SetOptions(suppportedAlgorithms).
					SetRequired(true).
					Build(),
				"text": autoform.NewShortTextField().
					SetDisplayName("Text").
					SetDescription("Text to be hashed").
					SetRequired(true).
					Build(),
			},
			SampleOutput: map[string]interface{}{
				"name":         "cryptography-hash-text",
				"usage_mode":   "operation",
				"algorithm_id": "MD5",
				"text":         "This is the string to be hashed",
				"hashed_text":  "b2ba2dca5dd2ae6a18bb3f72ec1b9389",
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *HashTextOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	input := sdk.InputToType[hashTextOperationProps](ctx)

	algo := input.Algorithm
	inputText := input.Text

	var hashedString hash.Hash

	switch algo {
	case suppportedAlgorithms[0].Const:
		// #nosec
		hashedString = md5.New()
		hashedString.Write([]byte(inputText))
		break
	case suppportedAlgorithms[1].Const:
		// #nosec
		hashedString = sha1.New()
		hashedString.Write([]byte(inputText))
		break
	case suppportedAlgorithms[2].Const:
		hashedString = sha256.New()
		hashedString.Write([]byte(inputText))
		break
	case suppportedAlgorithms[3].Const:
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

func (c *HashTextOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *HashTextOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
