package zohoinventory

import (
	"errors"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type listItemsOperationProps struct {
	OrganizationID string `json:"organization_id"`
}

type ListItemsOperation struct {
	options *sdk.OperationInfo
}

func NewListItemsOperation() sdk.IOperation {
	return &ListItemsOperation{
		options: &sdk.OperationInfo{
			Name:        "List Items",
			Description: "Get list of items",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"organization_id": getOrganizationsInput(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *ListItemsOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Token == nil {
		return nil, errors.New("missing Zoho auth token")
	}

	input := sdk.InputToType[listItemsOperationProps](ctx)

	url := "/v1/items?organization_id=" + input.OrganizationID

	items, err := getZohoClient(ctx.Auth.AccessToken, url)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (c *ListItemsOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *ListItemsOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
