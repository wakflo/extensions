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

package shared

import (
	"encoding/json"
	"fmt"
)

func ConvertTextToJSON(text string) (any, error) {
	var result map[string]interface{}
	err := parseDoubleEncodedJSON(text, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

//  func convertJSONToText(jsonObj map[string]interface{}) (string, error) {
//	jsonData, err := json.Marshal(jsonObj)
//	if err != nil {
//		return "", err
//	}
//	return string(jsonData), nil
//}

// parseDoubleEncodedJSON parses a double-encoded JSON string.
func parseDoubleEncodedJSON(input string, output interface{}) error {
	// Step 1: Decode the outer layer to get the inner JSON string
	var rawJSON string
	err := json.Unmarshal([]byte(input), &rawJSON)
	if err != nil {
		return fmt.Errorf("failed to unmarshal outer JSON: %w", err)
	}

	// Step 2: Decode the inner JSON string into the expected structure
	err = json.Unmarshal([]byte(rawJSON), &output)
	if err != nil {
		return fmt.Errorf("failed to unmarshal inner JSON: %w", err)
	}

	return nil
}
