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

type createShipmentLabelActionProps struct {
	Rate string `json:"rate"`
}

type CreateShipmentLabelAction struct{}

// Metadata returns metadata about the action
func (a *CreateShipmentLabelAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_shipment_label",
		DisplayName:   "Create Shipment Label",
		Description:   "Create a shipping label based on a rate.",
		Type:          core.ActionTypeAction,
		Documentation: newShipmentLabelDocs,
		SampleOutput: map[string]any{
			"transaction": map[string]any{},
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *CreateShipmentLabelAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_shipment_label", "Create Shipment Label")

	form.TextField("rate", "Rate").
		Required(true).
		HelpText("The Shippo rate object ID.")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *CreateShipmentLabelAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *CreateShipmentLabelAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	if authCtx.Extra["api-key"] == "" {
		return nil, errors.New("missing shippo api key")
	}

	// Use the InputToTypeSafely helper function to convert the input to our struct
	input, err := sdk.InputToTypeSafely[createShipmentLabelActionProps](ctx)
	if err != nil {
		return nil, err
	}

	endpoint := "/transactions"

	shipmentData := map[string]interface{}{
		"rate":            input.Rate,
		"async":           false,
		"label_file_type": "PDF",
	}

	response, err := shared.CreateAShipment(endpoint, authCtx.Extra["api-key"], shipmentData)
	if err != nil {
		return nil, fmt.Errorf("error creating shipment: %v", err)
	}

	return response, nil
}

func NewCreateShipmentLabelAction() sdk.Action {
	return &CreateShipmentLabelAction{}
}
