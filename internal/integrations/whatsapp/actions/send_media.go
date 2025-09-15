package actions

import (
	"fmt"
	"strings"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/whatsapp/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type sendMediaActionProps struct {
	PhoneNumberID string `json:"phone_number_id,omitempty"`
	ToNumber      string `json:"to_number,omitempty"`
	MediaType     string `json:"media_type,omitempty"`
	MediaURL      string `json:"media_url,omitempty"`
	MediaID       string `json:"media_id,omitempty"`
	Caption       string `json:"caption,omitempty"`
	Filename      string `json:"filename,omitempty"`
}

type SendMediaAction struct{}

func (a *SendMediaAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "send_media",
		DisplayName:   "Send Media",
		Description:   "Send media files (images, videos, documents, audio) to a WhatsApp number. The recipient must have previously opted in to receive messages from your business or messaged you in the last 24 hours.",
		Type:          core.ActionTypeAction,
		Documentation: sendMediaDocs,
		SampleOutput: map[string]any{
			"messaging_product": "whatsapp",
			"contacts": []map[string]interface{}{
				{
					"input": "+1234567890",
					"wa_id": "1234567890",
				},
			},
			"messages": []map[string]interface{}{
				{
					"id": "wamid.HBgMMTIzNDU2Nzg5MBUCABIYFjNFQjBEN0EwNTAwNEUwRkI0MzU4QTkA",
				},
			},
		},
		Settings: core.ActionSettings{},
	}
}

func (a *SendMediaAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("send_media", "Send Media")

	form.TextField("to_number", "Recipient's Phone Number").
		Placeholder("Enter recipient's phone number").
		Required(true).
		HelpText("The recipient's phone number in international format (e.g., +1XXXXXXXXXX).")

	form.SelectField("media_type", "Media Type").
		Placeholder("Select media type").
		Required(true).
		AddOptions([]*smartform.Option{
			{Value: "image", Label: "Image"},
			{Value: "video", Label: "Video"},
			{Value: "audio", Label: "audio"},
			{Value: "document", Label: "Document"},
			{Value: "sticker", Label: "Sticker"},
		}...).
		HelpText("The type of media you want to send.")

	form.TextField("media_url", "Media URL").
		Placeholder("https://example.com/media.jpg").
		Required(true).
		HelpText("Public URL of the media file. Either Media URL or Media ID is required.")

	form.TextareaField("caption", "Caption").
		Placeholder("Enter caption for the media").
		Required(false).
		HelpText("Caption for the media (not applicable for audio or sticker).")

	form.TextField("filename", "Filename").
		Placeholder("document.pdf").
		Required(false).
		HelpText("Filename for document type media. Required when sending documents.")

	schema := form.Build()

	return schema
}

func (a *SendMediaAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[sendMediaActionProps](ctx)
	if err != nil {
		return nil, err
	}

	if input.MediaURL == "" {
		return nil, fmt.Errorf("Media Url must be provided")
	}

	toNumber := input.ToNumber
	if !strings.HasPrefix(toNumber, "+") {
		toNumber = "+" + toNumber
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	client := shared.NewWhatsAppClient(authCtx.Extra["token"])

	mediaObject := make(map[string]interface{})
	if input.MediaURL != "" {
		mediaObject["link"] = input.MediaURL
	}

	if input.Caption != "" && (input.MediaType == "image" || input.MediaType == "video" || input.MediaType == "document") {
		mediaObject["caption"] = input.Caption
	}

	if input.MediaType == "document" && input.Filename != "" {
		mediaObject["filename"] = input.Filename
	}

	messageBody := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                toNumber,
		"type":              input.MediaType,
		input.MediaType:     mediaObject,
	}

	endpoint := fmt.Sprintf("%s/messages", authCtx.Extra["phone-id"])

	response, err := client.SendRequest("POST", endpoint, messageBody)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (a *SendMediaAction) Auth() *core.AuthMetadata {
	return nil
}

func NewSendMediaAction() sdk.Action {
	return &SendMediaAction{}
}
