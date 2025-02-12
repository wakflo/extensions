package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/trello/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type createCardActionProps struct {
	BoardID     string `json:"board_id"`
	ListID      string `json:"idList"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Position    string `json:"position"`
}

type CreateCardAction struct{}

func (a *CreateCardAction) Name() string {
	return "Create Card"
}

func (a *CreateCardAction) Description() string {
	return "Create Card: Automatically generates a new card in your chosen project management tool (e.g., Trello, Asana) with pre-populated fields and custom details, streamlining task creation and organization."
}

func (a *CreateCardAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *CreateCardAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createCardDocs,
	}
}

func (a *CreateCardAction) Icon() *string {
	return nil
}

func (a *CreateCardAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"board_id": shared.GetBoardsInput(),
		"idList":   shared.GetBoardListsInput(),
		"name": autoform.NewShortTextField().
			SetDisplayName("Card Name").
			SetDescription("The name of the card to create").
			SetRequired(true).
			Build(),
		"description": autoform.NewLongTextField().
			SetDisplayName("Description").
			SetDescription("The description of the card to create").
			SetRequired(false).
			Build(),
		"position": autoform.NewSelectField().
			SetDisplayName("Position").
			SetDescription("Place the card on top or bottom of the list").
			SetOptions(shared.TrelloCardPosition).
			SetRequired(false).
			Build(),
	}
}

func (a *CreateCardAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createCardActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing trello api credentials")
	}
	apiKey := ctx.Auth.Extra["api-key"]
	apiToken := ctx.Auth.Extra["api-token"]

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

func (a *CreateCardAction) Auth() *sdk.Auth {
	return nil
}

func (a *CreateCardAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *CreateCardAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreateCardAction() sdk.Action {
	return &CreateCardAction{}
}
