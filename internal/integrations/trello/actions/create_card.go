package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/trello/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type createCardActionProps struct {
	BoardID     string `json:"boards"`
	ListID      string `json:"list"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Position    string `json:"position"`
}

type CreateCardAction struct{}

func (a *CreateCardAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_card",
		DisplayName:   "Create Card",
		Description:   "Create Card: Automatically generates a new card in your chosen project management tool (e.g., Trello, Asana) with pre-populated fields and custom details, streamlining task creation and organization.",
		Type:          core.ActionTypeAction,
		Documentation: createCardDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *CreateCardAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_card", "Create Card")
	shared.RegisterBoardsProp(form)
	shared.RegisterBoardListsProp(form)

	form.TextField("name", "Card Name").
		Placeholder("The name of the card to create").
		Required(true).
		HelpText("The name of the card to create")

	form.TextareaField("description", "Description").
		Placeholder("The description of the card to create").
		HelpText("The description of the card to create")

	form.SelectField("position", "Position").
		AddOption("top", "Top").
		AddOption("bottom", "Bottom").
		Placeholder("Place the card on top or bottom of the list").
		HelpText("Place the card on top or bottom of the list")

	schema := form.Build()
	return schema
}

func (a *CreateCardAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[createCardActionProps](ctx)
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

	fullURL := fmt.Sprintf("%s/cards?idList=%s&key=%s&token=%s", shared.BaseURL, input.ListID, apiKey, apiToken)

	payload := shared.CardRequest{
		Name: input.Name,
		Desc: input.Description,
		Pos:  input.Position,
	}

	payloadBytes, _ := json.Marshal(payload)

	response, err := shared.TrelloRequest(http.MethodPost, fullURL, payloadBytes)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}

	return response, nil
}

func (a *CreateCardAction) Auth() *core.AuthMetadata {
	return nil
}

func NewCreateCardAction() sdk.Action {
	return &CreateCardAction{}
}
