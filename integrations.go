package extensions

import (
	"github.com/wakflo/extensions/internal/integrations/aftership"
	"github.com/wakflo/extensions/internal/integrations/asana"
	"github.com/wakflo/extensions/internal/integrations/easyship"
	"github.com/wakflo/extensions/internal/integrations/facebookpages"
	"github.com/wakflo/extensions/internal/integrations/googledrive"
	"github.com/wakflo/extensions/internal/integrations/jsonconverter"
	"github.com/wakflo/extensions/internal/integrations/sendbox"
	"github.com/wakflo/extensions/internal/integrations/smartsheet"
	"github.com/wakflo/extensions/internal/integrations/square"
	"github.com/wakflo/extensions/internal/integrations/toggl"
	"github.com/wakflo/extensions/internal/integrations/trackingmore"
	"github.com/wakflo/extensions/internal/integrations/zohosalesiq"
	"github.com/wakflo/extensions/internal/integrations/zoom"
	"github.com/wakflo/go-sdk/integration"
)

func RegisterIntegrations() map[string]integration.RegistrationMap {
	// 🛑Do-Not-Edit
	reg := externalRegistry{
		integrations: make(map[string]integration.RegistrationMap),
	}

	plugins := []*integration.Registration{
		// 👋 Add connectors here
		googledrive.Integration,   // Google Drive
		asana.Integration,         // Asana
		aftership.Integration,     // AfterShip
		smartsheet.Integration,    // SmartSheet
		jsonconverter.Integration, // JsonConverter
		zohosalesiq.Integration,   // ZohoSales iq
		toggl.Integration,         // Toggl
		square.Integration,        // Square
		trackingmore.Integration,  // TrackingMore
		zoom.Integration,          // Zoom
		easyship.Integration,      // EasyShip
		sendbox.Integration,       // Sendbox
		facebookpages.Integration, // Facebook Pages
	}

	// 🛑Do-Not-Edit
	for _, plugin := range plugins {
		reg.addRegistration(plugin)
	}

	return reg.integrations
}

type externalRegistry struct {
	integrations map[string]integration.RegistrationMap
}

func (i *externalRegistry) addRegistration(plugin *integration.Registration) {
	_, ok := i.integrations[plugin.Info.Name]
	if !ok {
		i.integrations[plugin.Info.Name] = integration.RegistrationMap{
			Info: integration.RegistrationInfo{
				IntegrationSchemaModel: plugin.Info.IntegrationSchemaModel,
				DisplayName:            plugin.Info.DisplayName,
				Documentation:          plugin.Info.Documentation,
			},
			Versions: map[string]integration.Registration{},
		}
	}

	i.integrations[plugin.Info.Name].Versions[plugin.Info.Version] = *plugin
}
