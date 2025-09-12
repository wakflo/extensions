package actions

import (
	"errors"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/discord/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
)

type listGuildMembersActionProps struct {
	GuildID string `json:"guild-id"`
	After   string `json:"after"`
}

type ListGuildMembersAction struct{}

// Metadata returns metadata about the action
func (a *ListGuildMembersAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "list_guild_members",
		DisplayName:   "List Guild Members",
		Description:   "Get a list of guild members from a Discord server",
		Type:          core.ActionTypeAction,
		Documentation: listGuildMembersDocs,
		SampleOutput: map[string]any{
			"members": []interface{}{
				map[string]interface{}{
					"user": map[string]interface{}{
						"id":            "123456789012345678",
						"username":      "example_user",
						"discriminator": "1234",
						"avatar":        "a1b2c3d4e5f6g7h8i9j0",
						"bot":           false,
					},
					"nick":          "Nickname",
					"roles":         []string{"987654321098765432"},
					"joined_at":     "2021-01-01T00:00:00.000000+00:00",
					"premium_since": nil,
					"deaf":          false,
					"mute":          false,
				},
			},
			"count": 1,
		},
		Settings: core.ActionSettings{},
	}
}

func (a *ListGuildMembersAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("list_guild_members", "List Guild Members")

	shared.RegisterGuildsInput(form, "Guilds", "List of guilds", true)

	form.TextField("search", "Search").
		Required(true).
		HelpText("Search for a member")

	schema := form.Build()
	return schema
}

func (a *ListGuildMembersAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[listGuildMembersActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Build endpoint with query parameters
	endpoint := "/guilds/" + input.GuildID + "/members"

	// Add query parameters
	queryParams := "?"
	hasParams := false

	hasParams = true

	if input.After != "" {
		if hasParams {
			queryParams += "&"
		}
		queryParams += "after=" + input.After
		hasParams = true
	}

	if hasParams {
		endpoint += queryParams
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	if authCtx.Extra["token"] == "" {
		return nil, errors.New("missing discord bot token")
	}

	response, err := shared.GetDiscordClient(authCtx.Extra["token"], endpoint, "GET", nil)
	if err != nil {
		return nil, err
	}

	// Discord's guild members endpoint returns an array of member objects
	members := response

	// Return the members with additional metadata
	return map[string]interface{}{
		"members":  members,
		"count":    len(members),
		"guild_id": input.GuildID,
		"after":    input.After,
	}, nil
}

func (a *ListGuildMembersAction) Auth() *core.AuthMetadata {
	return nil
}

func NewListGuildMembersAction() sdk.Action {
	return &ListGuildMembersAction{}
}
