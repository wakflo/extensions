package actions

import (
	"errors"
	"fmt"
	"strings"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/monday/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type createColumnActionProps struct {
	ColumnTitle string `json:"column_title,omitempty"`
	BoardID     string `json:"board_id,omitempty"`
	ColumnType  string `json:"column_type,omitempty"`
}

type CreateColumnAction struct{}

func (a *CreateColumnAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_column",
		DisplayName:   "Create Column",
		Description:   "Create Column: Adds a new column to an existing table or spreadsheet, allowing you to customize your data structure and organization.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: createColumnDocs,
		SampleOutput: map[string]any{
			"data": map[string]interface{}{
				"id":   "123456789",
				"name": "New Column",
			},
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *CreateColumnAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("create_column", "Create Column")

	shared.GetWorkspaceProp(form)

	shared.GetBoardProp("board_id", "Board ID", "Select Board", form)

	form.TextField("column_title", "Column Title").
		Placeholder("New Column").
		Required(true).
		HelpText("Column name.")

	form.SelectField("column_type", "Column Type").
		Placeholder("Select Column Type").
		Required(true).
		AddOptions(shared.ColumnType...).
		HelpText("Column type.")

	schema := form.Build()

	return schema
}

func (a *CreateColumnAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createColumnActionProps](ctx)
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

	response, err := shared.MondayClient(ctx, mutation)
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

func (a *CreateColumnAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewCreateColumnAction() sdk.Action {
	return &CreateColumnAction{}
}
