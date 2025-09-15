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

type sendMessageActionProps struct {
	ToNumber    string `json:"to_number,omitempty"`
	MessageText string `json:"message_text,omitempty"`
}

type SendMessageAction struct{}

func (a *SendMessageAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "send_message",
		DisplayName:   "Send Message",
		Description:   "Send a text message to a WhatsApp number. The recipient must have previously opted in to receive messages from your business or messaged you in the last 24 hours.",
		Type:          core.ActionTypeAction,
		Documentation: sendMessageDocs,
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

func (a *SendMessageAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("send_message", "Send Message")

	form.TextField("to_number", "Recipient's Phone Number").
		Placeholder("Enter recipient's phone number").
		Required(true).
		HelpText("The recipient's phone number in international format (e.g., +1XXXXXXXXXX).")

	form.TextareaField("message_text", "Message Text").
		Placeholder("Enter message text").
		Required(true).
		HelpText("The text message you want to send.")

	schema := form.Build()

	return schema
}

func (a *SendMessageAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[sendMessageActionProps](ctx)
	if err != nil {
		return nil, err
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

	messageBody := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                toNumber,
		"type":              "text",
		"text": map[string]interface{}{
			"body": input.MessageText,
		},
	}

	endpoint := fmt.Sprintf("%s/messages", authCtx.Extra["phone-id"])

	response, err := client.SendRequest("POST", endpoint, messageBody)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (a *SendMessageAction) Auth() *core.AuthMetadata {
	return nil
}

func NewSendMessageAction() sdk.Action {
	return &SendMessageAction{}
}
