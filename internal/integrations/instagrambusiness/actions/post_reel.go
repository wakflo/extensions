package actions

import (
	"fmt"
	"time"

	"github.com/wakflo/extensions/internal/integrations/instagrambusiness/shared"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"

	"github.com/wakflo/go-sdk/sdk"
)

type postReelActionProps struct {
	InstagramAccountID string `json:"instagram_account_id"`
	Video              string `json:"video"`
	Caption            string `json:"caption"`
}

type PostReelAction struct{}

func (c *PostReelAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func (c PostReelAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (c PostReelAction) Name() string {
	return "Post Reel"
}

func (c PostReelAction) Description() string {
	return "Create and publish a new reel on your Instagram Business account"
}

func (c PostReelAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &postReelDocs,
	}
}

func (c PostReelAction) Icon() *string {
	return nil
}

func (c PostReelAction) SampleData() sdkcore.JSON {
	return nil
}

func (c PostReelAction) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"instagram_account_id": shared.GetInstagramAccountInput(
			"Select Instagram Account",
			"The Instagram Business account you want to post to",
			true,
		),
		"video": autoform.NewShortTextField().
			SetDisplayName("Video URL").
			SetDescription("A public URL for your video (max duration: 90 seconds)").
			SetRequired(true).
			Build(),
		"caption": autoform.NewLongTextField().
			SetDisplayName("Caption").
			SetDescription("Caption for your reel (include hashtags if desired)").
			SetRequired(false).
			Build(),
	}
}

func (c PostReelAction) Auth() *sdk.Auth {
	return &sdk.Auth{
		Inherit: true,
	}
}

func (c PostReelAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[postReelActionProps](ctx.BaseContext)
	if err != nil {
		return nil, fmt.Errorf("error parsing input: %w", err)
	}

	containerId, err := shared.CreateReelContainer(
		ctx.Auth.AccessToken,
		input.InstagramAccountID,
		input.Video,
		input.Caption,
	)
	if err != nil {
		return nil, err
	}

	time.Sleep(5 * time.Second)

	result, err := shared.PublishReel(
		ctx.Auth.AccessToken,
		input.InstagramAccountID,
		containerId,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func NewPostReelAction() sdk.Action {
	return &PostReelAction{}
}
