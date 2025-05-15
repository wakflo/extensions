package actions

import (
	"errors"
	"fmt"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/shippo/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type createNewShipmentActionProps struct {
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

type CreateNewShipmentAction struct{}

// Metadata returns metadata about the action
func (a *CreateNewShipmentAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_new_shipment",
		DisplayName:   "Create A New Shipment",
		Description:   "creates a new shipment",
		Type:          core.ActionTypeAction,
		Documentation: createNewShipmentDocs,
		SampleOutput: map[string]any{
			"shipment": map[string]any{},
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *CreateNewShipmentAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_new_shipment", "Create A New Shipment")

	form.TextField("sender-name", "Sender's Name").
		Required(true)

	form.TextField("sender-street1", "Sender's Street 1").
		Required(true).
		HelpText("street")

	form.TextField("sender-email", "Sender's Email").
		Required(true).
		HelpText("Sender's email")

	form.TextField("sender-city", "Sender's City").
		Required(true).
		HelpText("city")

	form.TextField("sender-state", "Sender's State").
		Required(true).
		HelpText("Sender's state")

	form.TextField("sender-zip", "Sender's Zip").
		Required(true).
		HelpText("Sender's zip")

	form.TextField("sender-phone", "Sender's Phone").
		Required(true).
		HelpText("Sender's phone number")

	form.TextField("receiver-name", "Receiver's Name").
		Required(true).
		HelpText("receiver name")

	form.TextField("receiver-street1", "Receiver's Street 1").
		Required(true).
		HelpText("receiver's street")

	form.TextField("receiver-email", "Receiver's Email").
		Required(true).
		HelpText("receiver's email")

	form.TextField("receiver-city", "Receiver's City").
		Required(true).
		HelpText("receiver's city")

	form.TextField("receiver-state", "Receiver's State").
		Required(true).
		HelpText("Receiver's state")

	form.TextField("receiver-zip", "Receiver's Zip").
		Required(true).
		HelpText("Receiver's zip")

	form.TextField("receiver-phone", "Receiver's Phone").
		Required(true).
		HelpText("Receiver's phone number")

	form.TextField("length", "Parcel length").
		Required(true).
		HelpText("parcel length")

	form.TextField("width", "Parcel width").
		Required(true).
		HelpText("parcel width")

	form.TextField("height", "Parcel height").
		Required(true).
		HelpText("parcel height")

	form.SelectField("distance-unit", "Distance Unit").
		Required(true).
		AddOption("in", "Inches").
		AddOption("cm", "Centimeters").
		AddOption("mm", "Millimeters").
		AddOption("yd", "Yards").
		AddOption("m", "Meters").
		AddOption("ft", "Feet").
		HelpText("distance unit")

	form.TextField("weight", "Parcel weight").
		Required(true).
		HelpText("parcel weight")

	form.SelectField("mass-unit", "Mass Unit").
		Required(true).
		AddOption("kg", "Kilogram").
		AddOption("lb", "Pound").
		AddOption("oz", "Ounce").
		AddOption("g", "Gram").
		HelpText("mass unit")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *CreateNewShipmentAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *CreateNewShipmentAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	// Get the auth context and API key
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	if authCtx.Extra["api-key"] == "" {
		return nil, errors.New("missing shippo api key")
	}

	endpoint := "/shipments"

	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[createNewShipmentActionProps](ctx)
	if err != nil {
		return nil, err
	}

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

	response, err := shared.CreateAShipment(endpoint, authCtx.Extra["api-key"], shipmentData)
	if err != nil {
		return nil, fmt.Errorf("error creating shipment: %v", err)
	}

	return response, nil
}

func NewCreateNewShipmentAction() sdk.Action {
	return &CreateNewShipmentAction{}
}
