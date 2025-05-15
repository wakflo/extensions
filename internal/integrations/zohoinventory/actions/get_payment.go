package actions

import (
	// "errors"
	"fmt"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/zohoinventory/shared"

	// "github.com/wakflo/go-sdk/autoform"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

type getPaymentActionProps struct {
	OrganizationID string `json:"organization_id"`
	PaymentID      string `json:"payment_id"`
}

type GetPaymentAction struct{}

// func (a *GetPaymentAction) Name() string {
// 	return "Get Payment"
// }

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

// func (a *GetPaymentAction) Description() string {
// 	return "Retrieves payment information from a specified source, such as an e-commerce platform or payment gateway, and stores it in the workflow's data storage."
// }

// func (a *GetPaymentAction) GetType() sdkcore.ActionType {
// 	return sdkcore.ActionTypeNormal
// }

// func (a *GetPaymentAction) Documentation() *sdk.OperationDocumentation {
// 	return &sdk.OperationDocumentation{
// 		Documentation: &getPaymentDocs,
// 	}
// }

// func (a *GetPaymentAction) Icon() *string {
// 	return nil
// }

// func (a *GetPaymentAction) Properties() map[string]*sdkcore.AutoFormSchema {
// 	return map[string]*sdkcore.AutoFormSchema{
// 		"organization_id": shared.GetOrganizationsInput(),
// 		"payment_id": autoform.NewShortTextField().
// 			SetDisplayName("Payment ID").
// 			SetDescription("The ID of the customer payment to retrieve").
// 			SetRequired(true).
// 			Build(),
// 	}
// }

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
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	// if ctx.Auth.AccessToken == "" {
	// 	return nil, errors.New("missing Zoho auth token")
	// }

	url := fmt.Sprintf(
		"/v1/customerpayments/%s?organization_id=%s",
		input.PaymentID,
		input.OrganizationID,
	)

	payments, err := shared.GetZohoClient(authCtx.AccessToken, url)
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func (a *GetPaymentAction) Auth() *sdkcore.AuthMetadata {
	return nil
}

// func (a *GetPaymentAction) SampleData() sdkcore.JSON {
// 	return map[string]any{
// 		"message": "Hello World!",
// 	}
// }

// func (a *GetPaymentAction) Settings() sdkcore.ActionSettings {
// 	return sdkcore.ActionSettings{}
// }

func NewGetPaymentAction() sdk.Action {
	return &GetPaymentAction{}
}
