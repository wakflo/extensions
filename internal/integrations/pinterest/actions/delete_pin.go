package actions

import (
	"errors"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/pinterest/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type deletePinActionProps struct {
	BoardID string `json:"board_id"`
	PinID   string `json:"pin_id"`
}

type DeletePinAction struct{}

// Metadata returns metadata about the action
func (a *DeletePinAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "delete_pin",
		DisplayName:   "Delete a Pin",
		Description:   "Delete a pin from a board",
		Type:          core.ActionTypeAction,
		Documentation: getPinDocs,
		Icon:          "mdi:package-variant-closed-multiple",
		SampleOutput:  "",
		Settings:      core.ActionSettings{},
	}
}

func (a *DeletePinAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("delete_pin", "Delete a Pin")

	shared.RegisterBoardsProps(form)

	shared.RegisterBoardPinsProps(form)

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *DeletePinAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *DeletePinAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[deletePinActionProps](ctx)
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

	product, err := shared.DeletePin(accessToken, input.PinID)

	return product, nil
}

func NewDeletePinAction() sdk.Action {
	return &DeletePinAction{}
}
