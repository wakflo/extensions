package actions

import (
	"fmt"
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/hubspot/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type searchOwnerByEmailActionProps struct {
	Email string `json:"email"`
}

type SearchOwnerByEmailAction struct{}

func (a *SearchOwnerByEmailAction) Name() string {
	return "Search Owner by Email"
}

func (a *SearchOwnerByEmailAction) Description() string {
	return "Find a HubSpot owner by their email address"
}

func (a *SearchOwnerByEmailAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *SearchOwnerByEmailAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &searchOwnerDocs,
	}
}

func (a *SearchOwnerByEmailAction) Icon() *string {
	return nil
}

func (a *SearchOwnerByEmailAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"email": autoform.NewShortTextField().
			SetDisplayName("Owner's Email").
			SetDescription("Email address of the HubSpot owner").
			SetRequired(true).Build(),
	}
}

func (a *SearchOwnerByEmailAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[searchOwnerByEmailActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}
	reqURL := "/crm/v3/owners"

	resp, err := shared.HubspotClient(reqURL, ctx.Auth.AccessToken, http.MethodGet, nil)
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

func (a *SearchOwnerByEmailAction) Auth() *sdk.Auth {
	return nil
}

func (a *SearchOwnerByEmailAction) SampleData() sdkcore.JSON {
	return map[string]interface{}{
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
	}
}

func (a *SearchOwnerByEmailAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewSearchOwnerByEmailAction() sdk.Action {
	return &SearchOwnerByEmailAction{}
}
