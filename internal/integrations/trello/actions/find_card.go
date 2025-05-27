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

type findCardActionProps struct {
	BoardID string `json:"boards"`
	ListID  string `json:"list"`
	CardID  string `json:"cardId"`
}

type FindCardAction struct{}

func (a *FindCardAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "find_card",
		DisplayName:   "Find Card",
		Description:   "Searches for a specific card within a specified deck or collection, returning the card's details if found.",
		Type:          core.ActionTypeAction,
		Documentation: findCardDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *FindCardAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("find_card", "Find Card")

	shared.RegisterBoardsProp(form)

	shared.RegisterBoardListsProp(form)

	shared.RegisterCardsProp(form)

	schema := form.Build()
	return schema
}

func (a *FindCardAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[findCardActionProps](ctx)
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

	response, err := shared.TrelloRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}

	return response, nil
}

func (a *FindCardAction) Auth() *core.AuthMetadata {
	return nil
}

func NewFindCardAction() sdk.Action {
	return &FindCardAction{}
}
