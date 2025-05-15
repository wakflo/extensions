package actions

import (
	"context"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/shopify/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type getLocationsActionProps struct {
	Name string `json:"name"`
}

type GetLocationsAction struct{}

func (a *GetLocationsAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_locations",
		DisplayName:   "Get Locations",
		Description:   "Retrieves a list of locations from a specified data source or system, allowing you to integrate with various mapping and geolocation services.",
		Type:          core.ActionTypeAction,
		Documentation: getLocationsDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *GetLocationsAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_locations", "Get Locations")

	// No properties in the original implementation

	schema := form.Build()
	return schema
}

func (a *GetLocationsAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	client, err := shared.CreateClient(ctx)
	if err != nil {
		return nil, err
	}
	locations, err := client.Location.List(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	simplifiedLocations := make([]map[string]interface{}, len(locations))
	for i, location := range locations {
		simplifiedLocations[i] = map[string]interface{}{
			"Id":       location.Id,
			"Name":     location.Name,
			"Address1": location.Address1,
			"Address2": location.Address2,
			"City":     location.City,
			"Zip":      location.Zip,
			"Country":  location.Country,
			"Phone":    location.Phone,
			"Active":   location.Active,
		}
	}
	return map[string]interface{}{
		"locations": simplifiedLocations,
	}, nil
}

func (a *GetLocationsAction) Auth() *core.AuthMetadata {
	return nil
}

func NewGetLocationsAction() sdk.Action {
	return &GetLocationsAction{}
}
