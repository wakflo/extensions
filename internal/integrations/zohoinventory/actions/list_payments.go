package actions

import (
	"github.com/wakflo/extensions/internal/integrations/zohoinventory/shared"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type listPaymentsActionProps struct {
	OrganizationID string `json:"organization_id"`
}

type ListPaymentsAction struct{}

func (a *ListPaymentsAction) Name() string {
	return "List Payments"
}

func (a *ListPaymentsAction) Description() string {
	return "Retrieve a list of payments made to or from an account, including payment dates, amounts, and statuses."
}

func (a *ListPaymentsAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *ListPaymentsAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &listPaymentsDocs,
	}
}

func (a *ListPaymentsAction) Icon() *string {
	return nil
}

func (a *ListPaymentsAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"organization_id": shared.GetOrganizationsInput(),
	}
}

func (a *ListPaymentsAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[listPaymentsActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	url := "/v1/customerpayments/?organization_id=" + input.OrganizationID

	paymentList, err := shared.GetZohoClient(ctx.Auth.AccessToken, url)
	if err != nil {
		return nil, err
	}

	return paymentList, nil
}

func (a *ListPaymentsAction) Auth() *sdk.Auth {
	return nil
}

func (a *ListPaymentsAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *ListPaymentsAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewListPaymentsAction() sdk.Action {
	return &ListPaymentsAction{}
}
