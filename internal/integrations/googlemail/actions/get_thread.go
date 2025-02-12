package actions

import (
	"context"

	"github.com/wakflo/extensions/internal/integrations/googlemail/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type getThreadActionProps struct {
	ThreadID string `json:"id"`
	Format   string `json:"format"`
}

type GetThreadAction struct{}

func (a *GetThreadAction) Name() string {
	return "Get Thread"
}

func (a *GetThreadAction) Description() string {
	return "Retrieves a specific thread or conversation from a messaging platform, allowing you to incorporate its contents into your automated workflow."
}

func (a *GetThreadAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *GetThreadAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &getThreadDocs,
	}
}

func (a *GetThreadAction) Icon() *string {
	return nil
}

func (a *GetThreadAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"id": autoform.NewShortTextField().
			SetDisplayName("Thread ID").
			SetDescription("The thread Id of the mail to read").
			SetRequired(true).
			Build(),
		"format": autoform.NewSelectField().
			SetDisplayName("Format").
			SetDescription("The format of the email to read").
			SetOptions(shared.ViewMailFormat).
			SetRequired(false).
			Build(),
	}
}

func (a *GetThreadAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[getThreadActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	gmailService, err := gmail.NewService(context.Background(), option.WithTokenSource(*ctx.Auth.TokenSource))
	if err != nil {
		return nil, err
	}

	chosenFormat := input.Format
	if chosenFormat == "" {
		chosenFormat = "full"
	}

	mail, err := gmailService.Users.Messages.Get("me", input.ThreadID).
		Format(chosenFormat).
		Do()
	if err != nil {
		return nil, err
	}

	return mail, nil
}

func (a *GetThreadAction) Auth() *sdk.Auth {
	return nil
}

func (a *GetThreadAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *GetThreadAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewGetThreadAction() sdk.Action {
	return &GetThreadAction{}
}
