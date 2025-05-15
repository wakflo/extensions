package cin7

import (
	_ "embed"

	"github.com/wakflo/extensions/internal/integrations/cin7/actions"
	"github.com/wakflo/extensions/internal/integrations/cin7/shared"
	"github.com/wakflo/extensions/internal/integrations/cin7/triggers"

	"github.com/wakflo/go-sdk/v2"
	"github.com/wakflo/go-sdk/v2/core"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(NewCin7())

type Cin7 struct{}

func (n *Cin7) Metadata() sdk.IntegrationMetadata {
	return sdk.LoadMetadataFromFlo(Flow, ReadME)
}

func (n *Cin7) Auth() *core.AuthMetadata {
	return &core.AuthMetadata{
		Required: true,
		Schema:   shared.Cin7SharedAuth,
	}
}

func (n *Cin7) Triggers() []sdk.Trigger {
	return []sdk.Trigger{
		triggers.NewNewSalesTrigger(),
	}
}

func (n *Cin7) Actions() []sdk.Action {
	return []sdk.Action{
		actions.NewGetCustomersAction(),
		actions.NewGetCustomersAction(),
		actions.NewGetPaymentAction(),
		actions.NewGetPurchaseInvoiceAction(),
		actions.NewGetSalesOrderAction(),
		actions.NewGetCustomersAction(),
		actions.NewGetSalesListAction(),
	}
}

func NewCin7() sdk.Integration {
	return &Cin7{}
}
