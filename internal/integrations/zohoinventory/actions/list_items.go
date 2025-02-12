package actions

import (
	"github.com/wakflo/extensions/internal/integrations/zohoinventory/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type listItemsActionProps struct {
	OrganizationID string `json:"organization_id"`
}

type ListItemsAction struct{}

func (a *ListItemsAction) Name() string {
	return "List Items"
}

func (a *ListItemsAction) Description() string {
	return "Retrieves a list of items from a specified data source or application, allowing you to collect and process data in your workflow."
}

func (a *ListItemsAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *ListItemsAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &listItemsDocs,
	}
}

func (a *ListItemsAction) Icon() *string {
	return nil
}

func (a *ListItemsAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"organization_id": shared.GetOrganizationsInput(),
	}
}

func (a *ListItemsAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[listItemsActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	url := "/v1/items?organization_id=" + input.OrganizationID

	items, err := shared.GetZohoClient(ctx.Auth.AccessToken, url)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (a *ListItemsAction) Auth() *sdk.Auth {
	return nil
}

func (a *ListItemsAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *ListItemsAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewListItemsAction() sdk.Action {
	return &ListItemsAction{}
}
