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

package jsonconverter

import (
	"encoding/json"
)

func convertTextToJSON(text string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(text), &result)
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
