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

package shopify

import (
	"context"
	"errors"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type GetLocationsOperation struct {
	options *sdk.OperationInfo
}

func NewGetLocationsOperation() *GetLocationsOperation {
	return &GetLocationsOperation{
		options: &sdk.OperationInfo{
			Name:        "Get Locations",
			Description: "get locations",
			RequireAuth: true,
			Auth:        sharedAuth,
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}
func (c *GetLocationsOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["token"] == "" {
		return nil, errors.New("missing shopify auth token")
	}
	domain := ctx.Auth.Extra["domain"]
	shopName := domain + ".myshopify.com"
	client := getShopifyClient(shopName, ctx.Auth.Extra["token"])
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
func (c *GetLocationsOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}
func (c *GetLocationsOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
