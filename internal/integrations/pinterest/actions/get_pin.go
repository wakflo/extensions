package actions

import (
	"errors"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/pinterest/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type getPinActionProps struct {
	BoardID string `json:"board_id"`
	PinID   string `json:"pin_id"`
}

type GetPinAction struct{}

// Metadata returns metadata about the action
func (a *GetPinAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "get_pin",
		DisplayName:   "Get Pin",
		Description:   "Get a pin from a board",
		Type:          core.ActionTypeAction,
		Documentation: getPinDocs,
		Icon:          "mdi:package-variant-closed-multiple",
		SampleOutput:  "",
		Settings:      core.ActionSettings{},
	}
}

func (a *GetPinAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("get_pin", "Get Pin")

	shared.RegisterBoardsProps(form)

	shared.RegisterBoardPinsProps(form)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *GetPinAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *GetPinAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[getPinActionProps](ctx)
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

	if input.BoardID == "" {
		return nil, errors.New("board ID is required")
	}

	if input.PinID == "" {
		return nil, errors.New("pin ID is required")
	}

	product, err := shared.GetPin(accessToken, input.PinID)

	return product, nil
}

func NewGetPinAction() sdk.Action {
	return &GetPinAction{}
}
