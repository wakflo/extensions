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

	//nolint:staticcheck,gosec
	"golang.org/x/crypto/ripemd160"
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
	{
		Const: "RIPEMD160",
		Title: "RIPEMD160",
	},
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
	case suppportedAlgorithms[4].Const:
		//nolint:gosec
		hashedString = ripemd160.New()
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
