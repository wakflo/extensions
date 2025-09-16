package actions

import (
	"errors"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/pinterest/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type updatePinActionProps struct {
	BoardID     string `json:"board_id"`
	PinID       string `json:"pin_id"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Link        string `json:"link,omitempty"`
	AltText     string `json:"alt_text,omitempty"`
	Note        string `json:"note,omitempty"`
}

type UpdatePinAction struct{}

// Metadata returns metadata about the action
func (a *UpdatePinAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "update_pin",
		DisplayName:   "Update Pin",
		Description:   "Update an existing pin's details",
		Type:          core.ActionTypeAction,
		Documentation: updatePinDocs,
		Icon:          "mdi:pin-edit",
		SampleOutput:  "",
		Settings:      core.ActionSettings{},
	}
}

func (a *UpdatePinAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("update_pin", "Update Pin")
	form.Description("Update the details of an existing pin")

	// Create validation builder
	v := smartform.NewValidationBuilder()

	// Board selection
	shared.RegisterBoardsProps(form)

	// Pin selection (dependent on board)
	shared.RegisterBoardPinsProps(form)

	// Section for pin details to update
	form.SectionField("pin_details", "Pin Details")

	// Title field
	form.TextField("title", "Title").
		Required(false).
		Placeholder("Enter pin title").
		AddValidation(v.MaxLength(100, "Title must be 100 characters or less"))

	// Description field
	form.TextareaField("description", "Description").
		Required(false).
		Placeholder("Enter pin description").
		AddValidation(v.MaxLength(800, "Description must be 800 characters or less"))

	// Link field
	form.TextField("link", "Link").
		Required(false).
		Placeholder("https://example.com").
		AddValidation(v.Pattern("^https?://.*", "Link must start with http:// or https://"))

	// Alt text field for accessibility
	form.TextField("alt_text", "Alt Text").
		Required(false).
		Placeholder("Alternative text for the pin image").
		AddValidation(v.MaxLength(500, "Alt text must be 500 characters or less"))

	// Note field
	form.TextareaField("note", "Note").
		Required(false).
		Placeholder("Add a note to this pin").
		AddValidation(v.MaxLength(500, "Note must be 500 characters or less"))

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *UpdatePinAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *UpdatePinAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[updatePinActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	if authCtx.AccessToken == "" {
		return nil, errors.New("missing Pinterest auth token")
	}
	accessToken := authCtx.Token.AccessToken

	if input.BoardID == "" {
		return nil, errors.New("board ID is required")
	}

	if input.PinID == "" {
		return nil, errors.New("pin ID is required")
	}

	// Build update data map with only the fields that have values
	updateData := make(map[string]interface{})

	if input.Title != "" {
		updateData["title"] = input.Title
	}

	if input.Description != "" {
		updateData["description"] = input.Description
	}

	if input.Link != "" {
		updateData["link"] = input.Link
	}

	if input.AltText != "" {
		updateData["alt_text"] = input.AltText
	}

	if input.Note != "" {
		updateData["note"] = input.Note
	}

	// Check if at least one field is being updated
	if len(updateData) == 0 {
		return nil, errors.New("at least one field must be provided for update")
	}

	// Call the shared function to update the pin
	updatedPin, err := shared.UpdatePin(accessToken, input.PinID, updateData)
	if err != nil {
		return nil, err
	}

	return updatedPin, nil
}

func NewUpdatePinAction() sdk.Action {
	return &UpdatePinAction{}
}
