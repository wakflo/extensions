package actions

import (
	"errors"
	"fmt"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/mailchimp/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type addNoteToSubscriberActionProps struct {
	ListID string `json:"list-id"`
	Email  string `json:"email"`
	Note   string `json:"note"`
}

type AddNoteToSubscriberAction struct{}

func (a *AddNoteToSubscriberAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		DisplayName:   "Add Note To Subscriber",
		Description:   "Adds a note to a subscriber's record, allowing you to store additional information and context about the subscriber. This integration action is useful for tracking important details or updates related to a specific subscriber, making it easier to manage and analyze their activity over time.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: addNoteToSubscriberDocs,
		SampleOutput: map[string]any{
			"success": "Note added!",
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *AddNoteToSubscriberAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("add_note_to_subscriber", "Add Note To Subscriber")

	form.TextField("list-id", "List ID").
		Placeholder("Enter a value for List ID.").
		Required(true).
		HelpText("The ID of the list to add the contact to.")

	form.TextField("note", "Note").
		Placeholder("Enter a value for Note.").
		Required(true).
		HelpText("The note to add to the subscriber.")

	form.TextField("email", "Email").
		Placeholder("Enter a value for Email.").
		Required(true).
		HelpText("The email of the subscriber.")

	schema := form.Build()

	return schema
}

func (a *AddNoteToSubscriberAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[addNoteToSubscriberActionProps](ctx)
	if err != nil {
		return nil, err
	}

	tokenSource := ctx.Auth().Token
	if tokenSource == nil {
		return nil, errors.New("missing authentication token")
	}
	token := tokenSource.AccessToken

	dc, err := shared.GetMailChimpServerPrefix(token)
	if err != nil {
		return nil, fmt.Errorf("unable to get mailchimp server prefix: %v", err)
	}

	err = shared.AddMemberNote(token, dc, input.ListID, input.Email, input.Note)
	if err != nil {
		return nil, err
	}

	return sdkcore.JSON(map[string]interface{}{
		"result": "note added!",
	}), err
}

func (a *AddNoteToSubscriberAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewAddNoteToSubscriberAction() sdk.Action {
	return &AddNoteToSubscriberAction{}
}
