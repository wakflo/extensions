package actions

import (
	"fmt"
	"net/url"

	"github.com/wakflo/extensions/internal/integrations/activecampaign/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type listContactsActionProps struct {
	Email  string `json:"email"`
	ListID string `json:"list-id"`
	TagID  string `json:"tag-id"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

type ListContactsAction struct{}

func (a *ListContactsAction) Name() string {
	return "List Contacts"
}

func (a *ListContactsAction) Description() string {
	return "Retrieve a list of contacts from your ActiveCampaign account with filtering options."
}

func (a *ListContactsAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *ListContactsAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &listContactsDocs,
	}
}

func (a *ListContactsAction) Icon() *string {
	icon := "mdi:account-group"
	return &icon
}

func (a *ListContactsAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"email": autoform.NewShortTextField().
			SetDisplayName("Email").
			SetDescription("Filter contacts by email address").
			Build(),
		"list_id": shared.GetActiveCampaignListsInput(),
		"tag-id": autoform.NewShortTextField().
			SetDisplayName("Tag ID").
			SetDescription("Filter contacts by a specific tag ID").
			Build(),
		"limit": autoform.NewNumberField().
			SetDisplayName("Limit").
			SetDescription("Maximum number of contacts to return (default: 100)").
			Build(),
	}
}

func (a *ListContactsAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[listContactsActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	endpoint := "contacts?"

	if input.Email != "" {
		endpoint += fmt.Sprintf("filters[email]=%s&", url.QueryEscape(input.Email))
	}

	if input.ListID != "" {
		endpoint += fmt.Sprintf("filters[listid]=%s&", input.ListID)
	}

	if input.TagID != "" {
		endpoint += fmt.Sprintf("filters[tagid]=%s&", input.TagID)
	}

	limit := 20
	maxLim := 100
	if input.Limit > 0 {
		if input.Limit > maxLim {
			limit = 100
		} else {
			limit = input.Limit
		}
	}

	offset := 0
	if input.Offset > 0 {
		offset = input.Offset
	}

	endpoint += fmt.Sprintf("limit=%d&offset=%d", limit, offset)

	response, err := shared.GetActiveCampaignClient(
		ctx.Auth.Extra["api_url"],
		ctx.Auth.Extra["api_key"],
		endpoint,
	)
	if err != nil {
		return nil, err
	}

	responseMap, ok := response.(map[string]interface{})
	if !ok {
		return nil, err
	}

	contacts, ok := responseMap["contacts"]
	if !ok {
		return nil, err
	}

	contactsArray, ok := contacts.([]interface{})
	if !ok {
		return nil, err
	}

	return contactsArray, nil
}

func (a *ListContactsAction) Auth() *sdk.Auth {
	return nil
}

func (a *ListContactsAction) SampleData() sdkcore.JSON {
	return []map[string]any{
		{
			"id":        "123",
			"email":     "sample1@example.com",
			"firstName": "John",
			"lastName":  "Doe",
			"phone":     "+1234567890",
		},
		{
			"id":        "124",
			"email":     "sample2@example.com",
			"firstName": "Jane",
			"lastName":  "Smith",
			"phone":     "+0987654321",
		},
	}
}

func (a *ListContactsAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewListContactsAction() sdk.Action {
	return &ListContactsAction{}
}
