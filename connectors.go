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

package extensions

import (
	"github.com/wakflo/extensions/internal/connectors/cryptography"
	"github.com/wakflo/extensions/internal/connectors/delay"
	googledrive "github.com/wakflo/extensions/internal/connectors/google-drive"
	"github.com/wakflo/extensions/internal/connectors/goscript"
	"github.com/wakflo/extensions/internal/connectors/javascript"
	"github.com/wakflo/extensions/internal/connectors/slack"
	"github.com/wakflo/extensions/internal/connectors/todoist"
	sdk "github.com/wakflo/go-sdk/connector"
)

func RegisterConnectors() []*sdk.ConnectorPlugin {
	var connectors []*sdk.ConnectorPlugin

	// Google Drive
	gd, err := googledrive.NewConnector()
	if err == nil {
		connectors = append(connectors, gd)
	}

	// Slack
	sk, err := slack.NewConnector()
	if err == nil {
		connectors = append(connectors, sk)
	}

	// JavaScript
	js, err := javascript.NewConnector()
	if err == nil {
		connectors = append(connectors, js)
	}

	// Cryptography
	crp, err := cryptography.NewConnector()
	if err == nil {
		connectors = append(connectors, crp)
	}

	// Go Lang
	gos, err := goscript.NewConnector()
	if err == nil {
		connectors = append(connectors, gos)
	}

	// Delay
	del, err := delay.NewConnector()
	if err == nil {
		connectors = append(connectors, del)
	}

	// Todoist
	todo, err := todoist.NewConnector()
	if err == nil {
		connectors = append(connectors, todo)
	}

	return connectors
}
