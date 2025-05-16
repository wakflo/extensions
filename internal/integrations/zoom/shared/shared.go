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
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/juicycleff/smartform/v1"
)

// #nosec
var tokenURL = ZoomBaseURL + "/oauth/token"

var (
	form = smartform.NewAuthForm("zoom-auth", "Zoom Oauth", smartform.AuthStrategyOAuth2)
	_    = form.OAuthField("oauth", "Zoom Oauth").
		AuthorizationURL(ZoomBaseURL + "/oauth/authorize").
		TokenURL(tokenURL).
		Scopes([]string{}).
		Build()
)

var ZoomSharedAuth = form.Build()

const ZoomBaseURL = "https://api.zoom.us"

func ZoomRequest(accessToken, reqURL string, request []byte) (interface{}, error) {
	fullURL := ZoomBaseURL + reqURL

	req, err := http.NewRequest(http.MethodPost, fullURL, bytes.NewBuffer(request))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	res, errs := client.Do(req)
	if errs != nil {
		return nil, errs
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response interface{}
	if newErrs := json.Unmarshal(body, &response); newErrs != nil {
		return nil, errors.New("error parsing response")
	}

	return response, nil
}
