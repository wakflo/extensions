package actions

import _ "embed"

//go:embed send_message.md
var sendMessageDocs string

//go:embed create_channel.md
var createChannelDocs string

//go:embed delete_channel.md
var deleteChannelDocs string

//go:embed rename_channel.md
var renameChannelDocs string

//go:embed find_channel.md
var findChannelDocs string

//go:embed find_guild_member.md
var listGuildMembersDocs string

//go:embed remove_member_from_guild.md
var removeGuildMemberDocs string

//go:embed ban_guild_member.md
var banGuildMemberDocs string

//go:embed remove_ban_from_user.md
var unbanGuildMemberDocs string

//go:embed add_role_to_member.md
var addRoleToMemberDocs string
