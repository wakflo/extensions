package actions

import (
	"errors"
	"fmt"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/discord/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type addRoleToMemberActionProps struct {
	GuildID string `json:"guild-id"`
	UserID  string `json:"user-id"`
	RoleID  string `json:"role-id"`
}

type AddRoleToMemberAction struct{}

// Metadata returns metadata about the action
func (a *AddRoleToMemberAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "add_role_to_member",
		DisplayName:   "Add Role to Member",
		Description:   "Add a role to a member in a Discord server",
		Type:          core.ActionTypeAction,
		Documentation: addRoleToMemberDocs,
		SampleOutput: map[string]any{
			"success":  true,
			"user_id":  "123456789012345678",
			"guild_id": "857347647235678912",
			"role_id":  "987654321098765432",
			"action":   "role_added",
			"reason":   "Promoted to moderator",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *AddRoleToMemberAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("add_role_to_member", "Add Role to Member")

	shared.RegisterGuildsInput(form, "Guilds", "List of guilds", true)

	form.TextField("user-id", "User ID").
		Required(true).
		HelpText("The ID of the user to add the role to")

	shared.RegisterRolesInput(form, "Roles", "List of roles", true)

	schema := form.Build()
	return schema
}

func (a *AddRoleToMemberAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[addRoleToMemberActionProps](ctx)
	if err != nil {
		return nil, err
	}
	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	if authCtx.Extra["token"] == "" {
		return nil, errors.New("missing discord bot token")
	}

	endpoint := "/guilds/" + input.GuildID + "/members/" + input.UserID + "/roles/" + input.RoleID

	response, err := shared.GetDiscordClient(authCtx.Extra["token"], endpoint, "PUT", nil)
	if err != nil {
		return nil, err
	}

	fmt.Println(response...)
	return map[string]interface{}{
		"success":  true,
		"user_id":  input.UserID,
		"guild_id": input.GuildID,
		"role_id":  input.RoleID,
		"action":   "role_added",
		"message":  "Role successfully added to member",
	}, nil
}

func (a *AddRoleToMemberAction) Auth() *core.AuthMetadata {
	return nil
}

func NewAddRoleToMemberAction() sdk.Action {
	return &AddRoleToMemberAction{}
}
