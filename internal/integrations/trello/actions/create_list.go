package actions

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/trello/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type createListActionProps struct {
	IDBoard string `json:"boards"`
	Name    string `json:"name"`
}

type CreateListAction struct{}

func (a *CreateListAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_list",
		DisplayName:   "Create List",
		Description:   "Create List: Automatically generates a list of items based on predefined criteria or data sources, allowing you to quickly and easily collect and organize information.",
		Type:          core.ActionTypeAction,
		Documentation: createListDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *CreateListAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_list", "Create List")
	shared.RegisterBoardsProp(form)

	form.TextField("name", "List Name").
		Placeholder("The name of the list to create").
		Required(true).
		HelpText("The name of the list to create")

	schema := form.Build()
	return schema
}

func (a *CreateListAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[createListActionProps](ctx)
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

	// Build query parameters
	params := url.Values{}
	params.Add("name", input.Name)
	params.Add("idBoard", input.IDBoard)
	params.Add("key", apiKey)
	params.Add("token", apiToken)

	fullURL := fmt.Sprintf("%s/lists?%s", shared.BaseURL, params.Encode())

	response, err := shared.TrelloRequest(http.MethodPost, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}

	return response, nil
}

func (a *CreateListAction) Auth() *core.AuthMetadata {
	return nil
}

func NewCreateListAction() sdk.Action {
	return &CreateListAction{}
}
