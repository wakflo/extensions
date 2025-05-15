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
	"github.com/juicycleff/smartform/v1"
)

var (
	googleDocsForm = smartform.NewAuthForm("google-docs-auth", "Google Docs OAuth", smartform.AuthStrategyOAuth2)
	_              = googleDocsForm.
			OAuthField("oauth", "Google Docs OAuth").
			AuthorizationURL("https://accounts.google.com/o/oauth2/auth").
			TokenURL("https://oauth2.googleapis.com/token").
			Scopes([]string{
			"https://www.googleapis.com/auth/documents https://www.googleapis.com/auth/drive.readonly",
		}).
		Build()
)

var SharedGoogleDocsAuth = googleDocsForm.Build()
