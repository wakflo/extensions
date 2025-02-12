package actions

import (
	"context"

	"github.com/wakflo/extensions/internal/integrations/shopify/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type getLocationsActionProps struct {
	Name string `json:"name"`
}

type GetLocationsAction struct{}

func (a *GetLocationsAction) Name() string {
	return "Get Locations"
}

func (a *GetLocationsAction) Description() string {
	return "Retrieves a list of locations from a specified data source or system, allowing you to integrate with various mapping and geolocation services."
}

func (a *GetLocationsAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *GetLocationsAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getLocationsDocs,
	}
}

func (a *GetLocationsAction) Icon() *string {
	return nil
}

func (a *GetLocationsAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{}
}

func (a *GetLocationsAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	client, err := shared.CreateClient(ctx.BaseContext)
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

func (a *GetLocationsAction) Auth() *sdk.Auth {
	return nil
}

func (a *GetLocationsAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *GetLocationsAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewGetLocationsAction() sdk.Action {
	return &GetLocationsAction{}
}
