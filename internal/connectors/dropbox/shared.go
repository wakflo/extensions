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

package dropbox

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/wakflo/go-sdk/autoform"
)

var (
	// #nosec
	tokenURL   = "https://api.dropboxapi.com/oauth2/token"
	sharedAuth = autoform.NewOAuthField("https://www.dropbox.com/oauth2/authorize", &tokenURL, []string{
		"files.metadata.write files.content.read files.metadata.read files.content.write",
	}).SetRequired(true).Build()
)

func dropBoxClient(reqURL, accessToken string, request []byte) (interface{}, error) {
	req, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewBuffer(request))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var dropboxResponse interface{}
	err = json.Unmarshal(body, &dropboxResponse)
	if err != nil {
		return nil, err
	}

	return dropboxResponse, nil
}

func listFolderContent(reqURL, accessToken string, request []byte) (interface{}, error) {
	req, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewBuffer(request))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil, nil
	}
	var folderContent map[string]interface{}
	err = json.Unmarshal(body, &folderContent)
	if err != nil {
		return nil, err
	}

	nodes, ok := folderContent["entries"].([]interface{})
	if !ok {
		return nil, errors.New("failed to extract issues from response")
	}

	return nodes, nil
}