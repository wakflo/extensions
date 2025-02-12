package actions

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/trello/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type deleteCardActionProps struct {
	CardID string `json:"cardId"`
}

type DeleteCardAction struct{}

func (a *DeleteCardAction) Name() string {
	return "Delete Card"
}

func (a *DeleteCardAction) Description() string {
	return "Deletes a card from a specified board or list in Trello, allowing you to automate the removal of cards that meet specific criteria."
}

func (a *DeleteCardAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *DeleteCardAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &deleteCardDocs,
	}
}

func (a *DeleteCardAction) Icon() *string {
	return nil
}

func (a *DeleteCardAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"cardId": autoform.NewShortTextField().
			SetDisplayName("Card ID").
			SetDescription("The id of the card to delete").
			SetRequired(true).
			Build(),
	}
}

func (a *DeleteCardAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[deleteCardActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing trello api credentials")
	}
	apiKey := ctx.Auth.Extra["api-key"]
	apiToken := ctx.Auth.Extra["api-token"]
	fullURL := fmt.Sprintf("%s/cards/%s?key=%s&token=%s", shared.BaseURL, input.CardID, apiKey, apiToken)

	_, err = shared.TrelloRequest(http.MethodDelete, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}

	return map[string]interface{}{
		"Result": "Card deleted Successfully",
	}, nil
}

func (a *DeleteCardAction) Auth() *sdk.Auth {
	return nil
}

func (a *DeleteCardAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *DeleteCardAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewDeleteCardAction() sdk.Action {
	return &DeleteCardAction{}
}
