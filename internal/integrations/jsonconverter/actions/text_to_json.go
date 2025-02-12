package actions

import (
	"github.com/wakflo/extensions/internal/integrations/jsonconverter/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type convertToJsonActionProps struct {
	Text string `json:"text"`
}

type ConvertToJsonAction struct{}

func (c *ConvertToJsonAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func (c ConvertToJsonAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (c ConvertToJsonAction) Name() string {
	return "Convert to Json"
}

func (c ConvertToJsonAction) Description() string {
	return "Returns the text in JSON"
}

func (c ConvertToJsonAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &textToJsonDocs,
	}
}

func (c ConvertToJsonAction) Icon() *string {
	return nil
}

func (c ConvertToJsonAction) SampleData() sdkcore.JSON {
	return nil
}

func (c ConvertToJsonAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"text": autoform.NewLongTextField().
			SetDisplayName("Text").
			SetDescription("Enter text to convert to JSON, starting with '{'").
			SetRequired(true).Build(),
	}
}

func (c ConvertToJsonAction) Auth() *sdk.Auth {
	return nil
}

func (c ConvertToJsonAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[convertToJsonActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}
	sdkcore.PrettyPrint(input)

	result, err := shared.ConvertTextToJSON(input.Text)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func NewConvertToJsonAction() sdk.Action {
	return &ConvertToJsonAction{}
}
