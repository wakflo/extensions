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

type createUpdateActionProps struct {
	ItemID string `json:"item_id,omitempty"`
	Body   string `json:"body,omitempty"`
}

type CreateUpdateAction struct{}

func (a *CreateUpdateAction) Name() string {
	return "Create Update"
}

func (a *CreateUpdateAction) Description() string {
	return "Create Update: This integration action allows you to create or update records in your target system based on the data provided in the trigger event. It enables you to maintain accurate and up-to-date information by automatically updating existing records or creating new ones when necessary."
}

func (a *CreateUpdateAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *CreateUpdateAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &createUpdateDocs,
	}
}

func (a *CreateUpdateAction) Icon() *string {
	return nil
}

func (a *CreateUpdateAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"item_id": autoform.NewShortTextField().
			SetDisplayName("Item ID").
			SetDescription("Item ID").
			SetRequired(true).
			Build(),
		"body": autoform.NewLongTextField().
			SetDisplayName("Body").
			SetRequired(true).
			Build(),
	}
}

func (a *CreateUpdateAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createUpdateActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	fields := make(map[string]string)
	fields["item_id"] = fmt.Sprintf(`"%s"`, input.ItemID)
	fields["body"] = fmt.Sprintf(`"%s"`, input.Body)

	fieldStrings := make([]string, 0, len(fields))
	for key, value := range fields {
		fieldStrings = append(fieldStrings, fmt.Sprintf("%s: %s", key, value))
	}

	query := fmt.Sprintf(`
		mutation {
  create_update (%s) {
    	id
		body
  }
}`, strings.Join(fieldStrings, "\n"))

	response, err := shared.MondayClient(ctx.BaseContext, query)
	if err != nil {
		return nil, err
	}

	update, ok := response["data"].(map[string]interface{})["create_update"].(map[string]interface{})
	if !ok {
		return nil, errors.New("failed to extract item from response")
	}

	return update, nil
}

func (a *CreateUpdateAction) Auth() *sdk.Auth {
	return nil
}

func (a *CreateUpdateAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *CreateUpdateAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewCreateUpdateAction() sdk.Action {
	return &CreateUpdateAction{}
}
