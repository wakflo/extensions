package zohoinventory

import (
	"github.com/wakflo/extensions/internal/integrations/zohoinventory/actions"
	"github.com/wakflo/extensions/internal/integrations/zohoinventory/shared"
	"github.com/wakflo/extensions/internal/integrations/zohoinventory/triggers"
	"github.com/wakflo/go-sdk/sdk"
)

var Integration = sdk.Register(NewZohoInventory())

type ZohoInventory struct{}

func (n *ZohoInventory) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: true,
		Schema:   *shared.SharedAuth,
	}
}

func (n *ZohoInventory) Triggers() []sdk.Trigger {
	return []sdk.Trigger{
		triggers.NewNewPaymentTrigger(),
	}
}

func (n *ZohoInventory) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewListPaymentsAction(),
		actions.NewListInvoicesAction(),
		actions.NewListItemsAction(),
		actions.NewGetPaymentAction(),
		actions.NewGetInvoiceAction(),
	}
}

func NewZohoInventory() sdk.Integration {
	return &ZohoInventory{}
}
