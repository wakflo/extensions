package actions

import (
	"fmt"

	"github.com/wakflo/extensions/internal/integrations/zohosalesiq/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"

	"github.com/wakflo/go-sdk/sdk"
)

type getVisitorsDetailsActionProps struct {
	ScreenName string `json:"screen_name"`
	ViewID     string `json:"view_id"`
}

type GetVisitorsDetailsAction struct{}

func (c *GetVisitorsDetailsAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func (c GetVisitorsDetailsAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (c GetVisitorsDetailsAction) Name() string {
	return "Get Visitor Details"
}

func (c GetVisitorsDetailsAction) Description() string {
	return "Get visitor details"
}

func (c GetVisitorsDetailsAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getVisitorsDetails,
	}
}

func (c GetVisitorsDetailsAction) Icon() *string {
	return nil
}

func (c GetVisitorsDetailsAction) SampleData() sdkcore.JSON {
	return nil
}

func (c GetVisitorsDetailsAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"screen_name": autoform.NewShortTextField().
			SetDisplayName("Screen name").
			SetDescription("Screen name").
			SetRequired(true).
			Build(),
		"view_id": autoform.NewShortTextField().
			SetDisplayName("View ID").
			SetRequired(false).
			Build(),
	}
}

func (c GetVisitorsDetailsAction) Auth() *sdk.Auth {
	return &sdk.Auth{
		Inherit: true,
	}
}

func (c GetVisitorsDetailsAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getVisitorsDetailsActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	viewID := input.ViewID

	if viewID == "" {
		viewID = "-1"
	}

	url := fmt.Sprintf("/%s/visitorsview/%s/visitors", input.ScreenName, viewID)

	visitors, err := shared.GetZohoClient(ctx.Auth.AccessToken, url)
	if err != nil {
		return nil, err
	}
	return visitors, nil
}

func NewGetVisitorsDetailsAction() sdk.Action {
	return &GetVisitorsDetailsAction{}
}
