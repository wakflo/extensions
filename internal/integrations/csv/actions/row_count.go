package actions

import (
	"github.com/gocarina/gocsv"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type rowCountActionProps struct {
	Content string `json:"content"`
}

type RowCountAction struct{}

func (a *RowCountAction) Name() string {
	return "Row Count"
}

func (a *RowCountAction) Description() string {
	return "The Row Count integration action counts the number of rows in a specified table or dataset and returns the result as an output variable. This action is useful for tracking changes to data sets, monitoring data growth, or verifying data integrity."
}

func (a *RowCountAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *RowCountAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &rowCountDocs,
	}
}

func (a *RowCountAction) Icon() *string {
	return nil
}

func (a *RowCountAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"content": autoform.NewShortTextField().
			SetLabel("CSV Content").
			SetRequired(true).
			SetPlaceholder("csv content").
			Build(),
	}
}

func (a *RowCountAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[rowCountActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	var rows [][]string
	err = gocsv.UnmarshalString(input.Content, &rows)
	if err != nil {
		return nil, err
	}

	return map[string]any{"rowCount": len(rows)}, nil
}

func (a *RowCountAction) Auth() *sdk.Auth {
	return nil
}

func (a *RowCountAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"rowCount": 0,
	}
}

func (a *RowCountAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewRowCountAction() sdk.Action {
	return &RowCountAction{}
}
