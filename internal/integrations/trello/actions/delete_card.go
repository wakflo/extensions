package actions

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/trello/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type deleteCardActionProps struct {
	CardID string `json:"cardId"`
}

type DeleteCardAction struct{}

func (a *DeleteCardAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "delete_card",
		DisplayName:   "Delete Card",
		Description:   "Deletes a card from a specified board or list in Trello, allowing you to automate the removal of cards that meet specific criteria.",
		Type:          core.ActionTypeAction,
		Documentation: deleteCardDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *DeleteCardAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("delete_card", "Delete Card")

	form.TextField("cardId", "Card ID").
		Placeholder("The id of the card to delete").
		Required(true).
		HelpText("The id of the card to delete")

	schema := form.Build()
	return schema
}

func (a *DeleteCardAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[deleteCardActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	authExtra := authCtx.Extra
	if authExtra["api-key"] == "" {
		return nil, errors.New("missing trello api credentials")
	}
	apiKey := authExtra["api-key"]
	apiToken := authExtra["api-token"]

	fullURL := fmt.Sprintf("%s/cards/%s?key=%s&token=%s", shared.BaseURL, input.CardID, apiKey, apiToken)

	_, err = shared.TrelloRequest(http.MethodDelete, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}

	return map[string]interface{}{
		"Result": "Card deleted Successfully",
	}, nil
}

func (a *DeleteCardAction) Auth() *core.AuthMetadata {
	return nil
}

func NewDeleteCardAction() sdk.Action {
	return &DeleteCardAction{}
}
