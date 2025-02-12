package actions

import (
	"fmt"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type inviteCollaboratorActionProps struct {
	Name string `json:"name"`
}

type InviteCollaboratorAction struct{}

func (a *InviteCollaboratorAction) Name() string {
	return "Invite Collaborator"
}

func (a *InviteCollaboratorAction) Description() string {
	return "Invite Collaborator: Automatically send an invitation to a new team member or collaborator to join your workflow, ensuring seamless onboarding and efficient collaboration."
}

func (a *InviteCollaboratorAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *InviteCollaboratorAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &inviteCollaboratorDocs,
	}
}

func (a *InviteCollaboratorAction) Icon() *string {
	return nil
}

func (a *InviteCollaboratorAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"name": autoform.NewShortTextField().
			SetLabel("Name").
			SetRequired(true).
			SetPlaceholder("Your name").
			Build(),
	}
}

func (a *InviteCollaboratorAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[inviteCollaboratorActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	// implement action logic
	out := map[string]any{
		"message": fmt.Sprintf("Hello %s!", input.Name),
	}

	return out, nil
}

func (a *InviteCollaboratorAction) Auth() *sdk.Auth {
	return nil
}

func (a *InviteCollaboratorAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *InviteCollaboratorAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewInviteCollaboratorAction() sdk.Action {
	return &InviteCollaboratorAction{}
}
