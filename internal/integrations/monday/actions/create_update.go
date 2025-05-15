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

type createUpdateActionProps struct {
	ItemID string `json:"item_id,omitempty"`
	Body   string `json:"body,omitempty"`
}

type CreateUpdateAction struct{}

func (a *CreateUpdateAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "create_update",
		DisplayName:   "Create Update",
		Description:   "Create Update: This integration action allows you to create or update records in your target system based on the data provided in the trigger event. It enables you to maintain accurate and up-to-date information by automatically updating existing records or creating new ones when necessary.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: createUpdateDocs,
		SampleOutput: map[string]any{
			"message": "Hello World!",
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *CreateUpdateAction) Properties() *smartform.FormSchema {

	form := smartform.NewForm("create_update", "Create Update")

	form.TextField("item_id", "Item ID").
		Placeholder("Item ID").
		Required(true).
		HelpText("Item ID.")

	schema := form.Build()

	return schema
}

func (a *CreateUpdateAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[createUpdateActionProps](ctx)
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

	response, err := shared.MondayClient(ctx, query)
	if err != nil {
		return nil, err
	}

	update, ok := response["data"].(map[string]interface{})["create_update"].(map[string]interface{})
	if !ok {
		return nil, errors.New("failed to extract item from response")
	}

	return update, nil
}

func (a *CreateUpdateAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewCreateUpdateAction() sdk.Action {
	return &CreateUpdateAction{}
}
