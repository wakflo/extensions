package actions

import _ "embed"

//go:embed send_text_message.md
var sendMessageDocs string

//go:embed send_photo.md
var sendPhotoDocs string

//go:embed create_invite_link.md
var createInviteLinkDocs string

//go:embed get_chat_member.md
var getChatMemberDocs string

//go:embed get_chat_administrators.md
var getChatAdministratorsDocs string

// //go:embed get_updates.md
// var getUpdatesDocs string
