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

type createListActionProps struct {
	IDBoard  string `json:"idBoard"`
	Name     string `json:"name"`
	Position string `json:"position"`
}

type CreateListAction struct{}

func (a *CreateListAction) Name() string {
	return "Create List"
}

func (a *CreateListAction) Description() string {
	return "Create List: Automatically generates a list of items based on predefined criteria or data sources, allowing you to quickly and easily collect and organize information."
}

func (a *CreateListAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *CreateListAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createListDocs,
	}
}

func (a *CreateListAction) Icon() *string {
	return nil
}

func (a *CreateListAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"idBoard": shared.GetBoardsInput(),
		"name": autoform.NewShortTextField().
			SetDisplayName("List Name").
			SetDescription("The name of the list to create").
			SetRequired(true).
			Build(),
	}
}

func (a *CreateListAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createListActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	if ctx.Auth.Extra["api-key"] == "" {
		return nil, errors.New("missing trello api credentials")
	}
	apiKey := ctx.Auth.Extra["api-key"]
	apiToken := ctx.Auth.Extra["api-token"]

	fullURL := fmt.Sprintf("%s/lists?name=%s&idBoard=%s&key=%s&token=%s", shared.BaseURL, input.Name, input.IDBoard, apiKey, apiToken)

	response, err := shared.TrelloRequest(http.MethodPost, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}

	return response, nil
}

func (a *CreateListAction) Auth() *sdk.Auth {
	return nil
}

func (a *CreateListAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *CreateListAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreateListAction() sdk.Action {
	return &CreateListAction{}
}
