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

package shippo

import (
	"errors"
	"fmt"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type createNewShipmentOperationProps struct {
	SenderName      string `json:"sender-name"`
	SenderStreet    string `json:"sender-street1"`
	SenderPhone     string `json:"sender-phone"`
	SenderZip       string `json:"sender-zip"`
	SenderEmail     string `json:"sender-email"`
	SenderCity      string `json:"sender-city"`
	SenderState     string `json:"sender-state"`
	SenderCountry   string `json:"sender-country"`
	ReceiverName    string `json:"receiver-name"`
	ReceiverStreet  string `json:"receiver-street1"`
	ReceiverPhone   string `json:"receiver-phone"`
	ReceiverZip     string `json:"receiver-zip"`
	ReceiverEmail   string `json:"receiver-email"`
	ReceiverCity    string `json:"receiver-city"`
	ReceiverState   string `json:"receiver-state"`
	ReceiverCountry string `json:"receiver-country"`
	Length          string `json:"length"`
	Weight          string `json:"weight"`
	Width           string `json:"width"`
	Height          string `json:"height"`
	DistanceUnit    string `json:"distance-unit,omitempty"`
	MassUnit        string `json:"mass-unit,omitempty"`
}

type CreateNewShipmentOperation struct {
	options *sdk.OperationInfo
}

func NewCreateNewShipmentOperation() *CreateNewShipmentOperation {
	return &CreateNewShipmentOperation{
		options: &sdk.OperationInfo{
			Name:        "Create A New Shipment",
			Description: "creates a new shipment",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"sender-name": autoform.NewShortTextField().
					SetDisplayName("Sender's Name").
					SetDescription("").
					SetRequired(true).
					Build(),
				"sender-street1": autoform.NewLongTextField().
					SetDisplayName(" Sender's Street 1").
					SetDescription("street ").
					SetRequired(true).
					Build(),
				"sender-email": autoform.NewShortTextField().
					SetDisplayName("Sender's Email").
					SetDescription("Sender's email").
					SetRequired(true).
					Build(),
				"sender-city": autoform.NewShortTextField().
					SetDisplayName("Sender's City").
					SetDescription("city").
					SetRequired(true).
					Build(),
				"sender-country": getCountriesInput(),
				"sender-state": autoform.NewShortTextField().
					SetDisplayName(" Sender's State").
					SetDescription("Sender's state").
					SetRequired(true).
					Build(),
				"sender-zip": autoform.NewShortTextField().
					SetDisplayName("Sender's City").
					SetDescription("Sender's city").
					SetRequired(true).
					Build(),
				"sender-phone": autoform.NewShortTextField().
					SetDisplayName("Sender's Phone").
					SetDescription("Sender's phone number").
					SetRequired(true).
					Build(),
				"receiver-name": autoform.NewShortTextField().
					SetDisplayName("Receiver's Name").
					SetDescription("receiver name").
					SetRequired(true).
					Build(),
				"receiver-street1": autoform.NewLongTextField().
					SetDisplayName(" Receiver's Street 1").
					SetDescription("receiver's street ").
					SetRequired(true).
					Build(),
				"receiver-email": autoform.NewShortTextField().
					SetDisplayName("Receiver's Email").
					SetDescription("receiver's email").
					SetRequired(true).
					Build(),
				"receiver-city": autoform.NewShortTextField().
					SetDisplayName("Receiver's City").
					SetDescription("receiver's city").
					SetRequired(true).
					Build(),
				"receiver-country": getCountriesInput(),
				"receiver-state": autoform.NewShortTextField().
					SetDisplayName(" Receiver's State").
					SetDescription("Receiver's state").
					SetRequired(true).
					Build(),
				"receiver-zip": autoform.NewShortTextField().
					SetDisplayName("Receiver's City").
					SetDescription("Receiver's city").
					SetRequired(true).
					Build(),
				"receiver-phone": autoform.NewShortTextField().
					SetDisplayName("Receiver's Phone").
					SetDescription("Receiver's phone number").
					SetRequired(true).
					Build(),
				"length": autoform.NewShortTextField().
					SetDisplayName("Parcel length").
					SetDescription("parcel length").
					SetRequired(true).
					Build(),
				"width": autoform.NewShortTextField().
					SetDisplayName("Parcel width").
					SetDescription("parcel width").
					SetRequired(true).
					Build(),
				"height": autoform.NewShortTextField().
					SetDisplayName("Parcel height").
					SetDescription("parcel height").
					SetRequired(true).
					Build(),
				"distance-unit": autoform.NewSelectField().
					SetDisplayName("Distance Unit").
					SetDescription("distance unit").
					SetOptions(distanceUnit).
					SetRequired(true).
					Build(),
				"weight": autoform.NewShortTextField().
					SetDisplayName("Parcel weight").
					SetDescription("parcel weight").
					SetRequired(true).
					Build(),
				"mass-unit": autoform.NewSelectField().
					SetDisplayName("Mass Unit").
					SetDescription("mass unit").
					SetOptions(massUnit).
					SetRequired(true).
					Build(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *CreateNewShipmentOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing shippo api key")
	}

	endpoint := "/shipments"

	input := sdk.InputToType[createNewShipmentOperationProps](ctx)

	shipmentData := map[string]interface{}{
		"address_to": map[string]interface{}{
			"name":    input.SenderName,
			"street1": input.SenderStreet,
			"city":    input.SenderCity,
			"state":   input.SenderState,
			"zip":     input.SenderZip,
			"country": input.SenderCountry,
			"phone":   input.SenderPhone,
			"email":   input.SenderEmail,
		},
		"address_from": map[string]interface{}{
			"name":    input.ReceiverName,
			"street1": input.ReceiverStreet,
			"city":    input.ReceiverCity,
			"state":   input.ReceiverState,
			"zip":     input.ReceiverZip,
			"country": input.ReceiverCountry,
			"phone":   input.ReceiverPhone,
			"email":   input.ReceiverEmail,
		},
		"parcels": []map[string]interface{}{
			{
				"length":        input.Length,
				"width":         input.Width,
				"height":        input.Height,
				"distance_unit": input.DistanceUnit,
				"weight":        input.Weight,
				"mass_unit":     input.MassUnit,
			},
		},
	}

	response, err := CreateAShipment(endpoint, ctx.Auth.Extra["api-key"], shipmentData)
	if err != nil {
		return nil, fmt.Errorf("error creating shipment:  %v", err)
	}

	return response, nil
}

func (c *CreateNewShipmentOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *CreateNewShipmentOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
