package actions

import (
	"errors"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/pinterest/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type createPinActionProps struct {
	BoardID         string `json:"board_id"`
	Title           string `json:"title,omitempty"`
	Description     string `json:"description,omitempty"`
	Link            string `json:"link,omitempty"`
	AltText         string `json:"alt_text,omitempty"`
	Note            string `json:"note,omitempty"`
	MediaSourceType string `json:"media_source_type"`
	MediaSourceURL  string `json:"media_source_url,omitempty"`
	MediaSourceID   string `json:"media_source_id,omitempty"`
	DominantColor   string `json:"dominant_color,omitempty"`
}

type CreatePinAction struct{}

// Metadata returns metadata about the action
func (a *CreatePinAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_pin",
		DisplayName:   "Create Pin",
		Description:   "Create a new pin on Pinterest",
		Type:          core.ActionTypeAction,
		Documentation: createPinDocs,
		Icon:          "mdi:pin-plus",
		SampleOutput:  "",
		Settings:      core.ActionSettings{},
	}
}

func (a *CreatePinAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_pin", "Create Pin")
	form.Description("Create a new pin on Pinterest with image or video content")

	// Create validation builder
	v := smartform.NewValidationBuilder()

	// Board selection
	shared.RegisterBoardsProps(form)

	// Pin Details Section
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

	// Media Source Section
	form.SectionField("media_source", "Media Source")

	// Media source type
	form.SelectField("media_source_type", "Media Source Type").
		Required(true).
		AddOption("image_url", "Image URL").
		AddOption("video_id", "Video ID").
		AddOption("image_base64", "Image Base64")

	// Media source URL (for image_url type)
	form.TextField("media_source_url", "Media Source URL").
		Required(false).
		Placeholder("https://example.com/image.jpg").
		VisibleWhenEquals("media_source_type", "image_url").
		AddValidation(v.Pattern("^https?://.*", "URL must start with http:// or https://"))

	// Media source ID (for video_id type)
	form.TextField("media_source_id", "Media Source ID").
		Required(false).
		Placeholder("Enter video ID").
		VisibleWhenEquals("media_source_type", "video_id")

	// Optional dominant color
	form.TextField("dominant_color", "Dominant Color").
		Required(false).
		Placeholder("#FF5733").
		AddValidation(v.Pattern("^#[0-9A-Fa-f]{6}$", "Color must be in hex format (e.g., #FF5733)"))

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *CreatePinAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *CreatePinAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[createPinActionProps](ctx)
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

	if input.MediaSourceType == "" {
		return nil, errors.New("media source type is required")
	}

	// Build pin data
	pinData := make(map[string]interface{})

	// Required fields
	pinData["board_id"] = input.BoardID

	// Build media source based on type
	mediaSource := make(map[string]interface{})
	mediaSource["source_type"] = input.MediaSourceType

	switch input.MediaSourceType {
	case "image_url":
		if input.MediaSourceURL == "" {
			return nil, errors.New("media source URL is required for image URL type")
		}
		mediaSource["url"] = input.MediaSourceURL
	case "video_id":
		if input.MediaSourceID == "" {
			return nil, errors.New("media source ID is required for video ID type")
		}
		mediaSource["media_id"] = input.MediaSourceID
	case "image_base64":
		// Base64 handling would require additional logic
		// For now, we'll just note it's not implemented
		return nil, errors.New("image base64 upload is not yet implemented")
	default:
		return nil, errors.New("invalid media source type")
	}

	pinData["media_source"] = mediaSource

	// Optional fields
	if input.Title != "" {
		pinData["title"] = input.Title
	}

	if input.Description != "" {
		pinData["description"] = input.Description
	}

	if input.Link != "" {
		pinData["link"] = input.Link
	}

	if input.AltText != "" {
		pinData["alt_text"] = input.AltText
	}

	if input.Note != "" {
		pinData["note"] = input.Note
	}

	if input.DominantColor != "" {
		pinData["dominant_color"] = input.DominantColor
	}

	// Call the shared function to create the pin
	createdPin, err := shared.CreatePin(accessToken, pinData)
	if err != nil {
		return nil, err
	}

	return createdPin, nil
}

func NewCreatePinAction() sdk.Action {
	return &CreatePinAction{}
}
