package actions

import (
	"errors"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/zohoinventory/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type listItemsActionProps struct {
	OrganizationID string `json:"organization_id"`
}

type ListItemsAction struct{}

func (a *ListItemsAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "list_items",
		DisplayName:   "List Items",
		Description:   "Retrieves a list of items from a specified data source or application, allowing you to collect and process data in your workflow.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: listItemsDocs,
		SampleOutput: map[string]any{
			"message": "hello world",
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *ListItemsAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("list_items", "List Items")

	shared.GetOrganizationsProp(form)

	schema := form.Build()

	return schema
}

func (a *ListItemsAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[listItemsActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the token source from the auth context
	tokenSource := ctx.Auth().Token
	if tokenSource == nil {
		return nil, errors.New("missing authentication token")
	}
	token := tokenSource.AccessToken

	url := "/v1/items?organization_id=" + input.OrganizationID

	items, err := shared.GetZohoClient(token, url)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (a *ListItemsAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewListItemsAction() sdk.Action {
	return &ListItemsAction{}
}
