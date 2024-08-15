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
	"github.com/wakflo/extensions/internal/connectors/calculator"
	"github.com/wakflo/extensions/internal/connectors/cin7"
	"github.com/wakflo/extensions/internal/connectors/clickup"
	"github.com/wakflo/extensions/internal/connectors/cryptography"
	"github.com/wakflo/extensions/internal/connectors/delay"
	googledocs "github.com/wakflo/extensions/internal/connectors/google_docs"
	googlesheets "github.com/wakflo/extensions/internal/connectors/google_sheets"
	"github.com/wakflo/extensions/internal/connectors/googledrive"
	"github.com/wakflo/extensions/internal/connectors/googlemail"
	"github.com/wakflo/extensions/internal/connectors/goscript"
	"github.com/wakflo/extensions/internal/connectors/javascript"
	"github.com/wakflo/extensions/internal/connectors/mailchimp"
	"github.com/wakflo/extensions/internal/connectors/manual"
	"github.com/wakflo/extensions/internal/connectors/shopify"
	"github.com/wakflo/extensions/internal/connectors/slack"
	"github.com/wakflo/extensions/internal/connectors/todoist"
	"github.com/wakflo/extensions/internal/connectors/webhook"
	"github.com/wakflo/extensions/internal/connectors/woocommerce"
	"github.com/wakflo/extensions/internal/connectors/xero"
	"github.com/wakflo/extensions/internal/connectors/zohoinventory"
	"github.com/wakflo/extensions/internal/logger"
	sdk "github.com/wakflo/go-sdk/connector"
)

var log = logger.NewDefaultLogger("connectors")

func RegisterConnectors() []*sdk.ConnectorPlugin {
	// 🛑Do-Not-Edit
	reg := internalRegistry{
		connectors: []*sdk.ConnectorPlugin{},
	}

	plugins := []func() (*sdk.ConnectorPlugin, error){
		// 👋 Add connectors here
		googledrive.NewConnector,   // Google Drive
		googlesheets.NewConnector,  // Google Sheets
		googledocs.NewConnector,    // Google Docs
		googlemail.NewConnector,    // Gmail
		slack.NewConnector,         // Slack
		javascript.NewConnector,    // Javascript
		cryptography.NewConnector,  // Cryptography
		goscript.NewConnector,      // Go Lang
		delay.NewConnector,         // Delay
		todoist.NewConnector,       // Todoist
		manual.NewConnector,        // Manual
		calculator.NewConnector,    // Calculator
		shopify.NewConnector,       // Shopify
		webhook.NewConnector,       // Webhook
		zohoinventory.NewConnector, // Zoho Inventory
		cin7.NewConnector,          // Cin7
		woocommerce.NewConnector,   // Woocommerce
		mailchimp.NewConnector,     // MailChimp
    xero.NewConnector,          // Xero
		clickup.NewConnector,       // Clickup
	}

	// 🛑Do-Not-Edit
	for _, plugin := range plugins {
		reg.insertConnector(plugin())
	}

	return reg.connectors
}

type internalRegistry struct {
	connectors []*sdk.ConnectorPlugin
}

func (i *internalRegistry) insertConnector(connector *sdk.ConnectorPlugin, err error) {
	if err == nil {
		i.connectors = append(i.connectors, connector)
	} else {
		log.Error().Err(err).Msgf(err.Error())
	}
}
