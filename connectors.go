package extensions

// Copyright 2022-present Wakflo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	"github.com/wakflo/extensions/internal/connectors/airtable"
	"github.com/wakflo/extensions/internal/connectors/calculator"
	"github.com/wakflo/extensions/internal/connectors/calendly"
	"github.com/wakflo/extensions/internal/connectors/cin7"
	"github.com/wakflo/extensions/internal/connectors/clickup"
	"github.com/wakflo/extensions/internal/connectors/cryptography"
	"github.com/wakflo/extensions/internal/connectors/delay"
	"github.com/wakflo/extensions/internal/connectors/dropbox"
	"github.com/wakflo/extensions/internal/connectors/flexport"
	"github.com/wakflo/extensions/internal/connectors/freshdesk"
	"github.com/wakflo/extensions/internal/connectors/freshworkscrm"
	"github.com/wakflo/extensions/internal/connectors/github"
	googledocs "github.com/wakflo/extensions/internal/connectors/google_docs"
	googlesheets "github.com/wakflo/extensions/internal/connectors/google_sheets"
	"github.com/wakflo/extensions/internal/connectors/googlecalendar"
	"github.com/wakflo/extensions/internal/connectors/googledrive"
	"github.com/wakflo/extensions/internal/connectors/googlemail"
	"github.com/wakflo/extensions/internal/connectors/harvest"
	"github.com/wakflo/extensions/internal/connectors/hubspot"
	"github.com/wakflo/extensions/internal/connectors/jiracloud"
	"github.com/wakflo/extensions/internal/connectors/linear"
	"github.com/wakflo/extensions/internal/connectors/mailchimp"
	"github.com/wakflo/extensions/internal/connectors/monday"
	"github.com/wakflo/extensions/internal/connectors/notion"
	"github.com/wakflo/extensions/internal/connectors/openai"
	"github.com/wakflo/extensions/internal/connectors/prisync"
	"github.com/wakflo/extensions/internal/connectors/shippo"
	"github.com/wakflo/extensions/internal/connectors/shopify"
	"github.com/wakflo/extensions/internal/connectors/slack"
	"github.com/wakflo/extensions/internal/connectors/smartsheet"
	"github.com/wakflo/extensions/internal/connectors/stripe"
	"github.com/wakflo/extensions/internal/connectors/todoist"
	"github.com/wakflo/extensions/internal/connectors/trello"
	"github.com/wakflo/extensions/internal/connectors/woocommerce"
	"github.com/wakflo/extensions/internal/connectors/xero"
	"github.com/wakflo/extensions/internal/connectors/zohoinventory"
	"github.com/wakflo/extensions/internal/connectors/zoom"
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
		googledrive.NewConnector,  // Google Drive
		googlesheets.NewConnector, // Google Sheets
		googledocs.NewConnector,   // Google Docs
		googlemail.NewConnector,   // Gmail
		slack.NewConnector,        // Slack
		// javascript.NewConnector,     // Javascript
		cryptography.NewConnector,   // Cryptography
		delay.NewConnector,          // Delay
		todoist.NewConnector,        // Todoist
		calculator.NewConnector,     // Calculator
		shopify.NewConnector,        // Shopify
		zohoinventory.NewConnector,  // Zoho Inventory
		cin7.NewConnector,           // Cin7
		woocommerce.NewConnector,    // Woocommerce
		mailchimp.NewConnector,      // MailChimp
		xero.NewConnector,           // Xero
		clickup.NewConnector,        // Clickup
		freshdesk.NewConnector,      // Freshdesk
		linear.NewConnector,         // Linear
		freshworkscrm.NewConnector,  // Freshworks CRM
		calendly.NewConnector,       // Calendly
		shippo.NewConnector,         // Shippo
		dropbox.NewConnector,        // Dropbox
		harvest.NewConnector,        // HubStaff
		airtable.NewConnector,       // Airtable
		stripe.NewConnector,         // Stripe
		openai.NewConnector,         // OpenAI
		googlecalendar.NewConnector, // Google Calendar
		monday.NewConnector,         // Monday.com
		zoom.NewConnector,           // Zoom
		flexport.NewConnector,       // Flexport
		hubspot.NewConnector,        // Hubspot
		jiracloud.NewConnector,      // Jira Cloud
		prisync.NewConnector,        // Prisync
		github.NewConnector,         // Github
		trello.NewConnector,         // Trello
		notion.NewConnector,         // Notion
		smartsheet.NewConnector,     // Smartsheet
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
