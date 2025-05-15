package actions

import (
	"errors"
	"fmt"

	"github.com/juicycleff/smartform/v1"

	"github.com/wakflo/extensions/internal/integrations/zohosalesiq/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type getVisitorsDetailsActionProps struct {
	ScreenName string `json:"screen_name"`
	ViewID     string `json:"view_id"`
}

type GetVisitorsDetailsAction struct{}

func (c *GetVisitorsDetailsAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		DisplayName:   "Get Visitor Details",
		Description:   "Retrieve visitor details from the CRM system, providing insights into visitor behavior and preferences.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: getVisitorsDetails,
		SampleOutput: map[string]any{
			"message": "hello world",
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (c GetVisitorsDetailsAction) Properties() *smartform.FormSchema {

	form := smartform.NewForm("get_visitors_details", "Get Visitor Details")

	form.TextField("screen_name", "Screen name").
		Placeholder("Enter a value for Screen name.").
		Required(true).
		HelpText("The screen name")

	form.TextField("view_id", "View ID").
		Placeholder("Enter a value for View ID.").
		Required(false).
		HelpText("The view ID")

	schema := form.Build()

	return schema
}

func (c GetVisitorsDetailsAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func (c GetVisitorsDetailsAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getVisitorsDetailsActionProps](ctx)
	if err != nil {
		return nil, err
	}

	tokenSource := ctx.Auth().Token
	if tokenSource == nil {
		return nil, errors.New("missing authentication token")
	}
	token := tokenSource.AccessToken

	viewID := input.ViewID

	if viewID == "" {
		viewID = "-1"
	}

	url := fmt.Sprintf("/%s/visitorsview/%s/visitors", input.ScreenName, viewID)

	visitors, err := shared.GetZohoClient(token, url)
	if err != nil {
		return nil, err
	}
	return visitors, nil
}

func NewGetVisitorsDetailsAction() sdk.Action {
	return &GetVisitorsDetailsAction{}
}
