package actions

import (
	"errors"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/gumroad/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type markasShippedActionProps struct {
	SaleID string `json:"sale_id"`
}

type MarkasShippedAction struct{}

// Metadata returns metadata about the action
func (a *MarkasShippedAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "mark_as_shipped",
		DisplayName:   "Mark Product as Shipped",
		Description:   "Mark a physical product as shipped.",
		Type:          core.ActionTypeAction,
		Documentation: markasShippedDocs,
		Icon:          "mdi:package-variant-closed-multiple",
		SampleOutput:  nil,
		Settings:      core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *MarkasShippedAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("mark_as_shipped", "Mark Product as Shipped")

	// Register sales selection field
	shared.RegisterSalesProps(form)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *MarkasShippedAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *MarkasShippedAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[markasShippedActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	if authCtx.AccessToken == "" {
		return nil, errors.New("missing Gumroad auth token")
	}
	accessToken := authCtx.Token.AccessToken

	if input.SaleID == "" {
		return nil, errors.New("sale ID is required")
	}

	sale, err := shared.DisableProduct(accessToken, input.SaleID)

	return sale, nil
}

func NewMarkasShippedAction() sdk.Action {
	return &MarkasShippedAction{}
}
