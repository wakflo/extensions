package zohoinventory

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/zohoinventory/actions"
	"github.com/wakflo/extensions/internal/integrations/zohoinventory/shared"
	"github.com/wakflo/extensions/internal/integrations/zohoinventory/triggers"
	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewZohoInventory())

type ZohoInventory struct{}

func (n *ZohoInventory) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *ZohoInventory) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.SharedAuth,
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
