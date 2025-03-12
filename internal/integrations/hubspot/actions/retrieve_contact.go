package actions

import (
	"encoding/json"
	"net/http"

	"github.com/wakflo/extensions/internal/integrations/hubspot/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type retrieveContactActionProps struct {
	Email string `json:"email"`
}

type RetrieveContactAction struct{}

func (a *RetrieveContactAction) Name() string {
	return "Retrieve Contact"
}

func (a *RetrieveContactAction) Description() string {
	return "retrieve a specific Contact"
}

func (a *RetrieveContactAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *RetrieveContactAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &retrieveContactDocs,
	}
}

func (a *RetrieveContactAction) Icon() *string {
	return nil
}

func (a *RetrieveContactAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"email": autoform.NewShortTextField().
			SetDisplayName("Contact's Email").
			SetDescription("contact's email").
			SetRequired(true).Build(),
	}
}

func (a *RetrieveContactAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[retrieveContactActionProps](ctx.BaseContext)
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

	resp, err := shared.HubspotClient(reqURL, ctx.Auth.AccessToken, http.MethodPost, requestBody)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (a *RetrieveContactAction) Auth() *sdk.Auth {
	return nil
}

func (a *RetrieveContactAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"results": []map[string]any{
			{
				"id": "51",
				"properties": map[string]any{
					"firstname": "John",
					"lastname":  "Doe",
					"email":     "john.doe@example.com",
				},
			},
		},
		"paging": map[string]any{
			"next": map[string]any{
				"after": "NTI=",
				"link":  "https://api.hubapi.com/crm/v3/objects/contacts?after=NTI=",
			},
		},
	}
}

func (a *RetrieveContactAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewRetrieveContactAction() sdk.Action {
	return &RetrieveContactAction{}
}
