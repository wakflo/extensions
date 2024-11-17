package zohosalesiq

import (
	"errors"
	"fmt"

	"github.com/wakflo/go-sdk/autoform"

	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type getVisitorsDetailsOperationProps struct {
	ScreenName string `json:"screen_name"`
	ViewID     string `json:"view_id"`
}

type GetVisitorsDetailsOperation struct {
	options *sdk.OperationInfo
}

func NewGetVisitorsDetailsOperation() sdk.IOperation {
	return &ListChatsOperation{
		options: &sdk.OperationInfo{
			Name:        "Get Visitor Details",
			Description: "Get visitor details",
			RequireAuth: true,
			Auth:        sharedAuth,
			Input: map[string]*sdkcore.AutoFormSchema{
				"screen_name": autoform.NewShortTextField().
					SetDisplayName("Screen name").
					SetDescription("Screen name").
					SetRequired(true).
					Build(),
				"view_id": autoform.NewShortTextField().
					SetDisplayName("View ID").
					SetRequired(false).
					Build(),
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *GetVisitorsDetailsOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	if ctx.Auth.Token == nil {
		return nil, errors.New("missing Zoho auth token")
	}

	input := sdk.InputToType[getVisitorsDetailsOperationProps](ctx)

	viewID := input.ViewID

	if viewID == "" {
		viewID = "-1"
	}

	url := fmt.Sprintf("/%s/visitorsview/%s/visitors", input.ScreenName, viewID)

	visitors, err := getZohoClient(ctx.Auth.AccessToken, url)
	if err != nil {
		return nil, err
	}
	return visitors, nil
}

func (c *GetVisitorsDetailsOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *GetVisitorsDetailsOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
