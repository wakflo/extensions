package actions

import (
	"errors"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/zohoinventory/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type listPaymentsActionProps struct {
	OrganizationID string `json:"organization_id"`
}

type ListPaymentsAction struct{}

func (a *ListPaymentsAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "list_payments",
		DisplayName:   "List Payments",
		Description:   "Retrieve a list of payments made to or from an account, including payment dates, amounts, and statuses.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: listPaymentsDocs,
		SampleOutput: map[string]any{
			"message": "hello world",
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *ListPaymentsAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("list_payments", "List Payments")

	shared.GetOrganizationsProp(form)

	schema := form.Build()

	return schema
}

func (a *ListPaymentsAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[listPaymentsActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the token source from the auth context
	tokenSource := ctx.Auth().Token
	if tokenSource == nil {
		return nil, errors.New("missing authentication token")
	}
	token := tokenSource.AccessToken

	url := "/v1/customerpayments/?organization_id=" + input.OrganizationID

	paymentList, err := shared.GetZohoClient(token, url)
	if err != nil {
		return nil, err
	}

	return paymentList, nil
}

func (a *ListPaymentsAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewListPaymentsAction() sdk.Action {
	return &ListPaymentsAction{}
}
