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

type searchOwnerByEmailActionProps struct {
	Email string `json:"email"`
}

type SearchOwnerByEmailAction struct{}

// Metadata returns metadata about the action
func (a *SearchOwnerByEmailAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "search_owner_by_email",
		DisplayName:   "Search Owner by Email",
		Description:   "Find a HubSpot owner by their email address",
		Type:          core.ActionTypeAction,
		Documentation: searchOwnerDocs,
		SampleOutput: map[string]interface{}{
			"found": true,
			"owner": map[string]interface{}{
				"id":        "12345",
				"email":     "owner@example.com",
				"firstName": "John",
				"lastName":  "Smith",
				"userId":    "67890",
				"createdAt": "2023-01-01T12:00:00Z",
				"updatedAt": "2023-01-01T12:00:00Z",
				"archived":  false,
				"teams":     []interface{}{},
			},
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *SearchOwnerByEmailAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("search_owner_by_email", "Search Owner by Email")

	form.TextField("email", "Owner's Email").
		Required(true).
		HelpText("Email address of the HubSpot owner")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *SearchOwnerByEmailAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *SearchOwnerByEmailAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[searchOwnerByEmailActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	reqURL := "/crm/v3/owners"

	resp, err := shared.HubspotClient(reqURL, authCtx.Token.AccessToken, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	// Convert response to a map to work with the results
	responseMap, ok := resp.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	// Extract the results array
	results, ok := responseMap["results"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unable to parse owners results")
	}

	var matchingOwner interface{}
	for _, result := range results {
		owner, ok := result.(map[string]interface{})
		if !ok {
			continue
		}

		ownerEmail, ok := owner["email"].(string)
		if !ok {
			continue
		}

		if ownerEmail == input.Email {
			matchingOwner = owner
			break
		}
	}

	if matchingOwner == nil {
		return map[string]interface{}{
			"found":   false,
			"message": "No owner found with that email address",
		}, nil
	}

	return map[string]interface{}{
		"found": true,
		"owner": matchingOwner,
	}, nil
}

func NewSearchOwnerByEmailAction() sdk.Action {
	return &SearchOwnerByEmailAction{}
}
