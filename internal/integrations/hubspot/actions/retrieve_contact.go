package actions

import (
	"encoding/json"
	"net/http"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/hubspot/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type retrieveContactActionProps struct {
	Email string `json:"email"`
}

type RetrieveContactAction struct{}

// Metadata returns metadata about the action
func (a *RetrieveContactAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "retrieve_contact",
		DisplayName:   "Retrieve Contact",
		Description:   "retrieve a specific Contact",
		Type:          core.ActionTypeAction,
		Documentation: retrieveContactDocs,
		SampleOutput: map[string]any{
			"results": []map[string]any{
				{
					"id": "51",
					"properties": map[string]any{
						"firstname": "John",
						"lastname":  "Doe",
						"email":     "john.doe@example.com",
						"phone":     "+1234567890",
						"company":   "Acme Inc.",
						"website":   "https://example.com",
						"address":   "123 Main St",
						"city":      "New York",
						"state":     "NY",
						"zipcode":   "10001",
						"country":   "USA",
					},
					"createdAt": "2023-01-10T12:00:00Z",
					"updatedAt": "2023-04-15T09:45:00Z",
				},
			},
			"total": 1,
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *RetrieveContactAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("retrieve_contact", "Retrieve Contact")

	form.TextField("email", "Contact's Email").
		Required(true).
		HelpText("contact's email")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *RetrieveContactAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *RetrieveContactAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[retrieveContactActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	// Use the search endpoint instead
	reqURL := "/crm/v3/objects/contacts/search"

	// Create the search request body
	searchRequest := map[string]interface{}{
		"filterGroups": []map[string]interface{}{
			{
				"filters": []map[string]interface{}{
					{
						"propertyName": "email",
						"operator":     "EQ",
						"value":        input.Email,
					},
				},
			},
		},
	}

	// Convert the request to JSON
	requestBody, err := json.Marshal(searchRequest)
	if err != nil {
		return nil, err
	}

	resp, err := shared.HubspotClient(reqURL, authCtx.Token.AccessToken, http.MethodPost, requestBody)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func NewRetrieveContactAction() sdk.Action {
	return &RetrieveContactAction{}
}
