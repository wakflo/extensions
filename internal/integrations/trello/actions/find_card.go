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

type findCardActionProps struct {
	CardID string `json:"cardId"`
}

type FindCardAction struct{}

func (a *FindCardAction) Name() string {
	return "Find Card"
}

func (a *FindCardAction) Description() string {
	return "Searches for a specific card within a specified deck or collection, returning the card's details if found."
}

func (a *FindCardAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *FindCardAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &findCardDocs,
	}
}

func (a *FindCardAction) Icon() *string {
	return nil
}

func (a *FindCardAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"name": autoform.NewShortTextField().
			SetLabel("Name").
			SetRequired(true).
			SetPlaceholder("Your name").
			Build(),
	}
}

func (a *FindCardAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[findCardActionProps](ctx.BaseContext)
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

func (a *FindCardAction) Auth() *sdk.Auth {
	return nil
}

func (a *FindCardAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *FindCardAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewFindCardAction() sdk.Action {
	return &FindCardAction{}
}
