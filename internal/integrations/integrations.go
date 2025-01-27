package integrations

import (
	"github.com/wakflo/extensions/internal/integrations/googledrive"
	"github.com/wakflo/extensions/internal/logger"
	"github.com/wakflo/go-sdk/integration"
)

var log = logger.NewDefaultLogger("external-connectors")

func RegisterIntegrations() map[string]integration.RegistrationMap {
	// 🛑Do-Not-Edit
	reg := externalRegistry{
		integrations: make(map[string]integration.RegistrationMap),
	}

	plugins := []*integration.Registration{
		// 👋 Add connectors here
		googledrive.Integration, // Google Drive
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
