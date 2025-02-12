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
	"github.com/wakflo/extensions/internal/connectors/cin7"
	"github.com/wakflo/extensions/internal/connectors/clickup"
	"github.com/wakflo/extensions/internal/connectors/dropbox"
	"github.com/wakflo/extensions/internal/connectors/flexport"
	"github.com/wakflo/extensions/internal/connectors/freshdesk"
	"github.com/wakflo/extensions/internal/connectors/freshworkscrm"
	"github.com/wakflo/extensions/internal/connectors/hubspot"
	"github.com/wakflo/extensions/internal/connectors/jiracloud"
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
		// javascript.NewConnector,     // Javascript
		// delay.NewConnector,         // Delay
		cin7.NewConnector,          // Cin7
		clickup.NewConnector,       // Clickup
		freshdesk.NewConnector,     // Freshdesk
		freshworkscrm.NewConnector, // Freshworks CRM
		dropbox.NewConnector,       // Dropbox
		flexport.NewConnector,      // Flexport
		hubspot.NewConnector,       // Hubspot
		jiracloud.NewConnector,     // Jira Cloud
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
