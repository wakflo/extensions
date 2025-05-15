package actions

import (
	"fmt"
	"net/http"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/hubspot/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type listContactsActionProps struct {
	Limit int `json:"limit"`
}

type ListContactsAction struct{}

// Metadata returns metadata about the action
func (a *ListContactsAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "list_contacts",
		DisplayName:   "List Contacts",
		Description:   "List all contacts",
		Type:          core.ActionTypeAction,
		Documentation: listContactsDocs,
		SampleOutput: map[string]any{
			"results": []map[string]any{
				{
					"id": "51",
					"properties": map[string]any{
						"firstname": "John",
						"lastname":  "Doe",
						"email":     "john.doe@example.com",
					},
				},
				{
					"id": "52",
					"properties": map[string]any{
						"firstname": "Jane",
						"lastname":  "Smith",
						"email":     "jane.smith@example.com",
					},
				},
			},
			"paging": map[string]any{
				"next": map[string]any{
					"after": "NTI=",
					"link":  "https://api.hubapi.com/crm/v3/objects/contacts?after=NTI=",
				},
			},
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *ListContactsAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("list_contacts", "List Contacts")

	form.NumberField("limit", "Limit").
		Required(false).
		HelpText("Limit")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *ListContactsAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *ListContactsAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[listContactsActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	if input.Limit <= 0 {
		input.Limit = 20
	}

	url := fmt.Sprintf("/crm/v3/objects/contacts?limit=%d", input.Limit)

	resp, err := shared.HubspotClient(url, authCtx.Token.AccessToken, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NewListContactsAction() sdk.Action {
	return &ListContactsAction{}
}
