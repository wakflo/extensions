package actions

import _ "embed"

//go:embed send_direct_message.md
var sendDirectMessageDocs string

//go:embed send_private_channel_message.md
var sendPrivateChannelMessageDocs string

//go:embed send_public_channel_message.md
var sendPublicChannelMessageDocs string
