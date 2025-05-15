package actions

import (
	"errors"
	"fmt"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/zohoinventory/shared"

	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type getPaymentActionProps struct {
	OrganizationID string `json:"organization_id"`
	PaymentID      string `json:"payment_id"`
}

type GetPaymentAction struct{}

func (a *GetPaymentAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_payment",
		DisplayName:   "Get Payment",
		Description:   "Retrieves payment information from a specified source, such as an e-commerce platform or payment gateway, and stores it in the workflow's data storage.",
		Type:          sdkcore.ActionTypeAction,
		Documentation: getPaymentDocs,
		SampleOutput: map[string]any{
			"message": "hello world",
		},
		Settings: sdkcore.ActionSettings{},
	}
}

func (a *GetPaymentAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_payment", "Get Payment")

	shared.GetOrganizationsProp(form)

	form.TextField("payment_id", "Payment ID").
		Placeholder("Enter a payment ID").
		Required(true).
		HelpText("The ID of the customer payment to retrieve.")

	schema := form.Build()

	return schema
}

func (a *GetPaymentAction) Perform(ctx sdkcontext.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getPaymentActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Get the token source from the auth context
	tokenSource := ctx.Auth().Token
	if tokenSource == nil {
		return nil, errors.New("missing authentication token")
	}
	token := tokenSource.AccessToken

	url := fmt.Sprintf(
		"/v1/customerpayments/%s?organization_id=%s",
		input.PaymentID,
		input.OrganizationID,
	)

	payments, err := shared.GetZohoClient(token, url)
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func (a *GetPaymentAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

func NewGetPaymentAction() sdk.Action {
	return &GetPaymentAction{}
}
