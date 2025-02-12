package actions

import (
	"errors"
	"fmt"
	"strings"

	"github.com/wakflo/extensions/internal/integrations/monday/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type createColumnActionProps struct {
	ColumnTitle string `json:"column_title,omitempty"`
	BoardID     string `json:"board_id,omitempty"`
	ColumnType  string `json:"column_type,omitempty"`
}

type CreateColumnAction struct{}

func (a *CreateColumnAction) Name() string {
	return "Create Column"
}

func (a *CreateColumnAction) Description() string {
	return "Create Column: Adds a new column to an existing table or spreadsheet, allowing you to customize your data structure and organization."
}

func (a *CreateColumnAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *CreateColumnAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createColumnDocs,
	}
}

func (a *CreateColumnAction) Icon() *string {
	return nil
}

func (a *CreateColumnAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"workspace_id": shared.GetWorkspaceInput(),
		"board_id":     shared.GetBoardInput("Board ID", "Select Board"),
		"column_title": autoform.NewShortTextField().
			SetDisplayName("Column Title").
			SetDescription("Group name").
			SetRequired(true).
			Build(),
		"column_type": autoform.NewSelectField().
			SetDisplayName("Column Type").
			SetOptions(shared.ColumnType).
			SetRequired(true).
			Build(),
	}
}

func (a *CreateColumnAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createColumnActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	fields := make(map[string]string)
	fields["board_id"] = fmt.Sprintf(`"%s"`, input.BoardID)
	fields["title"] = fmt.Sprintf(`"%s"`, input.ColumnTitle)
	fields["column_type"] = input.ColumnType

	fieldStrings := make([]string, 0, len(fields))
	for key, value := range fields {
		fieldStrings = append(fieldStrings, fmt.Sprintf("%s: %s", key, value))
	}

	mutation := fmt.Sprintf(`
		mutation {
  			create_column (%s) {
    			id
    			title
    			description
		}
}`, strings.Join(fieldStrings, "\n"))

	response, err := shared.MondayClient(ctx.BaseContext, mutation)
	if err != nil {
		return nil, err
		// return nil, errors.New("request not successful")
	}

	newColumn, ok := response["data"].(map[string]interface{})
	if !ok {
		return nil, errors.New("failed to extract column from response")
	}

	return newColumn, nil
}

func (a *CreateColumnAction) Auth() *sdk.Auth {
	return nil
}

func (a *CreateColumnAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *CreateColumnAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreateColumnAction() sdk.Action {
	return &CreateColumnAction{}
}
