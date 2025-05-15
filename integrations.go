package extensions

import (
	"github.com/gosimple/slug"
	// "github.com/wakflo/extensions/internal/integrations/activecampaign"
	// "github.com/wakflo/extensions/internal/integrations/aftership"
	// "github.com/wakflo/extensions/internal/integrations/airtable"
	// "github.com/wakflo/extensions/internal/integrations/asana"
	// "github.com/wakflo/extensions/internal/integrations/calculator"
	// "github.com/wakflo/extensions/internal/integrations/calendly"
	// "github.com/wakflo/extensions/internal/integrations/campaignmonitor"
	// "github.com/wakflo/extensions/internal/integrations/clickup"
	// "github.com/wakflo/extensions/internal/integrations/convertkit"
	// "github.com/wakflo/extensions/internal/integrations/cryptography"
	// "github.com/wakflo/extensions/internal/integrations/csv"
	// "github.com/wakflo/extensions/internal/integrations/easyship"
	// "github.com/wakflo/extensions/internal/integrations/freshdesk"
	// "github.com/wakflo/extensions/internal/integrations/gemini"
	// "github.com/wakflo/extensions/internal/integrations/github"
	// "github.com/wakflo/extensions/internal/integrations/googlecalendar"
	// "github.com/wakflo/extensions/internal/integrations/googledocs"
	"github.com/wakflo/extensions/internal/integrations/googledrive"
	// "github.com/wakflo/extensions/internal/integrations/googlemail"
	// "github.com/wakflo/extensions/internal/integrations/googlesheets"
	//
	// "github.com/wakflo/extensions/internal/integrations/gumroad"
	// "github.com/wakflo/extensions/internal/integrations/harvest"
	// "github.com/wakflo/extensions/internal/integrations/hubspot"
	// "github.com/wakflo/extensions/internal/integrations/instagram"
	// "github.com/wakflo/extensions/internal/integrations/jiracloudsoftware"
	// "github.com/wakflo/extensions/internal/integrations/jsonconverter"
	// "github.com/wakflo/extensions/internal/integrations/keapcrm"
	// "github.com/wakflo/extensions/internal/integrations/linear"
	// "github.com/wakflo/extensions/internal/integrations/mailchimp"
	"github.com/wakflo/extensions/internal/integrations/mailjet"
	// "github.com/wakflo/extensions/internal/integrations/monday"
	// "github.com/wakflo/extensions/internal/integrations/notion"
	// "github.com/wakflo/extensions/internal/integrations/openai"
	// "github.com/wakflo/extensions/internal/integrations/prisync"
	// "github.com/wakflo/extensions/internal/integrations/shopify"
	// "github.com/wakflo/extensions/internal/integrations/smartsheet"
	// "github.com/wakflo/extensions/internal/integrations/square"
	// "github.com/wakflo/extensions/internal/integrations/todoist"
	// "github.com/wakflo/extensions/internal/integrations/toggl"
	// "github.com/wakflo/extensions/internal/integrations/trackingmore"
	// "github.com/wakflo/extensions/internal/integrations/trello"
	// "github.com/wakflo/extensions/internal/integrations/typeform"
	// "github.com/wakflo/extensions/internal/integrations/woocommerce"
	// "github.com/wakflo/extensions/internal/integrations/wrike"
	// "github.com/wakflo/extensions/internal/integrations/xero"
	// "github.com/wakflo/extensions/internal/integrations/zohocrm"
	// "github.com/wakflo/extensions/internal/integrations/zohoinventory"
	// "github.com/wakflo/extensions/internal/integrations/zohosalesiq"
	// "github.com/wakflo/extensions/internal/integrations/zoom"
	"github.com/wakflo/go-sdk/v2"
)

func RegisterIntegrations() sdk.IntegrationsRegistrar {
	// 🛑Can-Edit
	plugins := []sdk.Integration{
		// 👋 Add connectors here
		googledrive.Integration, // Google Drive
		// asana.Integration,             // Asana
		// aftership.Integration,         // AfterShip
		// smartsheet.Integration,        // SmartSheet --fix
		// jsonconverter.Integration,     // JsonConverter --fix
		// zohosalesiq.Integration,       // ZohoSales iq --fix
		// toggl.Integration,             // Toggl --fix
		// square.Integration,            // Square --fix
		// trackingmore.Integration,      // TrackingMore --fix
		// zoom.Integration,              // Zoom
		// easyship.Integration,          // EasyShip
		// airtable.Integration,          // Airtable
		// calendly.Integration,          // Calendly
		// calculator.Integration,        // Calculator
		// zohoinventory.Integration,     // ZohoInventory
		// xero.Integration,              // Xero
		// woocommerce.Integration,       // WooCommerce
		// trello.Integration,            // Trello
		// shopify.Integration,           // Shopify
		// prisync.Integration,           // Prisync
		// openai.Integration,            // OpenAI
		// gemini.Integration,            // Gemini
		// monday.Integration, // Monday
		// mailchimp.Integration,         // Mailchimp
		// csv.Integration,               // CSV
		// cryptography.Integration,      // Cryptography
		// notion.Integration, // Notion
		// harvest.Integration,           // Harvest
		// googlesheets.Integration,      // Google Sheets
		// googlemail.Integration,        // Google Mail
		// googlecalendar.Integration,    // Google Calendar
		// googledocs.Integration,        // Google Docs
		// todoist.Integration,           // Todoist
		// linear.Integration,            // Linear
		// github.Integration,            // Github
		// instagram.Integration,         // Instagram
		// hubspot.Integration,           // Hubspot
		// zohocrm.Integration,           // Zoho CRM
		// freshdesk.Integration,         // Freshdesk
		// keapcrm.Integration,           // KeapCRM
		// activecampaign.Integration,    // ActiveCampaign
		// convertkit.Integration,        // ConvertKit
		// campaignmonitor.Integration,   // Campaign Monitor
		mailjet.Integration, // Mailjet
		// clickup.Integration,           // ClickUp
		// wrike.Integration,             // Wrike
		// jiracloudsoftware.Integration, // Jira Cloud Software
		// gumroad.Integration,           // Gumroad Software
		// typeform.Integration,          // Typeform Software
	}

	// 🛑Do-Not-Edit
	reg := externalRegistry{
		integrations: make(sdk.IntegrationsRegistrar, len(plugins)),
	}

	// 🛑Do-Not-Edit
	for _, plugin := range plugins {
		reg.addRegistration(plugin)
	}

	return reg.integrations
}

type externalRegistry struct {
	integrations map[string]sdk.RegistrationMap
}

func (i *externalRegistry) addRegistration(plugin sdk.Integration) {
	name := slug.Make(plugin.Metadata().Name)
	_, ok := i.integrations[name]
	if !ok {
		i.integrations[name] = sdk.RegistrationMap{
			Versions: make(map[string]sdk.Integration),
		}
	}

	i.integrations[name].Versions[plugin.Metadata().Version] = plugin
}
